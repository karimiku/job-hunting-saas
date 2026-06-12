CREATE TABLE IF NOT EXISTS ai_access_tokens (
    id           UUID        PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    token_hash   TEXT        NOT NULL UNIQUE,
    token_prefix TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ,
    revoked_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ai_access_tokens_user_created_at
    ON ai_access_tokens(user_id, created_at DESC);

ALTER TABLE ai_access_tokens ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'anon') THEN
        REVOKE ALL ON TABLE ai_access_tokens FROM anon;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'authenticated') THEN
        REVOKE ALL ON TABLE ai_access_tokens FROM authenticated;
    END IF;

    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'runtime_user') THEN
        GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE ai_access_tokens TO runtime_user;

        DROP POLICY IF EXISTS ai_access_tokens_runtime_all ON ai_access_tokens;
        CREATE POLICY ai_access_tokens_runtime_all
            ON ai_access_tokens
            FOR ALL
            TO runtime_user
            USING (true)
            WITH CHECK (true);
    END IF;
END $$;
