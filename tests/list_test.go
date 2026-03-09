package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

const (
	listService1 = "service_a"
	listService2 = "service_c"
	listService3 = "service_b"
)

type ListSuite struct {
	suite.Suite

	ctx context.Context
	app AppEnv
}

func (s *ListSuite) SetupTest() {
	s.ctx = context.Background()
	s.app = InitAppEnvironment(s.T())

	s.app.createEmptyService(s.T(), listService1, matreshka_api.ConfigType_verv)
	s.app.createEmptyService(s.T(), listService2, matreshka_api.ConfigType_verv)
	s.app.createEmptyService(s.T(), listService3, matreshka_api.ConfigType_verv)
}

func (s *ListSuite) Test_List() {
	type testCase struct {
		pattern string
		paging  *matreshka_api.Paging
		sorting *matreshka_api.Sort

		namesOrder []string
		total      uint64
	}

	testCases := map[string]testCase{
		"default": {
			namesOrder: []string{listService1, listService2, listService3},
			total:      3,
		},

		"order_by_name": {
			sorting: &matreshka_api.Sort{
				Type: matreshka_api.Sort_by_name,
			},
			namesOrder: []string{listService1, listService3, listService2},
			total:      3,
		},
		"order_by_name_desc": {
			sorting: &matreshka_api.Sort{
				Type: matreshka_api.Sort_by_name,
				Desc: true,
			},
			namesOrder: []string{listService2, listService3, listService1},
			total:      3,
		},

		"order_by_updated_at": {
			sorting: &matreshka_api.Sort{
				Type: matreshka_api.Sort_by_updated_at,
			},
			namesOrder: []string{listService1, listService2, listService3},
			total:      3,
		},
		"order_by_updated_at_desc": {
			sorting: &matreshka_api.Sort{
				Type: matreshka_api.Sort_by_updated_at,
			},
			namesOrder: []string{listService1, listService2, listService3},
			total:      3,
		},

		"filter_by_name_one_to_one": {
			pattern:    listService2,
			namesOrder: []string{listService2},
			total:      1,
		},

		"filter_by_name_part": {
			pattern:    listService2[len(listService2)-2:],
			namesOrder: []string{listService2},
			total:      1,
		},

		"paging_limit": {
			paging: &matreshka_api.Paging{
				Limit:  1,
				Offset: 0,
			},
			namesOrder: []string{listService1},
			total:      3,
		},
		"paging_limit_offset": {
			paging: &matreshka_api.Paging{
				Limit:  1,
				Offset: 1,
			},
			namesOrder: []string{listService2},
			total:      3,
		},
	}

	for name, tc := range testCases {
		s.T().Run(name, func(t *testing.T) {
			req := &matreshka_api.ListConfigs_Request{
				Paging:        tc.paging,
				Sort:          tc.sorting,
				SearchPattern: &tc.pattern,
			}

			resp, err := s.app.matreshkaApi.ListConfigs(s.ctx, req)
			require.NoError(t, err)

			s.verifyOrder(t, resp, tc.namesOrder)

			require.Equal(t, tc.total, resp.TotalRecords)
		})

	}

}

func (s *ListSuite) verifyOrder(t *testing.T, resp *matreshka_api.ListConfigs_Response, namesOrder []string) {
	actualNames := make([]string, 0, len(namesOrder))

	for idx := range resp.Configs {
		actualNames = append(actualNames, resp.Configs[idx].Name)
	}

	require.Equal(t, namesOrder, actualNames)
}

func Test_List(t *testing.T) {
	suite.Run(t, new(ListSuite))
}
