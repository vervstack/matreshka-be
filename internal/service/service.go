package service

//go:generate minimock -i Subscriber -o ./../../tests/mocks -g -s "_mock.go"

import (
	"context"

	"go.vervstack.ru/matreshka/internal/domain"
)

type Services interface {
	ConfigService() ConfigService
	PubSubService() PubSubService
}

type ConfigService interface {
	Create(ctx context.Context, req domain.CreateConfigRequest) error

	Patch(ctx context.Context, configPatch domain.PatchConfigRequest) error

	Save(ctx context.Context, req domain.SaveConfigReq) error

	Rename(ctx context.Context, oldName, newName string) error

	List(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error)
	GetConfigWithNodes(ctx context.Context, name string, version string) (domain.ConfigWithNodes, error)

	Delete(ctx context.Context, name string, version string) error
}

type PubSubService interface {
	PublisherService
	SubscriberService
}

type PublisherService interface {
	Publish(event domain.PatchConfigRequest)
}

type SubscriberService interface {
	Subscribe(c Subscriber, serviceNames ...string)
	Unsubscribe(c Subscriber, serviceNames ...string)
	StopSubscription(c Subscriber)
}

type Subscriber interface {
	Consume(request domain.PatchConfigRequest)
	GetUpdateChan() chan domain.PatchConfigRequest
	Stop()
}
