-- name: CreateWorkspace :one
INSERT INTO workspaces (
    name,
    slug,
    description,
    created_by
) VALUES (
    $1, -- name
    $2, -- slug
    $3, -- description
    $4  -- created_by
)
RETURNING *;

-- name: GetWorkspaceByID :one
SELECT * FROM workspaces
WHERE id = $1
LIMIT 1;

-- name: GetWorkspaceBySlug :one
SELECT * FROM workspaces
WHERE slug = $1
LIMIT 1;

-- name: ListWorkspacesByUser :many
SELECT w.* FROM workspaces w
INNER JOIN workspace_members wm ON wm.workspace_id = w.id
WHERE wm.user_id = $1
ORDER BY w.created_at DESC;


-- name: AddWorkspaceMember :one
INSERT INTO workspace_members (
    workspace_id,
    user_id,
    role
) VALUES (
    $1, -- workspace_id
    $2, -- user_id
    $3  -- role
)
RETURNING *;

-- name: GetWorkspaceMember :one
SELECT * FROM workspace_members
WHERE workspace_id = $1 AND user_id = $2
LIMIT 1;

-- name: ListWorkspaceMembers :many
SELECT * FROM workspace_members
WHERE workspace_id = $1
ORDER BY joined_at ASC;

-- name: UpdateWorkspaceMemberRole :one
UPDATE workspace_members
SET role = $3
WHERE workspace_id = $1 AND user_id = $2
RETURNING *;

-- name: RemoveWorkspaceMember :exec
DELETE FROM workspace_members
WHERE workspace_id = $1 AND user_id = $2;