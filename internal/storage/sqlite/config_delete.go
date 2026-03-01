package sqlite

import (
	"context"

	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
)

func (p *Provider) Delete(ctx context.Context, name domain.ConfigName) error {
	err := p.querier.DeleteConfig(ctx, name.Name())
	if err != nil {
		return rerrors.Wrap(err, "error deleting config")
	}

	return nil
}
