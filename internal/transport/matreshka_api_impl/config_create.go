package matreshka_api_impl

import (
	"context"

	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (s *Impl) CreateConfig(ctx context.Context, req *api.CreateConfig_Request) (*api.CreateConfig_Response, error) {
	cfgType, err := fromConfigType(req.GetConfigType())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	domainReq := domain.CreateConfigRequest{
		Name: req.GetConfigName(),
		Type: cfgType,
	}

	err = s.configService.Create(ctx, domainReq)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &api.CreateConfig_Response{}, nil
}
