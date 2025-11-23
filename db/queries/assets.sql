-- name: CreateAsset :one
INSERT INTO assets (id, upload_id, hls_root, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAsset :one
SELECT * FROM assets
WHERE id = $1 LIMIT 1;

-- name: GetAssetByUploadID :one
SELECT * FROM assets
WHERE upload_id = $1 LIMIT 1;

-- name: UpdateAssetStatus :one
UPDATE assets
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListAssets :many
SELECT * FROM assets
ORDER BY created_at DESC;
