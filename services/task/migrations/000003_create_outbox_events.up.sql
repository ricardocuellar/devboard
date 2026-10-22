--up
CREATE TABLE outbox_events (
    event_id        TEXT PRIMARY KEY,
    subject         TEXT NOT NULL CHECK (subject <> ''),
    payload         JSONB NOT NULL
                    CHECK (jsonb_typeof(payload) = 'object'),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at    TIMESTAMPTZ,
    attempts        INTEGER NOT NULL DEFAULT 0
                    CHECK (attempts >= 0),
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_error      TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_outbox_events_pending
    ON outbox_events (next_attempt_at, created_at, event_id)
    WHERE published_at IS NULL;