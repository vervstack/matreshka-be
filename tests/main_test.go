package tests

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"go.vervstack.ru/matreshka/internal/transport/matreshka_api_impl"
	"go.vervstack.ru/matreshka/pkg/app"
	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type Env struct {
	matreshkaApi matreshka_api.MatreshkaBeAPIClient

	a          app.App
	HttpServer *httptest.Server
}

//go:embed config/test.config.yaml
var fullConfigBytes []byte

var testEnv Env

func TestMain(m *testing.M) {
	defer closer.Close()

	err := initApp()
	if err != nil {
		logrus.Fatal(err)
	}

	var code int
	code = m.Run()
	os.Exit(code)
}

func initApp() (err error) {
	testEnv.a, err = app.New()
	if err != nil {
		return errors.Wrap(err, "error initializing config")
	}

	_, err = testEnv.a.Sqlite.Exec(`
		DELETE 
		FROM configs 	   
	    WHERE true;
		
		DELETE 
		FROM configs_values
		WHERE true;`)
	if err != nil {
		return errors.Wrap(err, "error db clean up")
	}

	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	serv := grpc.NewServer()
	matreshka_api.RegisterMatreshkaBeAPIServer(serv, testEnv.a.Custom.GrpcImpl)
	go func() {
		if err := serv.Serve(lis); err != nil {
			logrus.Fatalf("error serving grpc server for tests %s", err)
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("[::]:"+testEnv.a.Cfg.Servers.MASTER.Port,
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrus.Fatalf("error connecting to test grpc server: %s ", err)
	}

	testEnv.matreshkaApi = matreshka_api.NewMatreshkaBeAPIClient(conn)

	ping, err := testEnv.matreshkaApi.ApiVersion(testEnv.a.Ctx, &matreshka_api.ApiVersion_Request{})
	if err != nil {
		logrus.Fatalf("error pingin test server: %s", err)
	}

	if ping == nil {
		logrus.Fatalf("error pingin test server")
	}

	testEnv.HttpServer = httptest.NewServer(testEnv.a.Custom.WebApiImpl)
	closer.Add(func() error {
		testEnv.HttpServer.Close()
		return nil
	})
	return nil
}

func (e *Env) create(t *testing.T) string {
	configName := normalizeConfigName(t.Name())
	e.createWithName(t, configName)

	return configName
}

func (e *Env) createWithName(t *testing.T, configName string) {
	createReq := &matreshka_api.CreateConfig_Request{
		ConfigName: configName,
	}
	ctx := context.Background()

	postResp, err := testEnv.matreshkaApi.CreateConfig(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, postResp)
}

func (e *Env) updateConfigValues(t *testing.T, cfg matreshka.AppConfig) {
	req := &matreshka_api.PatchConfig_Request{
		ConfigName: cfg.ModuleName(),
	}

	nodes, err := evon.MarshalEnv(&cfg)
	require.NoError(t, err)

	storage := evon.NodesToStorage(nodes)

	for k, v := range storage {
		if v.Value != nil {
			req.Patches = append(req.Patches,
				&matreshka_api.PatchConfig_Patch{
					FieldName: k,
					Patch: &matreshka_api.PatchConfig_Patch_UpdateValue{
						UpdateValue: fmt.Sprint(v.Value),
					},
				})
		}
	}

	ctx := context.Background()

	_, err = e.matreshkaApi.PatchConfig(ctx, req)
	require.NoError(t, err)
}

func (e *Env) get(t *testing.T, configName string) matreshka.AppConfig {
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

	fullConfig.Name = getServiceNameFromTest(t)

	return fullConfig
}

func getServiceNameFromTest(t *testing.T) string {
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
