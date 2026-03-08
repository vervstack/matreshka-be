package pg

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	errors "go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/domain"
)

const defaultLimit uint64 = 20

func (c *ConfigStorage) ListConfigs(ctx context.Context, req domain.ListConfigsRequest) (list domain.ListConfigsResponse, err error) {
	q := sq.Select().
		From("configs")

	if req.SearchPattern != "" {
		q = q.Where(sq.ILike{"name": req.SearchPattern})
	}

	list.TotalRecords, err = c.countTotal(ctx, q)
	if err != nil {
		return domain.ListConfigsResponse{}, errors.Wrap(err)
	}

	list.Configs, err = c.selectConfigsInfo(ctx, q, req.Paging)
	if err != nil {
		return domain.ListConfigsResponse{}, errors.Wrap(err)
	}

	return list, nil
}

func (c *ConfigStorage) selectConfigsInfo(ctx context.Context, q sq.SelectBuilder, paging domain.Paging) ([]domain.ConfigInfo, error) {
	q = q.Columns(
		"id", "name", "type", "created_at", "updated_at",
		"array_agg(configs_content.version)",
	)

	q = q.InnerJoin("configs_content ON configs_content.config_id = configs.id")
	q = q.GroupBy("id", "name", "type", "created_at", "updated_at")

	if paging.Limit == 0 {
		paging.Limit = defaultLimit
	}

	q = q.Limit(min(paging.Limit, defaultLimit))

	q = q.Offset(paging.Offset)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	rows, err := c.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, wrapPgErr(err)
	}
	defer rows.Close()

	configs, err := scanConfigInfo(rows)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return configs, nil
}

func (c *ConfigStorage) countTotal(ctx context.Context, q sq.SelectBuilder) (uint64, error) {
	var count uint64

	q = q.Column("count(*)")

	query, args, err := q.ToSql()
	if err != nil {
		return count, errors.Wrap(err)
	}

	err = c.conn.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, wrapPgErr(err)
	}

	return count, nil
}

func scanConfigInfo(rows *sql.Rows) ([]domain.ConfigInfo, error) {
	var out []domain.ConfigInfo

	for rows.Next() {
		next := domain.ConfigInfo{}
		err := rows.Scan(
			&next.Id,
			&next.Name,
			&next.Type,
			&next.CreatedAt,
			&next.UpdatedAt,
			&next.ConfigVersions,
		)
		if err != nil {
			return nil, wrapPgErr(err)
		}

		out = append(out, next)
	}

	return out, nil
}
