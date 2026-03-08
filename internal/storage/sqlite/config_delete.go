package sqlite

import (
	"context"

	"go.redsock.ru/rerrors"
)

func (p *Provider) Delete(ctx context.Context, name string) error {
	err := p.querier.DeleteConfig(ctx, name)
	if err != nil {
		return rerrors.Wrap(err, "error deleting config")
	}

	return nil
}
