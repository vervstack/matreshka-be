package sqlite

import (
	"context"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) GetConfigNodes(ctx context.Context, serviceName string, version string) (*evon.Node, error) {
	params := config_queries.GetConfigNodesParams{
		MasterVersion: domain.MasterVersion,
		Version:       version,
		Name:          serviceName,
	}

	rows, err := p.querier.GetConfigNodes(ctx, params)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting config values")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	ns := evon.NodesToStorage(nil)
	for _, row := range rows {
		n := &evon.Node{
			Name:  row.Key,
			Value: row.Value,
		}

		ns.AddNode(n)
	}
	return ns[""], nil
}

func (p *Provider) GetConfigRawContent(ctx context.Context, name string, version string) ([]byte, error) {
	params := config_queries.GetRawContentParams{
		Version: version,
		Name:    name,
	}
	content, err := p.querier.GetRawContent(ctx, params)
	if err != nil {
		return nil, wrapError(err)
	}

	return content, nil
}
