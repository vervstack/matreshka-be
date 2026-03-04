package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type SubscriptionSuite struct {
	suite.Suite

	ctx        context.Context
	configName string
	apiClient  api.MatreshkaBeAPIClient
}

func (s *SubscriptionSuite) SetupTest() {
	s.ctx = context.Background()

	s.apiClient = testEnv.matreshkaApi

	s.configName = api.ConfigType_kv.String() + "_" + getServiceNameFromTest(s.T())
	testEnv.createWithName(s.T(), s.configName)
}

func (s *SubscriptionSuite) TestSubscribeOnChanges() {
	stream, err := s.apiClient.SubscribeOnChanges(s.ctx)
	require.NoError(s.T(), err)

	// Subscribe onto changes
	{
		subscribeRequest := &api.SubscribeOnChanges_Request{
			SubscribeConfigNames: []string{s.configName},
		}
		err = stream.Send(subscribeRequest)
		require.NoError(s.T(), err)
	}

	patchReq := &api.PatchConfig_Request{
		ConfigName: s.configName,
		Patches: []*api.PatchConfig_Patch{
			{
				FieldName: "ENVIRONMENT_SOME-VARIABLE_TYPE",
				Patch: &api.PatchConfig_Patch_UpdateValue{
					UpdateValue: "string",
				},
			},
			{
				FieldName: "ENVIRONMENT_SOME-VARIABLE",
				Patch: &api.PatchConfig_Patch_UpdateValue{
					UpdateValue: "123",
				},
			},
		},
	}

	expectedUpdates := &api.SubscribeOnChanges_Response{
		ConfigName: s.configName,
		Timestamp:  uint32(time.Now().UTC().Unix()),
		Patches:    patchReq.Patches,
	}

	doneC := make(chan struct{})

	go func() {
		defer close(doneC)

		actualUpdates, err := stream.Recv()
		require.NoError(s.T(), err)

		require.GreaterOrEqual(s.T(), actualUpdates.Timestamp, expectedUpdates.Timestamp,
			"Expected time of update to be greater that time before calling update")

		// Equalize it in order to pass next assertion
		actualUpdates.Timestamp = expectedUpdates.Timestamp

		if !proto.Equal(actualUpdates, expectedUpdates) {
			require.Equal(s.T(), actualUpdates, expectedUpdates)
		}
	}()
	// Perform change in configuration

	_, err = s.apiClient.PatchConfig(s.ctx, patchReq)
	require.NoError(s.T(), err)

	select {
	case <-time.After(time.Second):
		s.T().Fatal("timed out waiting for subscription to be received")
	case <-doneC:

	}

}

func Test_Subscription(t *testing.T) {
	suite.Run(t, new(SubscriptionSuite))
}
