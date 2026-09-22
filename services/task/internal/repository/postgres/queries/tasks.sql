-- name: CreateTask :one
INSERT INTO tasks (
    project_id,
    title,
    description,
    status,
    created_by
) VALUES (
    $1,-- project_id
    $2,-- title
    $3,-- description
    $4,-- status
    $5 -- created_by
)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1
LIMIT 1;

-- name: ListTasksByProject :many
SELECT * FROM tasks
WHERE project_id = $1
ORDER BY created_at ASC;


-- name: UpdateTask :one
UPDATE tasks
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    status = COALESCE($4, status),
    assignee_id = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

