-- name: CreateComment :one
INSERT INTO comments (
    task_id,
    author_id,
    content
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCommentByID :one
SELECT * FROM comments
WHERE id = $1
LIMIT 1;

-- name: ListCommentsByTask :many
SELECT * FROM comments
WHERE task_id = $1
ORDER BY created_at ASC;

-- name: UpdateComment :one
UPDATE comments
SET
    content = $2
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;