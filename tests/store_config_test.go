package tests

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

//go:embed data/loki.example.yaml
var lokiConfig []byte

type StoreConfigSuite struct {
	suite.Suite

	ctx        context.Context
	configName string
	apiClient  api.MatreshkaApiClient
}

func (s *StoreConfigSuite) SetupTest() {
	s.ctx = context.Background()
	s.apiClient = testEnv.matreshkaApi

	s.configName = testEnv.create(s.T())
}

func (s *StoreConfigSuite) TestStoreViaGrpc() {
	storeReq := &api.StoreConfig_Request{
		Format:     api.Format_yaml,
		ConfigName: s.configName,
		Config:     lokiConfig,
	}

	_, err := s.apiClient.StoreConfig(s.ctx, storeReq)
	s.Require().NoError(err)

	getReq := &api.GetConfig_Request{
		ConfigName: s.configName,
		Format:     api.Format_yaml,
	}

	cfg, err := s.apiClient.GetConfig(s.ctx, getReq)
	s.Require().NoError(err)
	s.Require().YAMLEq(string(lokiConfig), string(cfg.Config))
}

// TODO can be much better (test out different formats, configs, values and etc)
func (s *StoreConfigSuite) TestStoreViaHttp() {

	baseURL := "http://" + testEnv.HttpServer.Listener.Addr().String()

	uploadURL := baseURL + "/web_api/upload/" + s.configName

	uploadResp, err := http.Post(
		uploadURL,
		"application/json",
		bytes.NewBuffer(lokiConfig))
	require.NoError(s.T(), err)
	s.Require().Equal(http.StatusOK, uploadResp.StatusCode)

	downloadURL := baseURL + "/web_api/download/" + s.configName

	downloadResp, err := http.Get(downloadURL)
	s.Require().NoError(err)

	downloadedConfig, err := io.ReadAll(downloadResp.Body)
	s.Require().NoError(err)

	s.Require().YAMLEq(string(lokiConfig), string(downloadedConfig))
}

func Test_StoreConfig(t *testing.T) {
	t.Skip("REWORK IN PROGRESS")

	suite.Run(t, new(StoreConfigSuite))
}
