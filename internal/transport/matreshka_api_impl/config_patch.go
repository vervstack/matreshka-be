package matreshka_api_impl

import (
	"context"
	"strings"

	errors "go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
	"go.vervstack.ru/matreshka/internal/domain"
)

func (s *Impl) PatchConfig(ctx context.Context, req *api.PatchConfig_Request) (*api.PatchConfig_Response, error) {
	patch, err := fromPatch(req)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = s.configService.Patch(ctx, patch)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &api.PatchConfig_Response{}, nil
}

func fromPatch(req *api.PatchConfig_Request) (domain.PatchConfigRequest, error) {
	out := domain.PatchConfigRequest{
		ConfigVersion: toolbox.Coalesce(toolbox.FromPtr(req.Version), domain.MasterVersion),
	}

	out.ConfigName = req.ConfigName

	for _, patch := range req.Patches {

		switch v := patch.GetPatch().(type) {
		case *api.Patch_Rename:
			out.RenameTo = append(out.RenameTo,
				domain.PatchRename{
					OldName: patch.FieldName,
					NewName: v.Rename,
				})
		case *api.Patch_UpdateValue:
			out.Upsert = append(out.Upsert, domain.PatchUpdate{
				FieldName:  patch.FieldName,
				FieldValue: v.UpdateValue,
			})
		case *api.Patch_Delete:
			out.Delete = append(out.Delete, patch.FieldName)
		}
	}

	return out, nil
}

func ParseConfigName(name string) (*api.ConfigType, string) {
	nameSplited := strings.Split(name, "_")
	if len(nameSplited) < 2 {
		return nil, name
	}

	pref, ok := api.ConfigType_value[nameSplited[0]]
	if ok {
		return toolbox.ToPtr(api.ConfigType(pref)), strings.Join(nameSplited[1:], "_")
	}

	return nil, name

}
