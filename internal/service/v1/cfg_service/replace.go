package cfg_service

import (
	"context"
	"database/sql"
	"fmt"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
)

func (c *CfgService) Replace(ctx context.Context, req domain.ReplaceConfigReq) error {
	ns := evon.NodeStorage{}
	ns.AddNode(req.Config)

	upsert := domain.PatchConfigRequest{
		ConfigName:    req.Name,
		ConfigVersion: req.Version,
	}

	for k, n := range ns {
		if n.Value != nil {
			patch := domain.PatchUpdate{
				FieldName:  k,
				FieldValue: fmt.Sprint(n.Value),
			}

			upsert.Upsert = append(upsert.Upsert, patch)
		}
	}

	upsert.Upsert = append(upsert.Upsert)

	err := c.txManager.Execute(func(tx *sql.Tx) error {
		configStorage := c.configStorage.WithTx(tx)

		err := configStorage.ClearValues(ctx, req.Name, &req.Version)
		if err != nil {
			return rerrors.Wrap(err, "error clearing old values")
		}

		err = configStorage.UpsertValues(ctx, upsert)
		if err != nil {
			return rerrors.Wrap(err, "")
		}

		return nil
	})

	if err != nil {
		return rerrors.Wrap(err, "")
	}

	return nil
}
