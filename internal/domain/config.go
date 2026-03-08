package domain

import (
	"go.redsock.ru/evon"

	"go.vervstack.ru/matreshka/pkg/matreshka_api"
)

const MasterVersion = "master"

type ConfigWithNodes struct {
	Type matreshka_api.ConfigType

	Nodes    *evon.Node
	Versions []string
}

type CreateConfigRequest struct {
	Name   string
	Type   matreshka_api.ConfigType
	Format string
}
