package cfg_service

import (
	"context"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service/user_errors"
	"go.vervstack.ru/matreshka/pkg/parsers"
)

func (c *CfgService) List(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error) {
	resp, err := c.configStorage.ListConfigs(ctx, req)
	if err != nil {
		return domain.ListConfigsResponse{}, errors.Wrap(err)
	}

	return resp, nil
}

func (c *CfgService) GetConfigInfo(ctx context.Context, configName string) (domain.ConfigInfo, error) {
	return c.getConfigInfo(ctx, configName)
}

func (c *CfgService) GetConfigNodes(ctx context.Context, configName string, ver string) (domain.ConfigNodes, error) {
	info, err := c.GetConfigInfo(ctx, configName)
	if err != nil {
		return domain.ConfigNodes{}, errors.Wrap(err)
	}

	nodes, err := c.configStorage.GetConfigNodes(ctx, configName, ver)
	if err != nil {
		return domain.ConfigNodes{}, errors.Wrap(err)
	}

	if nodes == nil {
		return domain.ConfigNodes{}, nil
	}

	//switch configType {
	//case api.ConfigType_pg:
	//	toSnake(nodes)
	//}

	return domain.ConfigNodes{
		Info:  info,
		Nodes: nodes,
	}, nil
}
func (c *CfgService) GetConfigRaw(ctx context.Context, r domain.GetConfigRawReq) (domain.ConfigRawContent, error) {
	configInfo, err := c.getConfigInfo(ctx, r.Name)
	if err != nil {
		return domain.ConfigRawContent{}, errors.Wrap(err)
	}

	var rawContent []byte

	switch configInfo.Type {
	case api.ConfigType_kv, api.ConfigType_verv:
		var nodes *evon.Node
		nodes, err = c.configStorage.GetConfigNodes(ctx, r.Name, r.Version)
		if err != nil {
			return domain.ConfigRawContent{}, errors.Wrap(err)
		}

		switch r.Format {
		case api.Format_yaml:
			rawContent, err = parsers.NodeToYaml(nodes, configInfo.Type)
			if err != nil {
				return domain.ConfigRawContent{}, errors.Wrap(err)
			}
		case api.Format_env:
			rawContent = evon.Marshal(nodes.InnerNodes)
		default:
			return domain.ConfigRawContent{},
				errors.Wrap(user_errors.ErrNotImplemented, "config format unsupported")
		}

	default:
		rawContent, err = c.configStorage.GetConfigRawContent(ctx, r.Name, r.Version)
		if err != nil {
			return domain.ConfigRawContent{}, errors.Wrap(err)
		}
	}

	if rawContent == nil {
		return domain.ConfigRawContent{}, nil
	}

	return domain.ConfigRawContent{
		Info: configInfo,
		Raw:  rawContent,
	}, nil
}

func (c *CfgService) getConfigInfo(ctx context.Context, configName string) (domain.ConfigInfo, error) {
	cfgBaseInfo, err := c.configStorage.GetConfigByName(ctx, configName)
	if err != nil {
		return domain.ConfigInfo{}, errors.Wrap(err)
	}

	versions, err := c.configStorage.GetVersions(ctx, configName)
	if err != nil {
		return domain.ConfigInfo{}, errors.Wrap(err)
	}

	return domain.ConfigInfo{
		ConfigBase:     cfgBaseInfo,
		ConfigVersions: versions,
	}, nil
}
