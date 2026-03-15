package parsers

import (
	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

	errors "github.com/Red-Sock/trace-errors"

	"go.vervstack.ru/matreshka/internal/service/user_errors"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type parseMappingKey struct {
	Format     api.Format
	ConfigType api.ConfigType
}

var evonConfigParsersMapping map[parseMappingKey]func([]byte) (*evon.Node, error)

func init() {
	evonConfigParsersMapping = map[parseMappingKey]func([]byte) (*evon.Node, error){
		{
			Format:     api.Format_env,
			ConfigType: api.ConfigType_kv,
		}: ParseEvonConfig,

		{
			Format:     api.Format_env,
			ConfigType: api.ConfigType_verv,
		}: ParseEvonConfig,

		{
			Format:     api.Format_yaml,
			ConfigType: api.ConfigType_verv,
		}: ParseVervYamlToEvon,

		{
			Format:     api.Format_yaml,
			ConfigType: api.ConfigType_kv,
		}: ParseYamlToEvon,
	}
}

func ParseEvon(bytes []byte, format api.Format, configType api.ConfigType) (*evon.Node, error) {
	key := parseMappingKey{
		Format:     format,
		ConfigType: configType,
	}

	parser, ok := evonConfigParsersMapping[key]
	if !ok {
		return nil, rerrors.Wrap(user_errors.ErrNotImplemented, "unsupported config type and format pair", key)
	}

	return parser(bytes)
}

func ParseEvonConfig(cfg []byte) (*evon.Node, error) {
	n := evon.ParseToNodes(cfg)

	return n[""], nil
}

// TODO When marshalled from Matreshka original config to map - marshalles wrong.
// Must marshall as verv config, but puts all the values into innerNodes (must put main value in 'value' field in root)
func ParseYamlToEvon(cfg []byte) (*evon.Node, error) {
	yamlMap := map[string]any{}
	err := yaml.Unmarshal(cfg, yamlMap)
	if err != nil {
		return nil, rerrors.Wrap(err, "")
	}

	env, err := evon.MarshalEnv(yamlMap)
	if err != nil {
		return nil, rerrors.Wrap(err, "")
	}

	return env, nil
}

// TODO When marshalled from Matreshka original config to map - marshalles wrong.
// Must marshall as verv config, but puts all the values into innerNodes (must put main value in 'value' field in root)
func ParseVervYamlToEvon(cfgBytes []byte) (*evon.Node, error) {
	cfg := matreshka.NewEmptyConfig()
	err := yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "")
	}

	env, err := evon.MarshalEnv(&cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "")
	}

	return env, nil
}

func NodeToYaml(node *evon.Node, configType api.ConfigType) ([]byte, error) {
	switch configType {
	case api.ConfigType_kv:
		return kvToYaml(node)

	case api.ConfigType_verv:
		return vervToYaml(node)
	default:
		return nil, errors.Wrap(user_errors.ErrNotImplemented, "unsupported config type to parse from node to yaml")
	}
}

func vervToYaml(node *evon.Node) ([]byte, error) {
	nodeStorage := evon.NodesToStorage(node)

	matreshkaConf := matreshka.NewEmptyConfig()

	err := evon.UnmarshalWithNodes(nodeStorage, &matreshkaConf)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling config")
	}

	config, err := matreshkaConf.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling to yaml")
	}

	return config, nil
}

func kvToYaml(node *evon.Node) ([]byte, error) {
	nodeStorage := evon.NodesToStorage(node)

	m := make(map[string]any)
	err := evon.UnmarshalWithNodes(nodeStorage, m,
		evon.WithSnakeUnmarshal(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling config")
	}

	config, err := yaml.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling to yaml")
	}

	return config, nil
}
