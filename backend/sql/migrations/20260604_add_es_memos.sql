CREATE TABLE IF NOT EXISTS es_memos (
    id         UUID        PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entry_id   UUID        REFERENCES entries(id) ON DELETE SET NULL,
    category   TEXT        NOT NULL DEFAULT 'general',
    title      TEXT        NOT NULL,
    content    TEXT        NOT NULL,
    source     TEXT        NOT NULL DEFAULT 'mcp',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_es_memos_user_created_at ON es_memos(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_es_memos_user_entry ON es_memos(user_id, entry_id);
