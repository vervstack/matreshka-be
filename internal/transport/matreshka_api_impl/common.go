package matreshka_api_impl

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.vervstack.ru/matreshka/internal/domain"
	//"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func toConfigList(configs []domain.ConfigInfo) []*api.ConfigBase {
	out := make([]*api.ConfigBase, 0, len(configs))
	for _, c := range configs {
		out = append(out, toConfig(c))
	}

	return out
}

func toConfig(cfg domain.ConfigInfo) *api.ConfigBase {
	return &api.ConfigBase{
		Id:        cfg.Id,
		Name:      cfg.Name,
		CreatedAt: timestamppb.New(cfg.CreatedAt),
		UpdatedAt: timestamppb.New(cfg.UpdatedAt),
		Versions:  cfg.ConfigVersions,
	}
}
