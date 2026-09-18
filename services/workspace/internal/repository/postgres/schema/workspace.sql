CREATE TABLE workspaces (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE workspace_members (
    id UUID PRIMARY KEY,
    workspace_id UUID NOT NULL
        REFERENCES workspaces(id)
        ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL,
    UNIQUE(workspace_id, user_id)
);

CREATE TABLE projects (
    id UUID PRIMARY KEY,
    workspace_id UUID NOT NULL
        REFERENCES workspaces(id)
        ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);