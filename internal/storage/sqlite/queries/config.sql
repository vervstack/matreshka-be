-- name: CreateConfig :one
INSERT INTO configs
    (name)
VALUES (?)
ON CONFLICT
    DO UPDATE SET name = excluded.name
RETURNING id;
-- name: GetConfig :one
SELECT id,
       name,
       updated_at
FROM configs
WHERE name = ?
LIMIT 1;


-- name: ClearValues :exec
DELETE
FROM configs_values
WHERE config_id = (SELECT id FROM configs WHERE name = ?)
  AND version LIKE ?;
