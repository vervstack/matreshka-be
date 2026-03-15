package storage

import (
	"context"
	"database/sql"
	"time"

	"go.redsock.ru/evon"

	"go.vervstack.ru/matreshka/internal/domain"
	//"go.vervstack.ru/matreshka/internal/storage/pg/queries/config_queries"
)

type Data interface {
	GetConfigByName(ctx context.Context, name string) (domain.ConfigBase, error)

	// GetConfigNodes returns root node of parsed config
	GetConfigNodes(ctx context.Context, name string, ver string) (*evon.Node, error)
	GetConfigRawContent(ctx context.Context, name string, version string) ([]byte, error)
	GetVersions(ctx context.Context, name string) ([]string, error)
	ListConfigs(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error)

	Create(ctx context.Context, request domain.CreateConfigRequest) (int64, error)

	UpsertValues(ctx context.Context, req domain.PatchConfigRequest) error
	DeleteValues(ctx context.Context, req domain.PatchConfigRequest) error
	RenameValues(ctx context.Context, req domain.PatchConfigRequest) error
	// ClearValues - removes values from storage
	// if version is not selected - consider to delete in
	ClearValues(ctx context.Context, name string, version *string) error

	SetUpdatedAt(ctx context.Context, name string, req time.Time) error

	Rename(ctx context.Context, oldName, newName string) error
	Delete(ctx context.Context, name string) error

	WithTx(tx *sql.Tx) Data
}
