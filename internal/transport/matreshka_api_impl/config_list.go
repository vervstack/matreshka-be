package matreshka_api_impl

import (
	"context"

	"go.redsock.ru/toolbox"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.vervstack.ru/matreshka/internal/domain"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (s *Impl) ListConfigs(ctx context.Context, req *api.ListConfigs_Request) (*api.ListConfigs_Response, error) {
	listReq := domain.ListConfigsRequest{
		Paging: domain.Paging{
			Limit:  toolbox.Coalesce(req.GetPaging().GetLimit(), 10),
			Offset: req.GetPaging().GetOffset(),
		},
		Sort: domain.Sort{
			SortType: req.Sort.GetType(),
			Desc:     req.Sort.GetDesc(),
		},

		SearchPattern: req.GetSearchPattern(),
	}

	list, err := s.configService.List(ctx, listReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &api.ListConfigs_Response{
		Configs:      toConfigList(list.Configs),
		TotalRecords: list.TotalRecords,
	}

	return resp, nil
}
