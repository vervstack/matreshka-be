package tests

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.vervstack.ru/matreshka/pkg/matreshka"
	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"
	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type PatchConfigSuite struct {
	suite.Suite

	ctx        context.Context
	configName string
	cfg        matreshka.AppConfig
	patchReq   *matreshka_api.PatchConfig_Request
}

func (s *PatchConfigSuite) SetupTest() {
	s.ctx = context.Background()

	s.configName = matreshka_api.ConfigType_verv.String() + "_" + getConfigNameFromTest(s.T())
	testEnv.createWithName(s.T(), s.configName)

	s.cfg = getFullConfig(s.T())
	s.cfg.Name = s.configName
	testEnv.updateConfigValues(s.T(), s.cfg)

	s.patchReq = &matreshka_api.PatchConfig_Request{
		ConfigName: s.configName,
	}
}

func (s *PatchConfigSuite) Test_PatchConfigEnvironment() {
	// Change old environment variable
	{
		s.cfg.Environment[5] = environment.MustNewVariable(
			s.cfg.Environment[5].Name, []int{50051})

		s.patchReq.Patches = append(s.patchReq.Patches,
			&matreshka_api.Patch{
				FieldName: "ENVIRONMENT_AVAILABLE-PORTS",
				Patch: &matreshka_api.Patch_UpdateValue{
					UpdateValue: fmt.Sprint(s.cfg.Environment[5].Value),
				},
			})
	}
	// Delete environment variable
	{
		s.patchReq.Patches = append(s.patchReq.Patches,
			&matreshka_api.Patch{
				FieldName: "ENVIRONMENT_CREDIT-PERCENTS-BASED-ON-YEAR-OF-BIRTH",
				Patch: &matreshka_api.Patch_Delete{
					Delete: true,
				},
			})
		s.cfg.Environment = s.cfg.Environment[:len(s.cfg.Environment)-1]
	}
	// Add new environment variable
	{
		someValue := "rand val"
		valueType := string(environment.VariableTypeStr)

		newEnvVar := environment.MustNewVariable("new_value", someValue)

		s.cfg.Environment = append(s.cfg.Environment, newEnvVar)

		s.patchReq.Patches = append(s.patchReq.Patches,
			&matreshka_api.Patch{
				FieldName: "ENVIRONMENT_NEW-VALUE",
				Patch: &matreshka_api.Patch_UpdateValue{
					UpdateValue: someValue,
				},
			},
			&matreshka_api.Patch{
				FieldName: "ENVIRONMENT_NEW-VALUE_TYPE",
				Patch: &matreshka_api.Patch_UpdateValue{
					UpdateValue: valueType,
				},
			},
		)
	}
}

func (s *PatchConfigSuite) Test_PatchConfigServers() {
	const port = 50051
	s.cfg.Servers[port].GRPC["/{GRPC}"].Gateway = "/api/v2"

	s.patchReq.Patches = append(s.patchReq.Patches,
		&matreshka_api.Patch{
			FieldName: "SERVERS_MASTER2_/{GRPC}_GATEWAY",
			Patch: &matreshka_api.Patch_UpdateValue{
				UpdateValue: s.cfg.Servers[port].GRPC["/{GRPC}"].Gateway,
			},
		})
}

func (s *PatchConfigSuite) Test_PatchConfigDataSources() {

	// Change old resource (pg) port
	{
		pg := s.cfg.DataSources[0].(*resources.Postgres)
		pg.Port = 5433

		s.patchReq.Patches = append(s.patchReq.Patches,
			&matreshka_api.Patch{
				FieldName: "DATA-SOURCES_POSTGRES_PORT",
				Patch: &matreshka_api.Patch_UpdateValue{
					UpdateValue: strconv.Itoa(int(pg.Port)),
				},
			})
	}

	// Delete old resource (telegram) data source
	// Add new telegram data source
	{
		tg := s.cfg.DataSources[2].(*resources.Telegram)
		tg.Name = "telegram_bot"
		tg.ApiKey = "jjggwwkk"

		s.patchReq.Patches = append(s.patchReq.Patches,
			&matreshka_api.Patch{
				FieldName: "DATA-SOURCES_TELEGRAM",
				Patch: &matreshka_api.Patch_Delete{
					Delete: true,
				},
			})
		s.patchReq.Patches = append(s.patchReq.Patches,
			&matreshka_api.Patch{
				FieldName: "DATA-SOURCES_TELEGRAM-BOT_API-KEY",
				Patch: &matreshka_api.Patch_UpdateValue{
					UpdateValue: tg.ApiKey,
				},
			})
	}
}

func (s *PatchConfigSuite) TearDownTest() {
	patchResp, err := testEnv.matreshkaApi.PatchConfig(s.ctx, s.patchReq)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), patchResp)

	patchedConfig := testEnv.get(s.T(), s.configName)

	sort.Slice(s.cfg.Environment, func(i, j int) bool {
		return s.cfg.Environment[i].Name < s.cfg.Environment[j].Name
	})

	require.Equal(s.T(), s.cfg.Environment, patchedConfig.Environment)
}

func Test_PatchConfig(t *testing.T) {
	t.Skip("REWORK IN PROGRESS")

	suite.Run(t, new(PatchConfigSuite))
}
