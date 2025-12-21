package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgres_ParseFromDsn(t *testing.T) {
	t.Run("err_no_protocol", func(t *testing.T) {
		dsn := ""
		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.ErrorIs(t, err, ErrParsingPgDsn)
		require.Contains(t, err.Error(), "dsn must start with '"+pgProtocol+"'")
	})
	t.Run("err_no_user_colon", func(t *testing.T) {
		dsn := pgProtocol + ""
		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.ErrorIs(t, err, ErrParsingPgDsn)
		require.Contains(t, err.Error(), "dsn must have a colon for user:password separation")
	})
	t.Run("err_no_@", func(t *testing.T) {
		dsn := pgProtocol + "user:password"
		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.ErrorIs(t, err, ErrParsingPgDsn)
		require.Contains(t, err.Error(), "dsn must have a @ symbol for user:password@host:port separation")
	})
	t.Run("err_no_host_colon", func(t *testing.T) {
		dsn := pgProtocol + "user:password@host"
		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.ErrorIs(t, err, ErrParsingPgDsn)
		require.Contains(t, err.Error(), "dsn must have a colon for host:port separation")
	})
	t.Run("err_no_post_slash", func(t *testing.T) {
		dsn := pgProtocol + "user:password@host:port"
		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.ErrorIs(t, err, ErrParsingPgDsn)
		require.Contains(t, err.Error(), "dsn must have a slash symbol for user:password@host:port/database separation")
	})
	t.Run("err_not_a_number_port", func(t *testing.T) {
		dsn := pgProtocol + "user:password@host:546m/db"
		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.Contains(t, err.Error(), "error parsing port")
	})

	t.Run("OK_with_ssl", func(t *testing.T) {
		expected := Postgres{
			Host:    "local.host",
			Port:    5432,
			User:    "SomeGuy",
			Pwd:     "SomePass",
			DbName:  "GreatDb",
			SslMode: "disable",
		}

		dsn := expected.ConnectionString()

		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.NoError(t, err)

		require.Equal(t, pg, expected)
	})
	t.Run("OK_without_ssl", func(t *testing.T) {
		expected := Postgres{
			Host:    "local.host",
			Port:    5432,
			User:    "SomeGuy",
			Pwd:     "SomePass",
			DbName:  "GreatDb",
			SslMode: "disable",
		}

		dsn := "postgresql://SomeGuy:SomePass@local.host:5432/GreatDb"

		pg := Postgres{}
		err := pg.ParseFromDsn(dsn)
		require.NoError(t, err)

		require.Equal(t, pg, expected)
	})
}
