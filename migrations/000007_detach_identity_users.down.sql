ALTER INDEX IF EXISTS idx_legacy_users_email
    RENAME TO idx_users_email;

ALTER TABLE legacy_users
    RENAME CONSTRAINT
    legacy_users_email_key
    TO users_email_key;

ALTER TABLE legacy_users
    RENAME CONSTRAINT
    legacy_users_pkey
    TO users_pkey;

ALTER TABLE legacy_users
    RENAME TO users;


ALTER TABLE tasks
    ADD CONSTRAINT tasks_assignee_id_fkey
    FOREIGN KEY (assignee_id)
    REFERENCES users(id)
    ON DELETE SET NULL;

ALTER TABLE tasks
    ADD CONSTRAINT tasks_created_by_fkey
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON DELETE RESTRICT;


ALTER TABLE workspaces
    ADD CONSTRAINT workspaces_created_by_fkey
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON DELETE RESTRICT;

ALTER TABLE workspace_members
    ADD CONSTRAINT workspace_members_user_id_fkey
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;


ALTER TABLE projects
    ADD CONSTRAINT projects_created_by_fkey
    FOREIGN KEY (created_by)
    REFERENCES users(id)
    ON DELETE RESTRICT;


ALTER TABLE comments
    ADD CONSTRAINT comments_author_id_fkey
    FOREIGN KEY (author_id)
    REFERENCES users(id)
    ON DELETE RESTRICT;