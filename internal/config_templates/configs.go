package config_templates

import (
	"time"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/pkg/matreshka"
)

var ErrUnsupportedConfigType = errors.New("unsupported config type")

func InitNewConfig(configName string, configType matreshka_api.ConfigType) (*evon.Node, error) {
	switch configType {
	case matreshka_api.ConfigType_verv:
		return InitVervConfig(configName)
	case matreshka_api.ConfigType_kv:
		return &evon.Node{}, nil
	default:
		return nil, ErrUnsupportedConfigType
	}
}

func InitVervConfig(name string) (*evon.Node, error) {
	newCfg := matreshka.NewEmptyConfig()
	newCfg.AppInfo = matreshka.AppInfo{
		Name:            name,
		Version:         "v0.0.1",
		StartupDuration: time.Second * 5,
	}
	nodes, err := evon.MarshalEnv(&newCfg)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling config")
	}
	return nodes, nil
}
