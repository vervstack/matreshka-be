package tests

import (
	"context"
	"testing"

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
	s.api = testEnv.matreshkaApi
}

func (s *CreateSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *CreateSuite) Test_InvalidName() {
	type testCase struct {
		name         string
		expectedCode codes.Code
		message      string
	}

	testCases := map[string]testCase{
		"short": {
			name:         "12",
			expectedCode: codes.InvalidArgument,
			message:      "Validation error\nService name must be at least 3 symbols long\n\nerror creating config\n\n",
		},
		"invalid_char": {
			name:         "12+a",
			expectedCode: codes.InvalidArgument,
			message:      "Validation error\nName contains invalid character: +\n\nerror creating config\n\n",
		},
		"invalid_chars": {
			name:         "12+a)",
			expectedCode: codes.InvalidArgument,
			message:      "Validation error\nName contains invalid characters: +,)\n\nerror creating config\n\n",
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
