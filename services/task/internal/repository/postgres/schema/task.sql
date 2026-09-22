CREATE TABLE tasks (
    id UUID PRIMARY KEY,

    project_id UUID NOT NULL,

    title TEXT NOT NULL,

    description TEXT NOT NULL,

    status TEXT NOT NULL,

    assignee_id UUID,

    created_by UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,

    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE comments (
    id UUID PRIMARY KEY,

    task_id UUID NOT NULL,

    author_id UUID NOT NULL,

    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,

    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT comments_task_id_fkey
        FOREIGN KEY (task_id)
        REFERENCES tasks(id)
        ON DELETE CASCADE
);