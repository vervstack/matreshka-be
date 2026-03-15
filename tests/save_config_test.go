package tests

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

var (
	//go:embed data/loki.example.yaml
	lokiConfigYaml []byte
	//go:embed data/loki.example.env
	lokiConfigEnv []byte

	//go:embed data/verv.example.yaml
	vervConfigYaml []byte
	//go:embed data/verv.example.env
	vervConfigEnv []byte
)

type StoreConfigSuite struct {
	suite.Suite

	ctx context.Context

	apiClient api.MatreshkaApiClient
	appEnv    AppEnv
}

func (s *StoreConfigSuite) SetupSuite() {
	t := s.T()
	s.ctx = t.Context()

	app := InitAppEnvironment(t)

	s.apiClient = app.matreshkaApi
	s.appEnv = app
}

func (s *StoreConfigSuite) TestStoreKeyValueConfig_Yaml() {
	configName := s.createConfig(api.ConfigType_kv)

	storeReq := &api.SaveConfig_Request{
		Format:     api.Format_yaml,
		ConfigName: configName,
		Config:     lokiConfigYaml,
	}

	_, err := s.apiClient.SaveConfig(s.ctx, storeReq)
	s.Require().NoError(err)

	s.assertLokiConfig(configName)
}

func (s *StoreConfigSuite) TestStoreKeyValueConfig_Env() {
	configName := s.createConfig(api.ConfigType_kv)

	storeReq := &api.SaveConfig_Request{
		Format:     api.Format_env,
		ConfigName: configName,
		Config:     lokiConfigEnv,
	}

	_, err := s.apiClient.SaveConfig(s.ctx, storeReq)
	s.Require().NoError(err)

	s.assertLokiConfig(configName)
}

func (s *StoreConfigSuite) TestStoreVervConfig_Yaml() {
	configName := s.createConfig(api.ConfigType_verv)

	storeReq := &api.SaveConfig_Request{
		Format:     api.Format_yaml,
		ConfigName: configName,
		Config:     vervConfigYaml,
	}

	_, err := s.apiClient.SaveConfig(s.ctx, storeReq)
	s.Require().NoError(err)

	s.assertMatreshkaConfig(configName)
}

func (s *StoreConfigSuite) assertLokiConfig(configName string) {
	s.appEnv.assertYamlConfigContent(s.T(), configName, lokiConfigYaml)
	s.appEnv.assertEnvConfigContent(s.T(), configName, lokiConfigEnv)

	s.appEnv.assertConfigBase(s.T(), configName, api.ConfigType_kv)
}

func (s *StoreConfigSuite) assertMatreshkaConfig(configName string) {
	s.appEnv.assertYamlConfigContent(s.T(), configName, vervConfigYaml)
	s.appEnv.assertEnvConfigContent(s.T(), configName, vervConfigEnv)

	s.appEnv.assertConfigBase(s.T(), configName, api.ConfigType_verv)
}

func (s *StoreConfigSuite) createConfig(configType api.ConfigType) string {
	t := s.T()
	configName := getConfigNameFromTest(t) + "_" + configType.String()

	createReq := &api.CreateConfig_Request{
		ConfigName: configName,
		ConfigType: configType,
	}
	_, err := s.apiClient.CreateConfig(s.ctx, createReq)
	require.NoError(t, err)

	return configName
}

func Test_StoreConfig(t *testing.T) {
	suite.Run(t, new(StoreConfigSuite))
}

func (e *AppEnv) assertYamlConfigContent(t *testing.T, configName string, content []byte) {
	req := &api.GetConfig_Request{
		ConfigName: configName,
		Format:     api.Format_yaml,
	}
	resp, err := e.matreshkaApi.GetConfig(t.Context(), req)
	require.NoError(t, err)
	require.YAMLEq(t, string(content), string(resp.Config))
}

func (e *AppEnv) assertEnvConfigContent(t *testing.T, configName string, content []byte) {
	req := &api.GetConfig_Request{
		ConfigName: configName,
		Format:     api.Format_env,
	}
	resp, err := e.matreshkaApi.GetConfig(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, string(content), string(resp.Config))
}

func (e *AppEnv) assertConfigBase(t *testing.T, configName string, cType api.ConfigType) {
	getReq := &api.GetConfig_Request{
		ConfigName: configName,
	}

	actualResp, err := e.matreshkaApi.GetConfig(t.Context(), getReq)
	require.NoError(t, err)

	require.NotNil(t, actualResp)
	require.NotNil(t, actualResp.BaseInfo)

	require.Equal(t, configName, actualResp.BaseInfo.Name)
	require.Equal(t, cType, actualResp.BaseInfo.ConfigType)

	require.NotEmpty(t, actualResp.BaseInfo.Id)
	require.NotEmpty(t, actualResp.BaseInfo.CreatedAt)
	require.NotEmpty(t, actualResp.BaseInfo.UpdatedAt)
}
