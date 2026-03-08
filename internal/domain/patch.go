package domain

import (
	"go.redsock.ru/evon"

	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type PatchConfigRequest struct {
	ConfigName    string
	ConfigType    api.ConfigType
	ConfigVersion string

	Upsert   []PatchUpdate
	RenameTo []PatchRename
	Delete   []string
}

type PatchUpdate struct {
	FieldName  string
	FieldValue string
}

type PatchRename struct {
	OldName string
	NewName string
}

type ReplaceConfigReq struct {
	Name    string
	Version string
	Config  *evon.Node
}
