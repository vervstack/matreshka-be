package pg

import (
	"go.vervstack.ru/matreshka/internal/clients/sqldb"
	"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
)

type ConfigStorage struct {
	*config_queries.Queries
}

func New(conn sqldb.DB) *ConfigStorage {
	return &ConfigStorage{
		Queries: config_queries.New(conn),
	}
}
