-- name: CreateUpload :one
INSERT INTO uploads (id, title, s3_key, status, realm_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUpload :one
SELECT * FROM uploads
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListUploads :many
SELECT * FROM uploads
WHERE deleted_at IS NULL AND realm_id = $1
ORDER BY created_at DESC;

-- name: UpdateUploadStatusByKey :one
UPDATE uploads
SET status = $2, updated_at = NOW()
WHERE s3_key = $1
RETURNING *;

-- name: SoftDeleteUpload :one
UPDATE uploads
SET status = 'deleted', deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteUploadsByRealmID :exec
UPDATE uploads
SET status = 'deleted', deleted_at = NOW(), updated_at = NOW()
WHERE realm_id = $1 AND deleted_at IS NULL;
