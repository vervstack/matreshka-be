package evon

import (
	"context"
	"database/sql"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
)

func (c *CfgService) Delete(ctx context.Context,
	name domain.ConfigName, version string) error {

	var versionToDeleteIn *string
	if version != domain.MasterVersion {
		versionToDeleteIn = &version
	}

	err := c.txManager.Execute(func(tx *sql.Tx) error {
		err := c.configStorage.ClearValues(ctx, name, versionToDeleteIn)
		if err != nil {
			return rerrors.Wrap(err, "error deleting values from storage")
		}

		if version == domain.MasterVersion {
			err = c.configStorage.Delete(ctx, name)
			if err != nil {
				return rerrors.Wrap(err, "error deleting config from storage")
			}
		}

		return nil
	})
	if err != nil {
		return rerrors.Wrap(err)
	}

	return nil
}
