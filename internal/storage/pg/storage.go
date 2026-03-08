package pg

import (
	"go.vervstack.ru/matreshka/internal/clients/sqldb"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
)

type ConfigStorage struct {
	conn sqldb.DB
	*config_queries.Queries
}

func New(conn sqldb.DB) storage.ConfigStorage {
	return &ConfigStorage{
		Queries: config_queries.New(conn),

		conn: conn,
	}
}
