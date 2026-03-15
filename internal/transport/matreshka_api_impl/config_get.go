package matreshka_api_impl

import (
	"context"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.vervstack.ru/matreshka/internal/domain"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

func (s *Impl) GetConfig(ctx context.Context, req *api.GetConfig_Request) (*api.GetConfig_Response, error) {

	rawReq := domain.GetConfigRawReq{
		Name:    req.ConfigName,
		Version: toolbox.Coalesce(toolbox.FromPtr(req.Version), domain.MasterVersion),
		Format:  req.Format,
	}

	cfg, err := s.configService.GetConfigRaw(ctx, rawReq)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	resp := &api.GetConfig_Response{
		Config: cfg.Raw,
		BaseInfo: &api.ConfigBase{
			Id:         cfg.Info.Id,
			Name:       cfg.Info.Name,
			CreatedAt:  timestamppb.New(cfg.Info.CreatedAt),
			UpdatedAt:  timestamppb.New(cfg.Info.UpdatedAt),
			ConfigType: cfg.Info.Type,
		},
	}
	//
	//switch req.Format {
	//case api.Format_env:
	//	resp.Config = evon.Marshal(cfg.Nodes.InnerNodes)
	//default:
	//	switch cfg.BaseInfo.Type {
	//	case api.ConfigType_verv:
	//		resp.Config, err = vervToYaml(cfg.Nodes)
	//	default:
	//		resp.Config, err = kvToYaml(cfg.Nodes)
	//	}
	//	if err != nil {
	//		return nil, rerrors.Wrap(err)
	//	}
	//}

	return resp, nil
}
