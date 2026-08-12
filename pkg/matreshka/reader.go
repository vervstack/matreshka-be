package matreshka

import (
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
	"go.vervstack.ru/matreshka/pkg/matreshka/server"
)

const (
	VervName = "VERV_NAME"
)

func NewEmptyConfig() AppConfig {
	return AppConfig{
		AppInfo:     AppInfo{},
		DataSources: make(DataSources, 0),
		Servers:     make(Servers),
		Environment: make(Environment, 0),
	}
}

// configSource is a single input to ReadConfig: either a path to a YAML file
// on disk (path set) or raw in-memory YAML bytes (bytes set).
type configSource struct {
	path  string
	bytes []byte
}

// readOptions accumulates everything ReadOption values can configure:
// the ordered list of YAML/bytes sources, and an optional .env file path.
type readOptions struct {
	srcs        []configSource
	envFilePath string
}

// ReadOption configures how ReadConfig resolves its master config.
type ReadOption func(*readOptions)

// WithConfigPaths adds one or more YAML files to be read from disk, in order.
func WithConfigPaths(paths ...string) ReadOption {
	return func(opts *readOptions) {
		for _, p := range paths {
			opts.srcs = append(opts.srcs, configSource{path: p})
		}
	}
}

// WithConfigBytes adds a raw in-memory YAML config source (e.g. an embedded
// default/skeleton config compiled into the binary).
func WithConfigBytes(b []byte) ReadOption {
	return func(opts *readOptions) {
		opts.srcs = append(opts.srcs, configSource{bytes: b})
	}
}

// WithEnvFile sets an optional .env file whose KEY=VALUE values are overlaid
// on top of the merged YAML/bytes config, but are still overridden by real OS
// environment variables. If not passed, behavior is unchanged: only YAML/bytes
// sources and real OS env vars are used.
func WithEnvFile(path string) ReadOption {
	return func(opts *readOptions) {
		opts.envFilePath = path
	}
}

// ReadConfig builds a master config from the given sources (files and/or raw
// bytes), in the order the options were passed — earlier sources win
// per-entry over later ones (see MergeConfigs), then overlays OS env vars.
func ReadConfig(opts ...ReadOption) (masterConfig AppConfig, err error) {
	masterConfig = NewEmptyConfig()

	var ro readOptions
	for _, o := range opts {
		o(&ro)
	}

	var errs []error
	for _, s := range ro.srcs {
		var cfg AppConfig
		var srcErr error

		if s.path != "" {
			cfg, srcErr = getFromFile(s.path)
			if srcErr != nil {
				srcErr = rerrors.Wrapf(srcErr, "error reading config at %s", s.path)
			}
		} else {
			cfg, srcErr = ParseConfig(s.bytes)
			if srcErr != nil {
				srcErr = rerrors.Wrap(srcErr, "error parsing config bytes")
			}
		}

		if srcErr != nil {
			errs = append(errs, srcErr)
			continue
		}

		masterConfig = MergeConfigs(masterConfig, cfg)
	}

	var envFileLines []string
	if ro.envFilePath != "" {
		var envFileErr error
		envFileLines, envFileErr = parseEnvFile(ro.envFilePath)
		if envFileErr != nil {
			errs = append(errs, rerrors.Wrapf(envFileErr, "error reading env file at %s", ro.envFilePath))
		}
	}

	if len(errs) != 0 {
		return masterConfig, stderrors.Join(errs...)
	}

	masterConfig, err = enrichWithEnv(masterConfig, envFileLines)
	if err != nil {
		return masterConfig, rerrors.Wrap(err, "error enriching master config with env vars")
	}

	return masterConfig, nil
}

// ReadConfigs reads and merges YAML config files from paths, then overlays
// OS env vars.
//
// Deprecated: use ReadConfig(WithConfigPaths(paths...)) instead.
func ReadConfigs(paths ...string) (AppConfig, error) {
	return ReadConfig(WithConfigPaths(paths...))
}

func ParseConfig(in []byte) (AppConfig, error) {
	a := NewEmptyConfig()

	err := a.Unmarshal(in)
	if err != nil {
		return a, err
	}

	return a, nil
}

func MergeConfigs(master, slave AppConfig) AppConfig {
	if master.Name == "" {
		master.Name = slave.Name
	}
	if master.Version == "" {
		master.Version = slave.Version
	}
	if master.StartupDuration == 0 {
		master.StartupDuration = slave.StartupDuration
	}

	for _, slaveVal := range slave.Environment {
		var mv *environment.Variable
		for _, masterVal := range master.Environment {
			if masterVal.Name == slaveVal.Name {
				mv = masterVal
				break
			}
		}
		if mv != nil {
			continue
		}

		master.Environment = append(master.Environment, slaveVal)
	}

	serversByName := map[string]*server.Server{}

	for _, s := range master.Servers {
		serversByName[strings.ToLower(ServerName(s.Name))] = s
	}

	for slavePort, slaveServer := range slave.Servers {
		_, ok := master.Servers[slavePort]
		if ok {
			continue
		}

		slaveServerName := strings.ToLower(ServerName(slaveServer.Name))

		_, ok = serversByName[slaveServerName]
		if ok {
			// already have a server that resolves to the same destination field
			continue
		}

		master.Servers[slavePort] = slaveServer
	}

	for i := range slave.DataSources {
		if master.DataSources.get(slave.DataSources[i].GetName()) == nil {
			master.DataSources = append(master.DataSources, slave.DataSources[i])
		}
	}

	master.ServiceDiscovery.MakoshUrl =
		toolbox.Coalesce(master.ServiceDiscovery.MakoshUrl, slave.ServiceDiscovery.MakoshUrl)

	master.ServiceDiscovery.MakoshToken =
		toolbox.Coalesce(master.ServiceDiscovery.MakoshToken, slave.ServiceDiscovery.MakoshToken)

	for _, slaveOverride := range slave.ServiceDiscovery.Overrides {
		found := false
		for _, masterOverride := range master.ServiceDiscovery.Overrides {
			if masterOverride.ServiceName == slaveOverride.ServiceName {
				found = true
				break
			}
		}

		if !found {
			master.ServiceDiscovery.Overrides = append(master.ServiceDiscovery.Overrides, slaveOverride)
		}
	}
	return master
}

// enrichWithEnv overlays the given .env file lines (lower priority) and then
// real OS environment variables (highest priority) onto masterConfig. envFileLines
// may be nil, in which case only real OS env vars are applied (unchanged behavior).
func enrichWithEnv(masterConfig AppConfig, envFileLines []string) (enrichedConfig AppConfig, err error) {
	projectName := strings.ToUpper(toolbox.Coalesce(os.Getenv(VervName), masterConfig.ModuleName()))

	// Storage in Evon format (e.g. object_sub-field-name_leaf-field-name)
	masterEvonCfg, err := evon.MarshalEnv(&masterConfig)
	if err != nil {
		return masterConfig, rerrors.Wrap(err, "error marshalling to env")
	}

	masterEvonStorage := evon.NodeStorage{}
	masterEvonStorage.AddNode(masterEvonCfg)

	// Pointers to evon nodes presented in default environment variable format
	envNamePointers := map[string]*evon.Node{}
	for name, node := range masterEvonStorage {
		envNamePointers[strings.ReplaceAll(name, "-", "_")] = node
	}

	// Lower priority: .env file values, applied first so real OS env vars can
	// still override them below.
	applyEnvOverlay(envFileLines, projectName, masterEvonStorage, envNamePointers)

	// Highest priority: real OS environment variables, applied last so they
	// always win on conflict.
	applyEnvOverlay(os.Environ(), projectName, masterEvonStorage, envNamePointers)

	masterConfig = NewEmptyConfig()
	err = evon.UnmarshalWithNodes(masterEvonStorage, &masterConfig)
	if err != nil {
		return masterConfig, rerrors.Wrap(err, "error unmarshalling back to config")
	}

	sort.Slice(masterConfig.Environment, func(i, j int) bool {
		return masterConfig.Environment[i].Name < masterConfig.Environment[j].Name
	})

	return masterConfig, nil
}

// applyEnvOverlay overlays "KEY=VALUE" lines onto the given evon node storage,
// mutating the nodes it finds matches for in place. Used both for .env file
// lines and for real OS environment variables (os.Environ()).
func applyEnvOverlay(
	lines []string,
	projectName string,
	masterEvonStorage evon.NodeStorage,
	envNamePointers map[string]*evon.Node,
) {
	const environmentEvonPart = "ENVIRONMENT"

	nodeFinders := []func(originalName string) *evon.Node{
		// Simply extract from storage
		func(originalName string) *evon.Node {
			return masterEvonStorage[originalName]
		},
		// Try to find inside environment object
		func(originalName string) *evon.Node {
			originalName = environmentEvonPart + evon.ObjectSplitter + originalName

			return masterEvonStorage[originalName]
		},
		// Use environment style to find from
		func(originalName string) *evon.Node {
			originalName = environmentEvonPart + evon.ObjectSplitter + originalName
			originalName = strings.ReplaceAll(originalName, "-", "_")
			return envNamePointers[originalName]
		},
		// Use full-underscore style directly, with no section prefix
		// (covers DataSources, Servers, etc. — same tolerance Environment already gets,
		// needed since real env var names cannot contain '-')
		func(originalName string) *evon.Node {
			return envNamePointers[originalName]
		},
	}

	for _, variable := range lines {
		idx := strings.Index(variable, "=")
		if idx == -1 {
			continue
		}

		// Just in case. Validate that env variable has certain project name prefix.
		// This allows user to set variables in short form.
		// Instead of VELEZ_SHUT_DOWN_ON_EXIT use SHUT_DOWN_ON_EXIT as name
		variableName := strings.ToUpper(variable[:idx])

		variableValue := variable[idx+1:]

		if strings.HasPrefix(variableName, projectName) {
			variableName = variableName[len(projectName)+1:]
		}

		var node *evon.Node

		for _, nf := range nodeFinders {
			node = nf(variableName)
			if node != nil {
				break
			}
		}

		if node == nil {
			continue
		}

		node.Value = variableValue
	}
}

// parseEnvFile reads a .env file and returns its content as a slice of
// trimmed "KEY=VALUE" strings. Blank lines and lines starting with '#' (after
// trimming) are skipped. No quote stripping or "export" keyword handling is
// performed — this matches the format evon.Marshal already emits.
func parseEnvFile(path string) ([]string, error) {
	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(bts), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func getFromFile(pth string) (AppConfig, error) {
	f, err := os.Open(pth)
	if err != nil {
		return NewEmptyConfig(), err
	}

	defer func() {
		closerErr := f.Close()
		if err == nil {
			err = closerErr
			return
		}

		err = stderrors.Join(err, closerErr)
	}()

	fi, err := f.Stat()
	if err != nil {
		return AppConfig{}, rerrors.Wrap(err, "error getting config file info")
	}

	if fi.Size() > 1_000_000 {
		return AppConfig{}, fmt.Errorf("config file too large (more than a 1 MB)")
	}

	c := NewEmptyConfig()

	bts, err := io.ReadAll(f)
	if err != nil {
		return c, rerrors.Wrap(err, "error reading file")
	}

	err = c.Unmarshal(bts)
	if err != nil {
		return c, rerrors.Wrap(err, "error decoding config to struct")
	}

	return c, nil
}
