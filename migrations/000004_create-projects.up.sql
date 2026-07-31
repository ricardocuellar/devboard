CREATE TABLE projects (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id)
                 ON DELETE CASCADE,  -- si se borra el workspace, se borran sus proyectos
    name         TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    created_by   UUID NOT NULL REFERENCES users(id)
                 ON DELETE RESTRICT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_workspace_id ON projects(workspace_id);