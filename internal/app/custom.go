package app

import (
	"context"
	"net/http"

	"go.vervstack.ru/matreshka/internal/middleware"
	"go.vervstack.ru/matreshka/internal/service"
	"go.vervstack.ru/matreshka/internal/service/v1"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/internal/storage/sqlite"

	//"go.vervstack.ru/matreshka/internal/storage/pg"
	"go.vervstack.ru/matreshka/internal/storage/tx_manager"
	"go.vervstack.ru/matreshka/internal/transport/matreshka_api_impl"
	"go.vervstack.ru/matreshka/internal/transport/web"
	"go.vervstack.ru/matreshka/internal/transport/web_api"
	"go.vervstack.ru/matreshka/internal/web/auth"
	"go.vervstack.ru/matreshka/pkg/docs"
)

type Custom struct {
	DataProvider  storage.Data
	ConfigStorage storage.ConfigStorage

	Service service.Services

	GrpcImpl   *matreshka_api_impl.Impl
	WebApiImpl http.Handler
}

func (c *Custom) Init(a *App) (err error) {
	// Repository, Service logic, transport registration happens here

	txManager := tx_manager.New(a.Postgres)

	c.ConfigStorage = sqlite.New(a.Sqlite)
	c.Service = v1.New(c.DataProvider, c.ConfigStorage, txManager)

	c.GrpcImpl = matreshka_api_impl.NewServer(a.Cfg, c.Service)
	c.WebApiImpl = web_api.New(c.GrpcImpl)

	a.ServerMaster.AddImplementation(c.GrpcImpl)

	if a.Cfg.Environment.Pass != "" {
		a.ServerMaster.AddServerOption(auth.Interceptor(a.Cfg.Environment.Pass))
	}

	a.ServerMaster.AddServerOption(
		middleware.PanicInterceptor(),
		middleware.LogInterceptor(),
	)

	a.ServerMaster.AddHttpHandler("/web_api/", c.WebApiImpl)
	a.ServerMaster.AddHttpHandler("/", web.NewServer())
	a.ServerMaster.AddHttpHandler(docs.Swagger())

	return nil
}

func (c *Custom) Start(ctx context.Context) error {
	return nil
}

func (c *Custom) Stop() error {
	return nil
}
