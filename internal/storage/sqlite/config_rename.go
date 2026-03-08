package sqlite

import (
	"context"
	"database/sql"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) Rename(ctx context.Context, oldName, newName string) error {
	err := p.querier.RenameConfig(ctx, config_queries.RenameConfigParams{
		NewName: newName,
		OldName: oldName,
	})
	if err != nil {
		return rerrors.Wrap(err, "error executing config rename sql")
	}

	return nil
}

func (p *Provider) RenameValues(ctx context.Context, req domain.PatchConfigRequest) error {
	if len(req.RenameTo) == 0 {
		return nil
	}

	cfgId, err := p.getIdByName(ctx, req.ConfigName)
	if err != nil {
		return rerrors.Wrap(err)
	}

	for _, b := range req.RenameTo {
		err := p.querier.RenameValues(ctx, config_queries.RenameValuesParams{
			NewKey: b.NewName,
			ConfigID: sql.NullInt64{
				Int64: cfgId,
				Valid: true,
			},
			OldKey:  b.OldName,
			Version: req.ConfigVersion,
		})
		if err != nil {
			return rerrors.Wrap(err, "error upserting config")
		}
	}
	return nil
}
