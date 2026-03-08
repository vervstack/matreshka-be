package v1

import (
	"go.vervstack.ru/matreshka/internal/service"
	"go.vervstack.ru/matreshka/internal/service/v1/config"
	"go.vervstack.ru/matreshka/internal/service/v1/evon"
	"go.vervstack.ru/matreshka/internal/service/v1/subscription"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/internal/storage/tx_manager"
)

type Services struct {
	configService service.BinaryConfigService

	evonService   *evon.CfgService
	pubSubService *subscription.PubSubService
}

func New(data storage.Data, configStorage storage.ConfigStorage, txManager *tx_manager.TxManager) service.Services {
	pubSubService := subscription.New()

	return &Services{
		configService: config.New(configStorage),

		evonService:   evon.New(data, txManager, pubSubService),
		pubSubService: pubSubService,
	}
}

func (s *Services) EvonService() service.EvonConfigService {
	return s.evonService
}

func (s *Services) PubSubService() service.PubSubService {
	return s.pubSubService
}
func (s *Services) ConfigService() service.BinaryConfigService {
	return s.configService
}
