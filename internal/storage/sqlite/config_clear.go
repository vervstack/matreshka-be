package sqlite

import (
	"context"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) ClearValues(ctx context.Context, req domain.ConfigName, version *string) error {
	if version == nil {
		version = toolbox.ToPtr("%%")
	}

	params := config_queries.ClearValuesParams{
		Name:    req.Name(),
		Version: *version,
	}

	err := p.querier.ClearValues(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error removing values from configs")
	}

	return nil
}
