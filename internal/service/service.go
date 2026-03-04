package service

//go:generate minimock -i Subscriber -o ./../../tests/mocks -g -s "_mock.go"

import (
	"context"

	"go.vervstack.ru/matreshka/internal/domain"
)

type Services interface {
	ConfigService() ConfigService
	// EvonService deprecated
	EvonService() EvonConfigService
	PubSubService() PubSubService
}

type ConfigService interface {
	Create(ctx context.Context, req domain.CreateConfigRequest) error
}

type EvonConfigService interface {
	Create(ctx context.Context, name domain.ConfigName) (domain.AboutConfig, error)

	Replace(ctx context.Context, req domain.ReplaceConfigReq) error
	Patch(ctx context.Context, configPatch domain.PatchConfigRequest) error

	Rename(ctx context.Context, oldName, newName domain.ConfigName) error

	GetConfigWithNodes(ctx context.Context, name domain.ConfigName, version string) (domain.ConfigWithNodes, error)
	ListConfigs(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error)
	Delete(ctx context.Context, name domain.ConfigName, version string) error
}

type PlainConfigService interface {
	Save(cfg string, version string) error
	Get(cfg string, version string) error
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
