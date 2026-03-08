package matreshka_api_impl

import (
	"context"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/matreshka/pkg/matreshka"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (s *Impl) StoreConfig(ctx context.Context, req *api.StoreConfig_Request) (
	*api.StoreConfig_Response, error) {

	//parsedConfigType, name := ParseConfigName(req.ConfigName)
	//
	//ConfigType := api.ConfigType_plain
	//
	//if parsedConfigType != nil {
	//	ConfigType = *parsedConfigType
	//}
	//
	//cfgName := domain.NewConfigName(ConfigType, name)
	//
	////_, err := s.evonConfigService.Create(ctx, cfgName)
	////if err != nil {
	////	if !rerrors.Is(err, user_errors.ErrAlreadyExists) {
	////		return nil, rerrors.Wrap(err, "error creating config")
	////	}
	////}
	//
	//replaceReq := domain.ReplaceConfigReq{
	//	Name:    cfgName,
	//	Version: toolbox.Coalesce(req.GetVersion(), domain.MasterVersion),
	//	Config:  nil,
	//}
	//
	//switch req.Format {
	//case api.Format_env:
	//	replaceReq.Config, err = fromEvon(req.Config)
	//default:
	//	switch ConfigType {
	//	case api.ConfigType_verv:
	//		replaceReq.Config, err = fromVervYamlToEvon(req.Config)
	//	default:
	//		replaceReq.Config, err = fromPlainYamlToEvon(req.Config)
	//	}
	//}
	//if err != nil {
	//	return nil, rerrors.Wrap(err, "error parsing config", codes.InvalidArgument)
	//}
	//
	//err = s.evonConfigService.Replace(ctx, replaceReq)
	//if err != nil {
	//	return nil, rerrors.Wrap(err, "error updating config")
	//}

	return &api.StoreConfig_Response{}, nil
}

// TODO When marshalled from Matreshka original config to map - marshalles wrong.
// Must marshall as verv config, but puts all the values into innerNodes (must put main value in 'value' field in root)
func fromPlainYamlToEvon(cfg []byte) (*evon.Node, error) {
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
func fromVervYamlToEvon(cfgBytes []byte) (*evon.Node, error) {
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

func fromEvon(cfg []byte) (*evon.Node, error) {
	n := &evon.Node{}
	err := evon.Unmarshal(cfg, n)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}
	return n, nil
}
