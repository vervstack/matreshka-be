package matreshka_api_impl

import (
	"context"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
)

func (s *Impl) DeleteConfig(ctx context.Context,
	req *matreshka_api.DeleteConfig_Request) (
	*matreshka_api.DeleteConfig_Response, error) {

	vers := toolbox.Coalesce(
		req.ConfigVersion,
		toolbox.ToPtr(domain.MasterVersion))

	err := s.configService.Delete(ctx, req.ConfigName, *vers)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting config")
	}

	return &matreshka_api.DeleteConfig_Response{}, nil
}
