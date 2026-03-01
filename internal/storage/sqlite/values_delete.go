package sqlite

import (
	"context"

	"database/sql"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) DeleteValues(ctx context.Context, req domain.PatchConfigRequest) error {
	if len(req.Delete) == 0 {
		return nil
	}

	cfgId, err := p.getIdByName(ctx, req.ConfigName.Name())
	if err != nil {
		return rerrors.Wrap(err)
	}

	for _, valueName := range req.Delete {
		err = p.querier.DeleteValues(ctx, config_queries.DeleteValuesParams{
			ConfigID: sql.NullInt64{
				Int64: cfgId,
				Valid: true,
			},
			Key:     valueName,
			Version: req.ConfigVersion,
		})
		if err != nil {
			return rerrors.Wrap(err, "error deleting value from db: "+valueName)
		}
	}

	return nil
}
