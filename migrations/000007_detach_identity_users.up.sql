-- tasks → Identity
ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS
    tasks_assignee_id_fkey;

ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS
    tasks_created_by_fkey;


-- workspaces → Identity
ALTER TABLE workspaces
    DROP CONSTRAINT IF EXISTS
    workspaces_created_by_fkey;

ALTER TABLE workspace_members
    DROP CONSTRAINT IF EXISTS
    workspace_members_user_id_fkey;


-- projects → Identity
ALTER TABLE projects
    DROP CONSTRAINT IF EXISTS
    projects_created_by_fkey;


-- comments → Identity
ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS
    comments_author_id_fkey;


-- users ya no pertenece a este database.
-- Conservamos un snapshot temporal para rollback.
ALTER TABLE users
    RENAME TO legacy_users;

ALTER TABLE legacy_users
    RENAME CONSTRAINT
    users_pkey
    TO legacy_users_pkey;

ALTER TABLE legacy_users
    RENAME CONSTRAINT
    users_email_key
    TO legacy_users_email_key;

ALTER INDEX IF EXISTS idx_users_email
    RENAME TO idx_legacy_users_email;