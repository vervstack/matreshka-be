package domain

import (
	"time"

	"go.redsock.ru/evon"

	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type ListConfigsRequest struct {
	Paging Paging
	Sort   Sort

	SearchPattern string
}

type ListConfigsResponse struct {
	List         []ConfigInfo
	TotalRecords uint32
}

type ConfigBase struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConfigInfo struct {
	ConfigBase
	ConfigVersions []string
}

type ConfigDescription struct {
	Id          int64
	ServiceName string
}

type ConfigEnvVals struct {
	ConfigDescription
	Nodes *evon.Node
}

type Paging struct {
	Limit  uint32
	Offset uint32
}

type Sort struct {
	SortType api.Sort_Type
	Desc     bool
}
