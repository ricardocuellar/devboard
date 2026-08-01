-- name: CreateProject :one
INSERT INTO projects (
    workspace_id,
    name,
    description,
    created_by
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1
LIMIT 1;

-- name: ListProjectsByWorkspace :many
SELECT * FROM projects
WHERE workspace_id = $1
ORDER BY created_at ASC;

-- name: UpdateProject :one
UPDATE projects
SET
    name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;