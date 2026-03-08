package cfg_service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service/user_errors"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (c *CfgService) Create(ctx context.Context, req domain.CreateConfigRequest) error {
	err := c.txManager.Execute(
		func(tx *sql.Tx) (err error) {
			err = c.createEvonConfigTx(ctx, c.configStorage.WithTx(tx), req.Name, req.Type)
			if err != nil {
				return errors.Wrap(err, "error creating config")
			}

			return nil
		})
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (c *CfgService) createEvonConfigTx(ctx context.Context, dataStorage storage.Data, name string, cfgType matreshka_api.ConfigType) error {
	err := c.validator.IsConfigNameValid(name)
	if err != nil {
		return errors.Wrap(err)
	}

	nodes, err := dataStorage.GetConfigNodes(ctx, name, domain.MasterVersion)
	if err != nil {
		return errors.Wrap(err, "error reading config from storage")
	}

	if nodes != nil {
		return errors.Wrap(user_errors.ErrAlreadyExists)
	}

	_, err = dataStorage.Create(ctx, name)
	if err != nil {
		return errors.Wrap(err, "error saving config")
	}

	newCfg, err := initNewConfig(name, cfgType)
	if err != nil {
		return errors.Wrap(err)
	}

	newCfgPatch, err := c.convertConfigToPatch(newCfg)
	if err != nil {
		return errors.Wrap(err, "error converting config to patch")
	}

	patchReq := domain.PatchConfigRequest{
		ConfigName:    name,
		Upsert:        newCfgPatch,
		ConfigVersion: domain.MasterVersion,
	}

	err = dataStorage.UpsertValues(ctx, patchReq)
	if err != nil {
		return errors.Wrap(err, "error upserting new config")
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

func initNewConfig(configName string, configType matreshka_api.ConfigType) (*evon.Node, error) {
	switch configType {
	case matreshka_api.ConfigType_verv:
		newCfg := matreshka.NewEmptyConfig()
		newCfg.AppInfo = matreshka.AppInfo{
			Name:            configName,
			Version:         "v0.0.1",
			StartupDuration: time.Second * 5,
		}
		nodes, err := evon.MarshalEnv(&newCfg)
		if err != nil {
			return nil, errors.Wrap(err, "error marshalling config")
		}
		return nodes, nil
	case matreshka_api.ConfigType_pg:
		return &evon.Node{
			InnerNodes: []*evon.Node{
				{
					Name:  "POSTGRES-USER",
					Value: "postgres",
				},
				{
					Name:       "POSTGRES-PASSWORD",
					Value:      "123",
					InnerNodes: nil,
				},
				{
					Name:  "POSTGRES-DB",
					Value: "postgres",
				},
				{
					Name:  "POSTGRES-HOST-AUTH-METHOD",
					Value: "trust",
				},
			},
		}, nil
	default:
		return &evon.Node{}, nil
	}
}
