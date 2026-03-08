package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox"

	"go.vervstack.ru/matreshka/internal/domain"
	api "go.vervstack.ru/matreshka/pkg/matreshka_api"
)

const defaultPageSize = 20

func (p *Provider) ListConfigs(ctx context.Context, req domain.ListConfigsRequest) (out domain.ListConfigsResponse, err error) {
	pattern := sql.NullString{
		String: req.SearchPattern,
		Valid:  true,
	}

	totalRecords, err := p.querier.ListConfigsCount(ctx, pattern)
	if err != nil {
		return domain.ListConfigsResponse{}, rerrors.Wrap(err, "error scanning total amount of configs")
	}
	out.TotalRecords = uint64(totalRecords)

	q := `
		WITH cfg AS (
			SELECT
				configs.id 			AS id,
				configs.updated_at 	AS updated_at,
				configs.name 		AS name
			FROM configs
			WHERE name LIKE '%'||$1||'%'
			GROUP BY configs.name
		),
		versions AS (
			SELECT
				cv.config_id as config_id,
				cv.version version
			FROM configs_values cv
			INNER JOIN cfg c on c.id = cv.config_id
			GROUP BY config_id, version
			UNION ALL
			SELECT
			    c.id,
			    'master'
			FROM cfg c
		)
		SELECT
			cfg.name 						    AS config_name,
			cfg.updated_at 					    AS last_updated_at,
			json_group_array(versions.version) AS config_versions
		FROM cfg
		LEFT JOIN versions ON versions.config_id = cfg.id

		GROUP BY cfg.id
		HAVING COUNT(cfg.id) > 0  -- Ensures only non-empty results are returned
		`
	args := []any{req.SearchPattern}

	q += "\nORDER BY " + extractSort(req.Sort)
	q += fmt.Sprintf("\nLIMIT %d OFFSET %d",
		toolbox.Coalesce(req.Paging.Limit, defaultPageSize),
		req.Paging.Offset)

	rows, err := p.conn.QueryContext(ctx, q, args...)
	if err != nil {
		return domain.ListConfigsResponse{}, rerrors.Wrap(err, "error listing configs")
	}
	defer rows.Close()

	out.Configs = make([]domain.ConfigInfo, 0, req.Paging.Limit)

	for rows.Next() {
		var item domain.ConfigInfo
		var versionsJSON string
		err = rows.Scan(
			&item.Name,
			&item.UpdatedAt,
			&versionsJSON,
		)
		if err != nil {
			return out, rerrors.Wrap(err, "error scanning row")
		}

		err = json.Unmarshal([]byte(versionsJSON), &item.ConfigVersions)
		if err != nil {
			return out, rerrors.Wrap(err, "error marshalling from json ")
		}
		sort.Slice(item.ConfigVersions, func(i, j int) bool {
			return item.ConfigVersions[i] < item.ConfigVersions[j]
		})

		for i := range item.ConfigVersions {
			if item.ConfigVersions[i] == domain.MasterVersion {
				item.ConfigVersions[0], item.ConfigVersions[i] =
					item.ConfigVersions[i], item.ConfigVersions[0]

				break
			}
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

func extractSort(sort domain.Sort) (field string) {
	switch sort.SortType {
	case api.Sort_default:
		field = "id"
	case api.Sort_by_name:
		field = "name"
	default:
		field = "updated_at"
	}
	if sort.Desc {
		field += " DESC"
	}

	return
}
