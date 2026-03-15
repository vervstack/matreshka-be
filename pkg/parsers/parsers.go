package parsers

import (
	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

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
