package matreshka_api_impl

import (
	errors "go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

var ErrUnknownConfigType = errors.New("unknown config type", codes.InvalidArgument)

var configTypeMappingToDomain map[api.ConfigType]config_queries.MatreshkaConfigType

func init() {
	configTypeMappingToDomain = map[api.ConfigType]config_queries.MatreshkaConfigType{
		api.ConfigType_plain: config_queries.MatreshkaConfigTypePlain,
		api.ConfigType_verv:  config_queries.MatreshkaConfigTypeVerv,
		api.ConfigType_minio: config_queries.MatreshkaConfigTypeMinio,
		api.ConfigType_pg:    config_queries.MatreshkaConfigTypePg,
		api.ConfigType_nginx: config_queries.MatreshkaConfigTypeNginx,
		api.ConfigType_kv:    config_queries.MatreshkaConfigTypeKv,
	}
}

func fromConfigType(configType api.ConfigType) (config_queries.MatreshkaConfigType, error) {
	tp, ok := configTypeMappingToDomain[configType]
	if !ok {
		return "", ErrUnknownConfigType
	}
	return tp, nil
}

func toConfigList(configs []domain.ConfigInfo) []*api.Config {
	out := make([]*api.Config, 0, len(configs))
	for _, c := range configs {
		out = append(out, toConfig(c))
	}

	return out
}

func toConfig(cfg domain.ConfigInfo) *api.Config {
	return &api.Config{
		Id:                    cfg.Id,
		Name:                  cfg.Name,
		CreatedAtUtcTimestamp: cfg.CreatedAt.UTC().Unix(),
		UpdatedAtUtcTimestamp: cfg.UpdatedAt.UTC().Unix(),
		Versions:              cfg.ConfigVersions,
	}
}
