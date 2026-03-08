package cfg_service

import (
	"context"

	errors "go.redsock.ru/rerrors"
)

func (c *CfgService) Rename(ctx context.Context, oldName, newName string) error {
	err := c.validator.IsConfigNameValid(newName)
	if err != nil {
		return errors.Wrap(err)
	}

	err = c.configStorage.Rename(ctx, oldName, newName)
	if err != nil {
		return errors.Wrap(err, "error during rename operation")
	}

	return nil
}
