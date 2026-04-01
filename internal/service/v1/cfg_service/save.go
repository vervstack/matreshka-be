package cfg_service

import (
	"context"
	"database/sql"
	"fmt"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service/user_errors"
	"go.vervstack.ru/matreshka/pkg/parsers"
)

func (c *CfgService) Save(ctx context.Context, req domain.SaveConfigReq) error {
	cfg, err := c.configStorage.GetConfigByName(ctx, req.ConfigName)
	if err != nil {
		return rerrors.Wrap(err, "error getting config by name")
	}

	switch cfg.Type {
	case api.ConfigType_kv, api.ConfigType_verv:
		var nodes *evon.Node
		nodes, err = parsers.ParseEvon(req.Content, req.Format, cfg.Type)
		if err != nil {
			return rerrors.Wrap(err, "error parsing evon")
		}

		err = c.saveEvon(ctx, req, nodes)
		if err != nil {
			return rerrors.Wrap(err, "error saving evon")
		}

		return nil
	default:
		//switch ConfigType {
		//case api.ConfigType_verv:
		//	replaceReq.Config, err = fromVervYamlToEvon(req.Config)
		//default:
		//	replaceReq.Config, err = fromPlainYamlToEvon(req.Config)
		//}
		return rerrors.Wrap(user_errors.ErrNotImplemented, "unsupported format", req.Format.String())
	}
}

func (c *CfgService) saveEvon(ctx context.Context, req domain.SaveConfigReq, nodes *evon.Node) (err error) {
	upsert := domain.PatchConfigRequest{
		ConfigName:    req.ConfigName,
		ConfigVersion: *toolbox.Coalesce(req.Version, toolbox.ToPtr(domain.MasterVersion)),
	}

	ns := evon.NodesToStorage(nodes)

	for k, n := range ns {
		if n.Value != nil {
			patch := domain.PatchUpdate{
				FieldName:  k,
				FieldValue: fmt.Sprint(n.Value),
			}

			upsert.Upsert = append(upsert.Upsert, patch)
		}
	}

	err = c.txManager.Execute(
		func(tx *sql.Tx) error {
			configStorage := c.configStorage.WithTx(tx)

			exErr := configStorage.ClearValues(ctx, req.ConfigName, req.Version)
			if exErr != nil {
				return rerrors.Wrap(err, "error clearing old values")
			}

			exErr = configStorage.UpsertValues(ctx, upsert)
			if exErr != nil {
				return rerrors.Wrap(err, "error upserting config values")
			}

			return nil
		})

	if err != nil {
		return rerrors.Wrap(err, "txManager returned an error")
	}

	return nil
}
