package cfg_service

import (
	"context"
	"database/sql"
	"time"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service/user_errors"
	"go.vervstack.ru/matreshka/internal/storage"
)

func (c *CfgService) Patch(ctx context.Context, req domain.PatchConfigRequest) error {
	err := c.txManager.Execute(func(tx *sql.Tx) error {
		dataStorage := c.configStorage.WithTx(tx)

		originalConfig, err := c.getEvonConfigNodes(ctx, dataStorage, req)
		if err != nil {
			return rerrors.Wrap(err, "error getting config")
		}

		err = c.validatePatch(originalConfig, &req)
		if err != nil {
			return rerrors.Wrap(err, "error validating patch")
		}

		err = c.patch(ctx, dataStorage, req)
		if err != nil {
			return rerrors.Wrap(err, "error patching config")
		}

		return nil
	})
	if err != nil {
		return rerrors.Wrap(err)
	}

	go c.pubService.Publish(req)

	return nil
}

func (c *CfgService) getEvonConfigNodes(ctx context.Context, dataStorage storage.Data, req domain.PatchConfigRequest) (*evon.Node, error) {
	ver := toolbox.Coalesce(req.ConfigVersion, domain.MasterVersion)

	cfgNodes, err := dataStorage.GetConfigNodes(ctx, req.ConfigName, ver)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting nodes")
	}

	if cfgNodes == nil {
		return nil, rerrors.Wrap(user_errors.ErrNotFound)
	}

	return cfgNodes, nil
}

func (c *CfgService) validatePatch(originalConfig *evon.Node, patch *domain.PatchConfigRequest) error {
	evonStorage := evon.NodesToStorage(originalConfig)

	err := c.validator.AsEvon(evonStorage, patch)
	if err != nil {
		// TODO Replace onto rerrors.UserError with documentation link here
		return rerrors.Wrap(err, "failed to validate EVON format")
	}

	switch patch.ConfigType {
	case api.ConfigType_verv:
		validationRes := c.validator.AsVerv(originalConfig, patch)
		if len(validationRes.Invalid) != 0 {
			return rerrors.New("error during patch validation: %v", validationRes.Invalid)
		}
	}

	return nil
}

func (c *CfgService) patch(ctx context.Context, configStorage storage.Data, patch domain.PatchConfigRequest) error {
	err := configStorage.DeleteValues(ctx, patch)
	if err != nil {
		return rerrors.Wrap(err, "error deleting values")
	}

	err = configStorage.UpsertValues(ctx, patch)
	if err != nil {
		return rerrors.Wrap(err, "error patching config in data storage")
	}

	err = configStorage.SetUpdatedAt(ctx, patch.ConfigName, time.Now())
	if err != nil {
		return rerrors.Wrap(err, "error updating time")
	}

	err = configStorage.RenameValues(ctx, patch)
	if err != nil {
		return rerrors.Wrap(err, "error renaming config")
	}

	return nil

}
