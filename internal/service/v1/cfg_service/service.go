package cfg_service

import (
	"go.vervstack.ru/matreshka/internal/service"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/internal/storage/tx_manager"
)

type CfgService struct {
	configStorage storage.Data
	txManager     *tx_manager.TxManager

	validator  Validator
	pubService service.PublisherService
}

func New(data storage.Data, txManager *tx_manager.TxManager, pubService service.PublisherService) *CfgService {
	return &CfgService{
		configStorage: data,
		txManager:     txManager,

		validator:  newValidator(),
		pubService: pubService,
	}
}
