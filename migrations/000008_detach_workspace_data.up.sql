-- Task ya no puede depender físicamente
-- de la tabla projects de Workspace.

ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS
        fk_tasks_project_id;


-- Conservamos snapshots temporales
-- para rollback y verificación.

ALTER TABLE workspace_members
    RENAME TO legacy_workspace_members;

ALTER TABLE projects
    RENAME TO legacy_projects;

ALTER TABLE workspaces
    RENAME TO legacy_workspaces;
