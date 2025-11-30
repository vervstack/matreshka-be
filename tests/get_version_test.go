package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type GetVersionSuite struct {
	suite.Suite
}

func (s *GetVersionSuite) Test_GetVersion() {
	ctx := context.Background()
	resp, err := testEnv.matreshkaApi.ApiVersion(ctx, &matreshka_api.ApiVersion_Request{})

	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func Test_GetVersion(t *testing.T) {
	suite.Run(t, new(GetVersionSuite))
}
