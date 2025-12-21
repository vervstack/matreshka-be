package resources

import (
	"fmt"
	"strconv"
	"strings"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
)

var ErrParsingPgDsn = rerrors.New("error parsing postgres dsn. Expected format: " + pgProtocol + "user:password@host:port/dbname")

const (
	PostgresResourceName = "postgres"
	pgProtocol           = "postgresql://"
)

type Postgres struct {
	Name `yaml:"resource_name" env:"-"`

	MigrationsFolder `yaml:"migrations_folder"`

	Host string `yaml:"host"`
	Port uint64 `yaml:"port"`

	User string `yaml:"user"`
	Pwd  string `yaml:"pwd"`

	DbName  string `yaml:"name"`
	SslMode string `yaml:"ssl_mode"`
}

func NewPostgres(n Name) Resource {
	return &Postgres{
		Name:   n,
		Host:   "0.0.0.0",
		Port:   5432,
		User:   "postgres",
		Pwd:    "",
		DbName: "postgres",
	}
}

func (p *Postgres) GetType() string {
	return PostgresResourceName
}

func (p *Postgres) MarshalYAML() (interface{}, error) {
	resourceType := strings.Split(p.GetName(), evon.ObjectSplitter)[0]
	if resourceType != "postgres" {
		return nil, rerrors.Wrap(ErrInvalidResourceName, "but got: "+resourceType)
	}

	return *p, nil
}

func (p *Postgres) ConnectionString() string {
	connstr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		p.User,
		p.Pwd,
		p.Host,
		p.Port,
		p.DbName,
	)

	if p.SslMode != "" {
		connstr += fmt.Sprintf("?sslmode=%s", p.SslMode)
	}

	return connstr
}

func (p *Postgres) ParseFromDsn(dsn string) (err error) {
	if !strings.HasPrefix(dsn, pgProtocol) {
		return rerrors.Wrap(ErrParsingPgDsn, "dsn must start with '"+pgProtocol+"'")
	}

	dsn = dsn[len(pgProtocol):]

	colonIdx := strings.Index(dsn, ":")
	if colonIdx == -1 {
		return rerrors.Wrap(ErrParsingPgDsn, "dsn must have a colon for user:password separation")
	}
	p.User = dsn[:colonIdx]
	dsn = dsn[colonIdx+1:]

	atIndex := strings.Index(dsn, "@")
	if atIndex == -1 {
		return rerrors.Wrap(ErrParsingPgDsn, "dsn must have a @ symbol for user:password@host:port separation")
	}

	p.Pwd = dsn[:atIndex]
	dsn = dsn[atIndex+1:]

	colonIdx = strings.Index(dsn, ":")
	if colonIdx == -1 {
		return rerrors.Wrap(ErrParsingPgDsn, "dsn must have a colon for host:port separation")
	}

	p.Host = dsn[:colonIdx]
	dsn = dsn[colonIdx+1:]

	slashIdx := strings.Index(dsn, "/")
	if slashIdx == -1 {
		return rerrors.Wrap(ErrParsingPgDsn, "dsn must have a slash symbol for user:password@host:port/database separation")
	}

	port := dsn[:slashIdx]
	p.Port, err = strconv.ParseUint(port, 10, 64)
	if err != nil {
		return rerrors.Wrap(err, "error parsing port")
	}

	dsn = dsn[slashIdx+1:]

	questionMarkIdx := strings.Index(dsn, "?")
	if questionMarkIdx == -1 {
		questionMarkIdx = len(dsn)
	}

	p.DbName = dsn[:questionMarkIdx]

	dsn = dsn[questionMarkIdx:]

	for len(dsn) != 0 {
		equalsIdx := strings.Index(dsn, "=")
		if equalsIdx == -1 {
			break
		}

		key := dsn[:equalsIdx]
		dsn = dsn[equalsIdx+1:]
		ampersantIdx := strings.Index(dsn, "&")
		if ampersantIdx == -1 {
			ampersantIdx = len(dsn)
		}

		value := dsn[:ampersantIdx]
		dsn = dsn[ampersantIdx:]

		switch key {
		case "sslmode":
			p.SslMode = value

		}
	}

	if p.SslMode == "" {
		p.SslMode = "disable"
	}

	return nil
}

func (p *Postgres) SqlDialect() string {
	return "postgres"
}
