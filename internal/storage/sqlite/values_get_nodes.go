package sqlite

import (
	"context"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) GetConfigNodes(ctx context.Context, serviceName string, version string) (*evon.Node, error) {
	rows, err := p.querier.GetConfigNodes(ctx, config_queries.GetConfigNodesParams{
		MasterVersion: domain.MasterVersion,
		Version:       version,
		Name:          serviceName,
	})
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting config values")
	}

	if len(rows) == 0 {
		return nil, nil
	}

	ns := evon.NodesToStorage(nil)
	for _, row := range rows {
		ns.AddNode(&evon.Node{
			Name:  row.Key,
			Value: row.Value,
		})
	}
	return ns[""], nil
}
