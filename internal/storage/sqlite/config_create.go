package sqlite

import (
	"context"

	errors "go.redsock.ru/rerrors"
)

func (p *Provider) Create(ctx context.Context, serviceName string) (int64, error) {
	cfgId, err := p.querier.CreateConfig(ctx, serviceName)
	if err != nil {
		return 0, errors.Wrap(err, "error upserting config")
	}
	if cfgId != 0 {
		return cfgId, nil
	}

	cfg, err := p.querier.GetConfig(ctx, serviceName)
	if err != nil {
		return 0, errors.Wrap(err, "error getting config")
	}

	return cfg.ID, nil
}
