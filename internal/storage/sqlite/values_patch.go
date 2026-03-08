package sqlite

import (
	"context"
	"database/sql"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) UpsertValues(ctx context.Context, req domain.PatchConfigRequest) error {
	if len(req.Upsert) == 0 {
		return nil
	}

	cfgId, err := p.getIdByName(ctx, req.ConfigName)
	if err != nil {
		return rerrors.Wrap(err)
	}

	for _, b := range req.Upsert {
		err := p.querier.UpsertValues(ctx, config_queries.UpsertValuesParams{
			ConfigID: sql.NullInt64{
				Int64: cfgId,
				Valid: true,
			},
			Key:     b.FieldName,
			Value:   b.FieldValue,
			Version: req.ConfigVersion,
		})
		if err != nil {
			return rerrors.Wrap(err, "error upserting config")
		}
	}

	return nil
}

func (p *Provider) getIdByName(ctx context.Context, name string) (cfgId int64, err error) {
	cfgId, err = p.querier.GetIdByName(ctx, name)
	if err != nil {
		return 0, rerrors.Wrap(err, "error getting config id by name")
	}

	return cfgId, nil
}
