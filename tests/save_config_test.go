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
}

func (s *StoreConfigSuite) SetupSuite() {
	t := s.T()
	s.ctx = t.Context()

	app := InitAppEnvironment(t,
		WithPersistentDb(),
	)

	s.apiClient = app.matreshkaApi
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
	getReqYaml := &api.GetConfig_Request{
		ConfigName: configName,
		Format:     api.Format_yaml,
	}

	cfgRespYaml, err := s.apiClient.GetConfig(s.ctx, getReqYaml)
	s.Require().NoError(err)
	s.Require().YAMLEq(string(lokiConfigYaml), string(cfgRespYaml.Config))

	getReqEnv := &api.GetConfig_Request{
		ConfigName: configName,
		Format:     api.Format_env,
	}

	cfgRespEnv, err := s.apiClient.GetConfig(s.ctx, getReqEnv)
	s.Require().NoError(err)
	s.Require().Equal(string(lokiConfigEnv), string(cfgRespEnv.Config))
}

func (s *StoreConfigSuite) assertMatreshkaConfig(configName string) {
	getReqYaml := &api.GetConfig_Request{
		ConfigName: configName,
		Format:     api.Format_yaml,
	}

	cfgRespYaml, err := s.apiClient.GetConfig(s.ctx, getReqYaml)
	s.Require().NoError(err)
	s.Require().YAMLEq(string(vervConfigYaml), string(cfgRespYaml.Config))

	getReqEnv := &api.GetConfig_Request{
		ConfigName: configName,
		Format:     api.Format_env,
	}

	cfgRespEnv, err := s.apiClient.GetConfig(s.ctx, getReqEnv)
	s.Require().NoError(err)
	s.Require().Equal(string(vervConfigEnv), string(cfgRespEnv.Config))
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
