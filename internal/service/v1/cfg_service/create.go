package cfg_service

import (
	"context"
	"database/sql"
	"fmt"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/config_templates"
	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (c *CfgService) Create(ctx context.Context, req domain.CreateConfigRequest) error {
	err := c.validator.IsConfigNameValid(req.Name)
	if err != nil {
		return errors.Wrap(err)
	}

	switch req.Type {
	case matreshka_api.ConfigType_verv:
		return c.createEmptyEvonConfig(ctx, req)
	default:
		return errors.New("currently no config types except for verv is supported")
	}
}

func (c *CfgService) createEmptyEvonConfig(ctx context.Context, req domain.CreateConfigRequest) error {
	newCfg, err := config_templates.InitNewConfig(req.Name, req.Type)
	if err != nil {
		return errors.Wrap(err)
	}

	newCfgPatch, err := c.convertConfigToPatch(newCfg)
	if err != nil {
		return errors.Wrap(err, "error converting config to patch")
	}

	err = c.txManager.Execute(func(tx *sql.Tx) error {
		dataStorage := c.configStorage.WithTx(tx)

		_, err = dataStorage.Create(ctx, req)
		if err != nil {
			return errors.Wrap(err, "error saving config")
		}

		patchReq := domain.PatchConfigRequest{
			ConfigName:    req.Name,
			Upsert:        newCfgPatch,
			ConfigVersion: domain.MasterVersion,
		}

		err = dataStorage.UpsertValues(ctx, patchReq)
		if err != nil {
			return errors.Wrap(err, "error upserting new config")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (c *CfgService) convertConfigToPatch(cfg *evon.Node) ([]domain.PatchUpdate, error) {
	newCfgNodesStore := evon.NodesToStorage(cfg)

	cfgPatch := make([]domain.PatchUpdate, 0, len(newCfgNodesStore))
	for _, node := range newCfgNodesStore {
		if node.Value != nil {
			cfgPatch = append(cfgPatch,
				domain.PatchUpdate{
					FieldName:  node.Name,
					FieldValue: fmt.Sprint(node.Value),
				})
		}
	}

	return cfgPatch, nil
}
