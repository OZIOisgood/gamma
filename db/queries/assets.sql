-- name: CreateAsset :one
INSERT INTO assets (id, upload_id, realm_id, hls_root, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAsset :one
SELECT * FROM assets
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetAssetByUploadID :one
SELECT * FROM assets
WHERE upload_id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: UpdateAssetStatus :one
UPDATE assets
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListAssets :many
SELECT * FROM assets
WHERE deleted_at IS NULL AND realm_id = $1
ORDER BY created_at DESC;

-- name: SoftDeleteAsset :one
UPDATE assets
SET status = 'deleted', deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteAssetsByRealmID :exec
UPDATE assets
SET status = 'deleted', deleted_at = NOW(), updated_at = NOW()
WHERE realm_id = $1 AND deleted_at IS NULL;
