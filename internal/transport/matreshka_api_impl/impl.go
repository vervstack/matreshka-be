package matreshka_api_impl

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/config"
	"go.vervstack.ru/matreshka/internal/service"
)

type Impl struct {
	version string

	configService service.ConfigService

	subService service.SubscriberService

	matreshka_api.UnimplementedMatreshkaApiServer
}

func NewServer(
	cfg config.Config,
	srv service.Services,
) *Impl {
	return &Impl{
		version: cfg.AppInfo.Version,

		configService: srv.ConfigService(),
		subService:    srv.PubSubService(),
	}
}

func (s *Impl) Register(srv grpc.ServiceRegistrar) {
	matreshka_api.RegisterMatreshkaApiServer(srv, s)
}

func (s *Impl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (route string, handler http.Handler) {
	gwHttpMux := runtime.NewServeMux()

	err := matreshka_api.RegisterMatreshkaApiHandlerFromEndpoint(
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
