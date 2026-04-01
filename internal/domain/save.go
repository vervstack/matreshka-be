package domain

import (
	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
)

type SaveConfigReq struct {
	ConfigName string
	Version    *string

	Format  api.Format
	Content []byte
}
