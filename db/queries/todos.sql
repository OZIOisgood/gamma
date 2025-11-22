-- name: CreateTodo :one
INSERT INTO todos (task, completed)
VALUES ($1, $2)
RETURNING *;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1 LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY created_at DESC;

-- name: UpdateTodo :one
UPDATE todos
SET task = $2, completed = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
