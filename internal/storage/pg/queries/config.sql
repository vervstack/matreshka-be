-- name: CreateConfig :one
INSERT INTO matreshka.configs
    (name, type)
VALUES ($1, $2)
RETURNING id;

-- name: GetConfig :one
SELECT id,
       name,
       type,
       created_at,
       updated_at
FROM matreshka.configs
WHERE name = $1
FETCH FIRST 1 ROW ONLY;


-- name: DeleteByName :exec
WITH config_id AS (SELECT id
                   FROM matreshka.configs
                   WHERE matreshka.configs.name = $1),
     _delete_values AS (
         DELETE FROM matreshka.configs_content
             WHERE config_id = config_id.id
             RETURNING true)
DELETE
FROM matreshka.configs
WHERE id = config_id.id;


-- name :Store