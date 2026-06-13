CREATE TABLE IF NOT EXISTS ai_access_tokens (
    id            UUID        PRIMARY KEY,
    user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    token_hash    TEXT        NOT NULL UNIQUE,
    token_preview TEXT        NOT NULL,
    last_used_at  TIMESTAMPTZ,
    revoked_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ai_access_tokens_user_created_at
    ON ai_access_tokens(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_ai_access_tokens_active_hash
    ON ai_access_tokens(token_hash)
    WHERE revoked_at IS NULL;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'runtime_user') THEN
        EXECUTE 'GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE ai_access_tokens TO runtime_user';
    END IF;
END $$;
