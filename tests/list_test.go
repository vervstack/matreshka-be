package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.redsock.ru/toolbox"
	"google.golang.org/protobuf/proto"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type ListSuite struct {
	suite.Suite

	ctx         context.Context
	serviceName string
	start       time.Time

	req      *matreshka_api.ListConfigs_Request
	expected *matreshka_api.ListConfigs_Response
}

func (s *ListSuite) SetupTest() {
	s.ctx = context.Background()

	s.start = time.Now().Add(-time.Minute).UTC()
	s.serviceName = matreshka_api.ConfigTypePrefix_kv.String() + "_" + getServiceNameFromTest(s.T())
}

func (s *ListSuite) Test_ListOneServiceWithOneVersion() {
	testEnv.createWithName(s.T(), s.serviceName)

	s.req = &matreshka_api.ListConfigs_Request{
		SearchPattern: s.serviceName,
	}

	s.expected = &matreshka_api.ListConfigs_Response{
		Configs: []*matreshka_api.Config{{
			Name:     s.serviceName,
			Version:  domain.MasterVersion,
			Versions: []string{domain.MasterVersion},
		}},
		TotalRecords: 1,
	}
}

func (s *ListSuite) Test_ListOneServiceWithTwoVersion() {
	testEnv.createWithName(s.T(), s.serviceName)

	patchReq := &matreshka_api.PatchConfig_Request{
		ConfigName: s.serviceName,
		Patches: []*matreshka_api.PatchConfig_Patch{
			{
				FieldName: "ENVIRONMENT_IS-CRON-ACTIVE",
				Patch: &matreshka_api.PatchConfig_Patch_UpdateValue{
					UpdateValue: "true",
				},
			},
			{
				FieldName: "ENVIRONMENT_IS-CRON-ACTIVE_TYPE",
				Patch: &matreshka_api.PatchConfig_Patch_UpdateValue{
					UpdateValue: "bool",
				},
			},
		},
		Version: toolbox.ToPtr("VERV-137"),
	}

	_, err := testEnv.matreshkaApi.PatchConfig(s.ctx, patchReq)
	require.NoError(s.T(), err)

	s.req = &matreshka_api.ListConfigs_Request{
		SearchPattern: s.serviceName,
	}

	s.expected = &matreshka_api.ListConfigs_Response{
		Configs: []*matreshka_api.Config{{
			Name:     s.serviceName,
			Version:  domain.MasterVersion,
			Versions: []string{domain.MasterVersion, "VERV-137"},
		}},
		TotalRecords: 1,
	}
}

func (s *ListSuite) TearDownTest() {
	resp, err := testEnv.matreshkaApi.ListConfigs(s.ctx, s.req)
	require.NoError(s.T(), err)

	tm := time.Unix(resp.Configs[0].UpdatedAtUtcTimestamp, 0).UTC()
	require.WithinRange(s.T(), tm, s.start, time.Now().UTC())
	resp.Configs[0].UpdatedAtUtcTimestamp = 0

	if !proto.Equal(s.expected, resp) {
		require.Equal(s.T(), s.expected, resp)
	}
}
func Test_List(t *testing.T) {
	suite.Run(t, new(ListSuite))
}
