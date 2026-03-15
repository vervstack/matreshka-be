package domain

import (
	"go.redsock.ru/evon"

	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

const MasterVersion = "master"

type ConfigNodes struct {
	Info  ConfigInfo
	Nodes *evon.Node
}

type GetConfigRawReq struct {
	Name    string
	Version string
	Format  matreshka_api.Format
}
type ConfigRawContent struct {
	Info ConfigInfo
	Raw  []byte
}

type CreateConfigRequest struct {
	Name   string
	Type   matreshka_api.ConfigType
	Format string
}
