package config

import (
	"go.vervstack.ru/matreshka/internal/service"
	"go.vervstack.ru/matreshka/internal/storage"
)

type Service struct {
	configStorage storage.ConfigStorage
}

func New(configStorage storage.ConfigStorage) service.BinaryConfigService {
	return &Service{
		configStorage: configStorage,
	}
}
