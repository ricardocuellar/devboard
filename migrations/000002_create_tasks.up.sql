
CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL,                  
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',        
    status      TEXT NOT NULL DEFAULT 'todo'
                CHECK (status IN ('todo', 'in_progress', 'done')), 
    assignee_id UUID REFERENCES users(id)       
                ON DELETE SET NULL,             
    created_by  UUID NOT NULL REFERENCES users(id)
                ON DELETE RESTRICT,             
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);