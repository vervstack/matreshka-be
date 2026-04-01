package matreshka_api_impl

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
)

func toConfigList(configs []domain.ConfigBase) []*api.ConfigBase {
	out := make([]*api.ConfigBase, 0, len(configs))
	for _, c := range configs {
		out = append(out, toConfig(c))
	}

	return out
}

func toConfig(cfg domain.ConfigBase) *api.ConfigBase {
	return &api.ConfigBase{
		Id:        cfg.Id,
		Name:      cfg.Name,
		CreatedAt: timestamppb.New(cfg.CreatedAt),
		UpdatedAt: timestamppb.New(cfg.UpdatedAt),
	}
}
