package tests

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	_ "modernc.org/sqlite"

	"go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/config"
	"go.vervstack.ru/matreshka/internal/transport/matreshka_api_impl"
	"go.vervstack.ru/matreshka/pkg/app"
	"go.vervstack.ru/matreshka/pkg/matreshka"
)

type AppEnv struct {
	matreshkaApi matreshka_api.MatreshkaApiClient

	app        app.App
	HttpServer *httptest.Server
}

//go:embed config/test.config.yaml
var fullConfigBytes []byte

var testEnv AppEnv

func TestMain(m *testing.M) {
	//defer closer.Close()

	//err := initApp()
	//if err != nil {
	//	logrus.Fatal(err)
	//}

	var code int
	code = m.Run()
	os.Exit(code)
}

func InitAppEnvironment(t *testing.T, opts ...opt) AppEnv {
	ctx := t.Context()

	stopFunc := func() {}

	o := appEnvOpts{}

	for _, opt := range opts {
		opt(&o)
	}

	inMemoryListener := initInMemoryListener()
	sqliteDb := initInMemorySqldb(t, o)

	applicationBase := app.App{
		Ctx:    ctx,
		Stop:   stopFunc,
		Cfg:    config.Config{},
		Sqlite: sqliteDb,
		MASTER: inMemoryListener,
		Custom: app.CustomApp{},
	}

	err := applicationBase.Custom.Init(&applicationBase)
	require.NoError(t, err)

	httpServer := httptest.NewServer(applicationBase.Custom.WebApiImpl)
	t.Cleanup(func() {
		httpServer.Close()
	})

	matreshkaClient := initClient(t, inMemoryListener, applicationBase.Custom.GrpcImpl)

	return AppEnv{
		matreshkaApi: matreshkaClient,
		HttpServer:   httpServer,
		app:          applicationBase,
	}
}

func (e *AppEnv) purge(t *testing.T) {
	_, err := testEnv.app.Sqlite.Exec(`
		DELETE 
		FROM configs 	   
	    WHERE true;
		
		DELETE 
		FROM configs_values
		WHERE true;`)
	require.NoError(t, err)
}

func (e *AppEnv) create(t *testing.T) string {
	configName := normalizeConfigName(t.Name())
	e.createWithName(t, configName)

	return configName
}

func (e *AppEnv) createWithName(t *testing.T, configName string) {
	createReq := &matreshka_api.CreateConfig_Request{
		ConfigName: configName,
	}
	ctx := context.Background()

	postResp, err := testEnv.matreshkaApi.CreateConfig(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, postResp)
}

func (e *AppEnv) createEmptyService(t *testing.T, name string, configType matreshka_api.ConfigType) {
	createReq := &matreshka_api.CreateConfig_Request{
		ConfigName: name,
		ConfigType: configType,
	}
	ctx := t.Context()

	postResp, err := e.matreshkaApi.CreateConfig(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, postResp)
}

func (e *AppEnv) updateConfigValues(t *testing.T, cfg matreshka.AppConfig) {
	req := &matreshka_api.PatchConfig_Request{
		ConfigName: cfg.ModuleName(),
	}

	nodes, err := evon.MarshalEnv(&cfg)
	require.NoError(t, err)

	storage := evon.NodesToStorage(nodes)

	for k, v := range storage {
		if v.Value != nil {
			req.Patches = append(req.Patches,
				&matreshka_api.Patch{
					FieldName: k,
					Patch: &matreshka_api.Patch_UpdateValue{
						UpdateValue: fmt.Sprint(v.Value),
					},
				})
		}
	}

	ctx := context.Background()

	_, err = e.matreshkaApi.PatchConfig(ctx, req)
	require.NoError(t, err)
}

func (e *AppEnv) get(t *testing.T, configName string) matreshka.AppConfig {
	ctx := context.Background()
	getReq := &matreshka_api.GetConfig_Request{
		ConfigName: configName,
	}
	getResp, err := testEnv.matreshkaApi.GetConfig(ctx, getReq)
	require.NoError(t, err)

	readConfig := matreshka.NewEmptyConfig()
	err = readConfig.Unmarshal(getResp.Config)
	require.NoError(t, err)

	return readConfig
}

func getFullConfig(t *testing.T) matreshka.AppConfig {
	fullConfig := matreshka.NewEmptyConfig()
	err := fullConfig.Unmarshal(fullConfigBytes)
	if err != nil {
		t.Fatal(errors.Wrap(err, "error during unmarshalling full config"))
	}

	fullConfig.Name = getConfigNameFromTest(t)

	return fullConfig
}

func getConfigNameFromTest(t *testing.T) string {
	return strings.ReplaceAll(t.Name(), "/", "_")
}

func normalizeConfigName(configName string) string {
	configName = strings.ReplaceAll(configName, "/", "__")

	pref, _ := matreshka_api_impl.ParseConfigName(configName)

	if pref == nil {
		configName = matreshka_api.ConfigType_kv.String() + "_" + configName
	}

	return configName
}

func initInMemoryListener() *bufconn.Listener {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	return lis
}

func initInMemorySqldb(t *testing.T, o appEnvOpts) *sql.DB {
	const dialect = "sqlite"

	var db *sql.DB
	var err error
	if o.withPersistentDb {
		dsFile := "./" + getConfigNameFromTest(t) + ".db"
		err = os.Remove(dsFile)
		if !errors.Is(err, os.ErrNotExist) {
			require.NoError(t, err)
		}

		db, err = sql.Open("sqlite", dsFile)
	} else {
		db, err = sql.Open("sqlite", ":memory:")
	}

	require.NoError(t, err)
	t.Cleanup(func() {
		err = db.Close()
		require.NoError(t, err)
	})

	goose.SetLogger(testlogger{t})

	err = goose.SetDialect(dialect)
	require.NoError(t, err)

	err = goose.Up(db, "./../migrations")
	require.NoError(t, err)

	return db
}

func initClient(t *testing.T, lis *bufconn.Listener, grpcImpl matreshka_api.MatreshkaApiServer) matreshka_api.MatreshkaApiClient {
	serv := grpc.NewServer()
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	matreshka_api.RegisterMatreshkaApiServer(serv, grpcImpl)
	go func() {
		err := serv.Serve(lis)
		require.NoError(t, err)
	}()

	conn, err := grpc.NewClient("[::]:50051",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	matreshkaClient := matreshka_api.NewMatreshkaApiClient(conn)
	ctx := t.Context()

	ping, err := matreshkaClient.Version(ctx, &matreshka_api.Version_Request{})
	require.NoError(t, err)
	require.NotNil(t, ping)

	return matreshkaClient
}

type testlogger struct {
	t *testing.T
}

func (t testlogger) Printf(format string, v ...interface{}) {

}

func (t testlogger) Fatalf(format string, v ...interface{}) {
	t.t.Fatalf(format, v...)
}

type appEnvOpts struct {
	withPersistentDb bool
}

type opt func(*appEnvOpts)

func WithPersistentDb() opt {
	return func(opts *appEnvOpts) {
		opts.withPersistentDb = true
	}
}
