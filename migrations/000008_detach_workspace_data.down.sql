ALTER TABLE legacy_workspaces
    RENAME TO workspaces;

ALTER TABLE legacy_projects
    RENAME TO projects;

ALTER TABLE legacy_workspace_members
    RENAME TO workspace_members;


ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_project_id
    FOREIGN KEY (project_id)
    REFERENCES projects(id)
    ON DELETE CASCADE;
