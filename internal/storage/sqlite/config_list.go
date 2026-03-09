package sqlite

import (
	"context"
	"encoding/json"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"go.redsock.ru/rerrors"

	"go.vervstack.ru/matreshka/internal/clients/sqldb"
	"go.vervstack.ru/matreshka/internal/domain"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

const defaultPageSize = 20

func (p *Provider) ListConfigs(ctx context.Context, req domain.ListConfigsRequest) (out domain.ListConfigsResponse, err error) {
	q := sq.Select().
		From("configs")

	if req.SearchPattern != "" {
		req.SearchPattern = strings.ReplaceAll(req.SearchPattern, "_", "\\_")

		q = q.Where(sq.Expr("name LIKE ? ESCAPE '\\'", "%"+req.SearchPattern+"%"))
	}

	totalRecords, err := p.countItems(ctx, q)
	if err != nil {
		return domain.ListConfigsResponse{}, rerrors.Wrap(err, "error scanning total amount of configs")
	}

	out.TotalRecords = uint64(totalRecords)

	q = q.Columns(
		"id",
		"name",
		"type_name",
		"created_at",
		"updated_at",
	)
	q = applySorting(q, req.Sort)
	q = applyPaging(q, req.Paging)

	query, args, err := q.ToSql()
	if err != nil {
		return domain.ListConfigsResponse{}, rerrors.Wrap(err, "error building query")
	}

	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return domain.ListConfigsResponse{}, rerrors.Wrap(err, "error listing configs")
	}
	defer closeRows(rows)

	out.Configs = make([]domain.ConfigBase, 0, req.Paging.Limit)

	for rows.Next() {
		var item domain.ConfigBase
		item, err = scanConfigBase(rows)
		if err != nil {
			return domain.ListConfigsResponse{}, rerrors.Wrap(err, "error scanning config")
		}

		out.Configs = append(out.Configs, item)
	}

	return out, nil
}

func (p *Provider) GetVersions(ctx context.Context, name string) ([]string, error) {
	versionsRaw, err := p.querier.GetVersions(ctx, name)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting versions")
	}

	versionsStr, ok := versionsRaw.(string)
	if !ok {
		// sqlc might return []byte if it's not sure
		if b, ok := versionsRaw.([]byte); ok {
			versionsStr = string(b)
		} else {
			return nil, rerrors.New("unexpected type from GetVersions")
		}
	}

	var versions []string
	err = json.Unmarshal([]byte(versionsStr), &versions)
	if err != nil {
		return nil, rerrors.Wrap(err, "error unmarshalling versions from json ")
	}

	return versions, nil
}

func (p *Provider) countItems(ctx context.Context, q sq.SelectBuilder) (int64, error) {
	var total int64
	query, args, err := q.Column("count(*)").ToSql()
	if err != nil {
		return 0, rerrors.Wrap(err, "error building query for countItems")
	}

	err = p.conn.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, wrapError(err)
	}

	return total, nil
}

func scanConfigBase(rows sqldb.RowScanner) (domain.ConfigBase, error) {
	var item domain.ConfigBase

	var typeName string
	err := rows.Scan(
		&item.Id,
		&item.Name,
		&typeName,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return item, rerrors.Wrap(err, "error scanning row")
	}

	item.Type = api.ConfigType(api.ConfigType_value[typeName])

	return item, nil
}

func applySorting(q sq.SelectBuilder, sort domain.Sort) sq.SelectBuilder {
	direction := ascSort
	if sort.Desc {
		direction = descSort
	}

	switch sort.SortType {
	case api.Sort_by_name:
		return q.OrderBy("name " + direction)
	case api.Sort_by_updated_at:
		return q.OrderBy("updated_at " + direction)
	default:
		return q.OrderBy("id " + direction)
	}
}

func applyPaging(q sq.SelectBuilder, paging domain.Paging) sq.SelectBuilder {
	if paging.Limit == 0 {
		paging.Limit = defaultPageSize
	} else {
		paging.Limit = min(paging.Limit, defaultPageSize)
	}

	q = q.Limit(paging.Limit).
		Offset(paging.Offset)

	return q
}
