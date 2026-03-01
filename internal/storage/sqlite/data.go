package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"go.vervstack.ru/matreshka/internal/clients/sqldb"
	"go.vervstack.ru/matreshka/internal/storage"
	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

type Provider struct {
	conn sqldb.DB

	querier *config_queries.Queries
}

func New(conn sqldb.DB) *Provider {
	return &Provider{
		conn: conn,

		querier: config_queries.New(conn),
	}
}

func (p *Provider) WithTx(tx *sql.Tx) storage.Data {
	return New(tx)
}
