package matreshka_api_impl

import (
	"context"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
)

func (s *Impl) Version(_ context.Context, _ *api.Version_Request) (*api.Version_Response, error) {
	resp := &api.Version_Response{
		Version: s.version,
	}

	return resp, nil
}
