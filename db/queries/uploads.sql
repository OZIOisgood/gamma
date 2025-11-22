-- name: CreateUpload :one
INSERT INTO uploads (id, title, s3_key, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUpload :one
SELECT * FROM uploads
WHERE id = $1 LIMIT 1;

-- name: ListUploads :many
SELECT * FROM uploads
ORDER BY created_at DESC;
