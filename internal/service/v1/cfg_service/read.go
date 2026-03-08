package cfg_service

import (
	"context"
	"strings"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
)

func (c *CfgService) GetConfigWithNodes(ctx context.Context, configName string, ver string) (domain.ConfigWithNodes, error) {
	nodes, err := c.configStorage.GetConfigNodes(ctx, configName, ver)
	if err != nil {
		return domain.ConfigWithNodes{}, errors.Wrap(err)
	}

	if nodes == nil {
		return domain.ConfigWithNodes{}, nil
	}

	//switch configType {
	//case api.ConfigType_pg:
	//	toSnake(nodes)
	//}

	versions, err := c.configStorage.GetVersions(ctx, configName)
	if err != nil {
		return domain.ConfigWithNodes{}, errors.Wrap(err, "error getting config by name")
	}

	return domain.ConfigWithNodes{
		Nodes:    nodes,
		Versions: versions,
	}, nil
}

func (c *CfgService) ListConfigs(
	ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error) {
	resp, err := c.configStorage.ListConfigs(ctx, req)
	if err != nil {
		return domain.ListConfigsResponse{}, errors.Wrap(err)
	}

	return resp, nil
}

func toSnake(root *evon.Node) {
	if root == nil {
		return
	}

	for _, n := range root.InnerNodes {
		toSnake(n)
	}
	root.Name = strings.ReplaceAll(root.Name, "-", "_")
}
