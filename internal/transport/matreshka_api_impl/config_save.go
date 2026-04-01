package matreshka_api_impl

import (
	"context"

	"go.redsock.ru/rerrors"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
)

func (s *Impl) SaveConfig(ctx context.Context, req *api.SaveConfig_Request) (
	*api.SaveConfig_Response, error) {

	r := domain.SaveConfigReq{
		ConfigName: req.ConfigName,
		Version:    req.Version,
		Format:     req.Format,
		Content:    req.Config,
	}

	err := s.configService.Save(ctx, r)
	if err != nil {
		return nil, rerrors.Wrap(err, "error creating config")
	}

	return &api.SaveConfig_Response{}, nil
}
