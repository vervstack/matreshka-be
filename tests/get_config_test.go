package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

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

func Test_GetConfig(t *testing.T) {
	suite.Run(t, new(GetConfigSuite))
}
