package matreshka_api_impl

import (
	"context"

	errors "go.redsock.ru/rerrors"

	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (s *Impl) RenameConfig(ctx context.Context, req *api.RenameConfig_Request) (*api.RenameConfig_Response, error) {
	err := s.configService.Rename(ctx, req.ConfigName, req.NewName)
	if err != nil {
		return nil, errors.Wrap(err, "error renaming config")
	}

	return &api.RenameConfig_Response{NewName: req.NewName}, nil
}
