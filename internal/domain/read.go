package domain

import (
	"time"

	"go.redsock.ru/evon"

	//"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

type ListConfigsRequest struct {
	Paging Paging
	Sort   Sort

	SearchPattern string
}

type ListConfigsResponse struct {
	Configs      []ConfigInfo
	TotalRecords uint64
}

type ConfigBase struct {
	Id   uint32
	Name string
	//Type      config_queries.MatreshkaConfigType
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
	Limit  uint64
	Offset uint64
}

type Sort struct {
	SortType api.Sort_Type
	Desc     bool
}
