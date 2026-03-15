package matreshka_api_impl

import (
	"context"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
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
