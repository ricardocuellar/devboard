-- name: RegisterProcessedEvent :execrows
INSERT INTO processed_events (event_id)
VALUES (sqlc.arg(event_id))
ON CONFLICT (event_id) DO NOTHING;

-- name: CreateNotification :one
INSERT INTO notifications (
    event_id,
    user_id,
    task_id,
    project_id,
    message
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(user_id),
    sqlc.arg(task_id),
    sqlc.arg(project_id),
    sqlc.arg(message)
)
RETURNING *;

-- name: ListNotificationsByUser :many
SELECT *
FROM notifications
WHERE user_id = sqlc.arg(user_id)
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(result_limit)::integer;