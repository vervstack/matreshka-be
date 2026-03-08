package cfg_service

import (
	"context"

	"go.vervstack.ru/matreshka/internal/domain"
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

func (c *CfgService) List(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func New(data storage.Data, txManager *tx_manager.TxManager, pubService service.PublisherService) *CfgService {
	return &CfgService{
		configStorage: data,
		txManager:     txManager,

		validator:  newValidator(),
		pubService: pubService,
	}
}
