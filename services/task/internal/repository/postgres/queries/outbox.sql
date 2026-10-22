-- name: InsertOutboxEvent :exec
INSERT INTO outbox_events (
    event_id,
    subject,
    payload
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(subject),
    sqlc.arg(payload)::jsonb
);

-- name: ListPendingOutboxEvents :many
SELECT *
FROM outbox_events
WHERE published_at IS NULL
  AND next_attempt_at <= NOW()
ORDER BY next_attempt_at, created_at, event_id
LIMIT sqlc.arg(batch_size)::integer;

-- name: MarkOutboxEventPublished :execrows
UPDATE outbox_events
SET
    published_at = NOW(),
    attempts = attempts + 1,
    last_error = ''
WHERE event_id = sqlc.arg(event_id)
  AND published_at IS NULL;

-- name: RecordOutboxEventFailure :execrows
UPDATE outbox_events
SET
    attempts = attempts + 1,
    next_attempt_at = sqlc.arg(next_attempt_at)::timestamptz,
    last_error = sqlc.arg(last_error)
WHERE event_id = sqlc.arg(event_id)
  AND published_at IS NULL;