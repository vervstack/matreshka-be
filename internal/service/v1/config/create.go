package config

import (
	"context"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
)

func (s *Service) Create(ctx context.Context, req domain.CreateConfigRequest) error {
	createParams := config_queries.CreateConfigParams{
		Name: req.Name,
		Type: req.Type,
	}

	_, err := s.configStorage.CreateConfig(ctx, createParams)
	if err != nil {
		return rerrors.Wrap(err, "create config")
	}

	return nil
}
