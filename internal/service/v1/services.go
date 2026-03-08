package v1

import (
	"go.vervstack.ru/matreshka/internal/service"
	"go.vervstack.ru/matreshka/internal/service/v1/cfg_service"
	"go.vervstack.ru/matreshka/internal/service/v1/subscription"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/internal/storage/tx_manager"
)

type Services struct {
	configService *cfg_service.CfgService

	pubSubService *subscription.PubSubService
}

func New(data storage.Data, txManager *tx_manager.TxManager) service.Services {
	pubSubService := subscription.New()

	return &Services{
		pubSubService: pubSubService,
		configService: cfg_service.New(data, txManager, pubSubService),
	}
}

func (s *Services) ConfigService() service.ConfigService {
	return s.configService
}

func (s *Services) PubSubService() service.PubSubService {
	return s.pubSubService
}
