package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type CreateSuite struct {
	suite.Suite

	ctx context.Context
	api matreshka_api.MatreshkaApiClient
}

func (s *CreateSuite) SetupSuite() {
	t := s.T()
	appEnv := InitAppEnvironment(t)

	s.api = appEnv.matreshkaApi
}

func (s *CreateSuite) SetupTest() {
	s.ctx = s.T().Context()
}

func (s *CreateSuite) Test_Ok() {
	t := s.T()

	expectedResults := map[matreshka_api.ConfigType]*matreshka_api.GetConfig_Response{}

	testCases := []matreshka_api.ConfigType{
		matreshka_api.ConfigType_verv,
	}

	for _, configType := range testCases {
		t.Run(configType.String(), func(t *testing.T) {
			configName := getConfigNameFromTest(t)

			createReq := &matreshka_api.CreateConfig_Request{

				ConfigName: configName,
				ConfigType: configType,
			}

			actualResponse, err := s.api.CreateConfig(s.ctx, createReq)
			require.NoError(t, err)
			require.NotNil(t, actualResponse)

			getConfigReq := &matreshka_api.GetConfig_Request{
				ConfigName: configName,
			}
			config, err := s.api.GetConfig(s.ctx, getConfigReq)
			require.NoError(t, err)
			require.NotNil(t, actualResponse)

			expectedResponse, ok := expectedResults[createReq.ConfigType]
			if ok {
				require.Equal(t, expectedResponse, config)
			}
		})
	}
}

func (s *CreateSuite) Test_AlreadyExist() {
	t := s.T()

	configName := getConfigNameFromTest(t)

	createReq := &matreshka_api.CreateConfig_Request{

		ConfigName: configName,
		ConfigType: matreshka_api.ConfigType_verv,
	}

	actualResponse, err := s.api.CreateConfig(s.ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, actualResponse)

	actualResponse, err = s.api.CreateConfig(s.ctx, createReq)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, st.Code(), codes.AlreadyExists)

	require.Nil(t, actualResponse)
}

func (s *CreateSuite) Test_InvalidCharsInName() {
	type testCase struct {
		name         string
		expectedCode codes.Code
		message      string
	}

	testCases := map[string]testCase{
		"short": {
			name:         "12",
			expectedCode: codes.InvalidArgument,
			message:      "Validation error;Service name must be at least 3 symbols long",
		},
		"invalid_char": {
			name:         "12+a",
			expectedCode: codes.InvalidArgument,
			message:      "Validation error;Variable name contains invalid character: +",
		},
		"invalid_chars": {
			name:         "12+a)",
			expectedCode: codes.InvalidArgument,
			message:      "Validation error;Variable name contains invalid characters: ),+",
		},
	}

	for name, tc := range testCases {
		tc := tc
		s.Run(name, func() {
			req := &matreshka_api.CreateConfig_Request{
				ConfigName: tc.name,
			}
			resp, err := s.api.CreateConfig(s.ctx, req)
			s.Nil(resp)
			s.Error(err)

			grpcStatus := status.Convert(err)
			s.Equal(tc.expectedCode, grpcStatus.Code())
			s.Equal(tc.message, grpcStatus.Message())
		})
	}
}

func Test_CreateConfig(t *testing.T) {
	suite.Run(t, new(CreateSuite))
}
