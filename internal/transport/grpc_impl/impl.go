package grpc_impl

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"go.vervstack.ru/matreshka/internal/config"
	"go.vervstack.ru/matreshka/internal/service"
	"go.vervstack.ru/matreshka/pkg/matreshka_be_api"
)

type Impl struct {
	version string

	configService service.ConfigService
	subService    service.SubscriberService

	matreshka_be_api.UnimplementedMatreshkaBeAPIServer
}

func NewServer(
	cfg config.Config,
	service service.Services,
) *Impl {
	return &Impl{
		version:       cfg.AppInfo.Version,
		configService: service.ConfigService(),
		subService:    service.PubSubService(),
	}
}

func (a *Impl) Register(srv grpc.ServiceRegistrar) {
	matreshka_be_api.RegisterMatreshkaBeAPIServer(srv, a)
}

func (a *Impl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (route string, handler http.Handler) {
	gwHttpMux := runtime.NewServeMux()

	err := matreshka_be_api.RegisterMatreshkaBeAPIHandlerFromEndpoint(
		ctx,
		gwHttpMux,
		endpoint,
		opts,
	)
	if err != nil {
		logrus.Errorf("error registering grpc2http handler: %s", err)
	}

	return "/api/", gwHttpMux
}
