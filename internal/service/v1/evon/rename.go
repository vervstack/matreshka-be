package evon

import (
	"context"

	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
)

func (c *CfgService) Rename(ctx context.Context, oldName, newName domain.ConfigName) error {
	err := c.validator.IsConfigNameValid(newName)
	if err != nil {
		return errors.Wrap(err)
	}

	err = c.configStorage.Rename(ctx, oldName.Name(), newName.Name())
	if err != nil {
		return errors.Wrap(err, "error during rename operation")
	}

	return nil
}
