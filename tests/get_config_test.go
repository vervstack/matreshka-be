package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type GetConfigSuite struct {
	suite.Suite

	ctx context.Context
	api matreshka_api.MatreshkaBeAPIClient
}

func (s *GetConfigSuite) SetupSuite() {
	s.ctx = context.Background()
	s.api = testEnv.matreshkaApi
}

func (s *GetConfigSuite) Test_GetNotExisting() {
	for _, typePrefix := range matreshka_api.ConfigTypePrefix_name {
		req := &matreshka_api.GetConfig_Request{
			ConfigName: typePrefix + "_unexisting",
		}

		resp, err := s.api.GetConfig(s.ctx, req)

		expectedErr := status.Error(codes.NotFound,
			fmt.Sprintf("Not found\nservice with name %s_unexisting not found\n", typePrefix))
		require.ErrorIs(s.T(), err, expectedErr)
		require.Nil(s.T(), resp)
	}
}

func Test_GetConfig(t *testing.T) {
	suite.Run(t, new(GetConfigSuite))
}
