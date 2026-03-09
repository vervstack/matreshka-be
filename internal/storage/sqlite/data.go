package sqlite

import (
	"database/sql"
	"errors"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"go.vervstack.ru/matreshka/internal/clients/sqldb"
	"go.vervstack.ru/matreshka/internal/service/user_errors"
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

func wrapError(err error) error {
	var e *sqlite.Error
	ok := errors.As(err, &e)
	if !ok {
		return err
	}

	switch e.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		return user_errors.ErrAlreadyExists
	}

	return err
}
