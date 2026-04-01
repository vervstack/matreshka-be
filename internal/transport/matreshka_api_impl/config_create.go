package matreshka_api_impl

import (
	"context"

	errors "go.redsock.ru/rerrors"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
)

func (s *Impl) CreateConfig(ctx context.Context, req *api.CreateConfig_Request) (*api.CreateConfig_Response, error) {
	domainReq := domain.CreateConfigRequest{
		Name: req.GetConfigName(),
		Type: req.GetConfigType(),
	}

	err := s.configService.Create(ctx, domainReq)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &api.CreateConfig_Response{}, nil
}
