package config

import (
	"context"

	"go.vervstack.ru/matreshka/internal/domain"
)

func (s *Service) List(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error) {
	list, err := s.configStorage.ListConfigs(ctx, req)
	if err != nil {
		return domain.ListConfigsResponse{}, err
	}

	return list, nil
}
