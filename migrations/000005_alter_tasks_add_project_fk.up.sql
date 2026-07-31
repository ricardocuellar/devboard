
ALTER TABLE tasks
    ADD CONSTRAINT fk_tasks_project_id
    FOREIGN KEY (project_id)
    REFERENCES projects(id)
    ON DELETE CASCADE;