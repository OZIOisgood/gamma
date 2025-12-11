-- name: ListRealms :many
SELECT * FROM realms
WHERE status != 'deleted'
ORDER BY name ASC;

-- name: GetRealmByName :one
SELECT * FROM realms
WHERE name = $1 AND status != 'deleted' LIMIT 1;

-- name: CreateRealm :one
INSERT INTO realms (name)
VALUES ($1)
RETURNING *;

-- name: GetRealm :one
SELECT * FROM realms
WHERE id = $1 AND status != 'deleted' LIMIT 1;

-- name: DeleteRealm :exec
UPDATE realms
SET status = 'deleted', deleted_at = NOW()
WHERE id = $1;
