package sqlite

import (
	"context"
	"time"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/storage/sqlite/queries/config_queries"
)

func (p *Provider) SetUpdatedAt(ctx context.Context, serviceName string, updatedAt time.Time) error {
	updatedAt = updatedAt.In(time.UTC)
	err := p.querier.SetUpdatedAt(ctx, config_queries.SetUpdatedAtParams{
		UpdatedAt: updatedAt,
		Name:      serviceName,
	})
	if err != nil {
		return rerrors.Wrap(err, "error updating updated_at for config")
	}

	return nil
}
