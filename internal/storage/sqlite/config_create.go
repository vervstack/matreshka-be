package sqlite

import (
	"context"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) Create(ctx context.Context, r domain.CreateConfigRequest) (int64, error) {
	params := config_queries.CreateConfigParams{
		Name:     r.Name,
		TypeName: r.Type.String(),
	}
	cfgId, err := p.querier.CreateConfig(ctx, params)
	if err != nil {
		return 0, wrapError(err)
	}

	return cfgId, nil
}
