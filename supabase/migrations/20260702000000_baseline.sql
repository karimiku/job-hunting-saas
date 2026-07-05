-- job-hunting-saas Supabase baseline migration.
--
-- This migration is for an empty Supabase database. Do not apply it to a
-- database that has already been initialized with backend/sql/schema.sql and
-- backend/sql/migrations/*.sql unless the migration history has been repaired
-- deliberately.
--
-- backend/sql/schema.sql remains the sqlc schema source of truth. Keep this
-- baseline synchronized when the production schema changes.

CREATE TYPE entry_status AS ENUM (
    'in_progress', 'offered', 'accepted', 'rejected', 'withdrawn'
);

CREATE TYPE stage_kind AS ENUM (
    'application', 'document', 'test', 'interview', 'group', 'offer', 'other'
);

CREATE TYPE task_type AS ENUM ('deadline', 'schedule');

CREATE TYPE task_status AS ENUM ('todo', 'done');

CREATE TYPE auth_provider AS ENUM ('google', 'supabase');

CREATE TABLE users (
    id         UUID        PRIMARY KEY,
    email      TEXT        NOT NULL UNIQUE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE companies (
    id         UUID        PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    memo       TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE entries (
    id          UUID         PRIMARY KEY,
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id  UUID         NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    route       TEXT         NOT NULL,
    source      TEXT         NOT NULL,
    source_url  TEXT         NOT NULL DEFAULT '',
    status      entry_status NOT NULL DEFAULT 'in_progress',
    stage_kind  stage_kind   NOT NULL DEFAULT 'application',
    stage_label TEXT         NOT NULL DEFAULT '',
    memo        TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE tasks (
    id         UUID        PRIMARY KEY,
    entry_id   UUID        NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
    title      TEXT        NOT NULL,
    task_type  task_type   NOT NULL,
    due_date   TIMESTAMPTZ,
    status     task_status NOT NULL DEFAULT 'todo',
    notify     BOOLEAN     NOT NULL DEFAULT false,
    memo       TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE stage_histories (
    id          UUID        PRIMARY KEY,
    entry_id    UUID        NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
    stage_kind  stage_kind  NOT NULL,
    stage_label TEXT        NOT NULL,
    note        TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE company_aliases (
    id         UUID        PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID        NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    alias      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT company_aliases_user_id_company_id_alias_key UNIQUE (user_id, company_id, alias)
);

CREATE TABLE external_identities (
    id         UUID          PRIMARY KEY,
    user_id    UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider   auth_provider NOT NULL,
    subject    TEXT          NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    UNIQUE (provider, subject)
);

CREATE TABLE password_credentials (
    id            UUID        PRIMARY KEY,
    user_id       UUID        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inbox_clips (
    id           UUID        PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url          TEXT        NOT NULL,
    title        TEXT        NOT NULL,
    source       TEXT        NOT NULL,
    guess        TEXT        NOT NULL DEFAULT '',
    content_text TEXT        NOT NULL DEFAULT '',
    captured_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT inbox_clips_user_id_url_key UNIQUE (user_id, url)
);

CREATE TABLE es_memos (
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

CREATE TABLE ai_access_tokens (
    id           UUID        PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    token_hash   TEXT        NOT NULL UNIQUE,
    token_prefix TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_used_at TIMESTAMPTZ,
    revoked_at   TIMESTAMPTZ
);

CREATE TABLE selection_flows (
    id                     UUID        PRIMARY KEY,
    entry_id               UUID        NOT NULL UNIQUE REFERENCES entries(id) ON DELETE CASCADE,
    source                 TEXT        NOT NULL,
    current_stage_position INTEGER     NOT NULL DEFAULT 1 CHECK (current_stage_position > 0),
    confidence             INTEGER     CHECK (confidence IS NULL OR (confidence >= 0 AND confidence <= 100)),
    inbox_clip_id          UUID        REFERENCES inbox_clips(id) ON DELETE SET NULL,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT selection_flows_source_check CHECK (source IN ('template', 'manual', 'ai_inbox', 'ai_paste'))
);

CREATE TABLE selection_stages (
    id            UUID        PRIMARY KEY,
    flow_id       UUID        NOT NULL REFERENCES selection_flows(id) ON DELETE CASCADE,
    position      INTEGER     NOT NULL CHECK (position > 0),
    stage_kind    stage_kind  NOT NULL,
    stage_label   TEXT        NOT NULL,
    evidence_text TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT selection_stages_flow_id_position_key UNIQUE (flow_id, position)
);

CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_entries_user_id ON entries(user_id);
CREATE INDEX idx_entries_user_id_status ON entries(user_id, status)
    WHERE status IN ('in_progress', 'offered');
CREATE INDEX idx_tasks_entry_id ON tasks(entry_id);
CREATE INDEX idx_tasks_due_date ON tasks(due_date)
    WHERE status = 'todo' AND due_date IS NOT NULL;
CREATE INDEX idx_stage_histories_entry_id ON stage_histories(entry_id);
CREATE INDEX idx_selection_flows_entry_id ON selection_flows(entry_id);
CREATE INDEX idx_selection_stages_flow_id_position ON selection_stages(flow_id, position);
CREATE INDEX idx_company_aliases_user_company ON company_aliases(user_id, company_id);
CREATE INDEX idx_inbox_clips_user_captured_at ON inbox_clips(user_id, captured_at DESC);
CREATE INDEX idx_es_memos_user_created_at ON es_memos(user_id, created_at DESC);
CREATE INDEX idx_es_memos_user_entry ON es_memos(user_id, entry_id);
CREATE INDEX idx_ai_access_tokens_user_created_at ON ai_access_tokens(user_id, created_at DESC);

-- Product data is served only through the Go API. Keep Supabase Data API roles
-- unable to read or mutate app tables, even if the API gateway is accidentally
-- enabled later.
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE stage_histories ENABLE ROW LEVEL SECURITY;
ALTER TABLE company_aliases ENABLE ROW LEVEL SECURITY;
ALTER TABLE external_identities ENABLE ROW LEVEL SECURITY;
ALTER TABLE password_credentials ENABLE ROW LEVEL SECURITY;
ALTER TABLE inbox_clips ENABLE ROW LEVEL SECURITY;
ALTER TABLE es_memos ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_access_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE selection_flows ENABLE ROW LEVEL SECURITY;
ALTER TABLE selection_stages ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    app_table text;
    app_tables text[] := ARRAY[
        'users',
        'companies',
        'entries',
        'tasks',
        'stage_histories',
        'company_aliases',
        'external_identities',
        'password_credentials',
        'inbox_clips',
        'es_memos',
        'ai_access_tokens',
        'selection_flows',
        'selection_stages'
    ];
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'anon') THEN
        EXECUTE 'REVOKE ALL ON SCHEMA public FROM anon';
        EXECUTE 'REVOKE ALL ON ALL TABLES IN SCHEMA public FROM anon';
        EXECUTE 'REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM anon';
        EXECUTE 'REVOKE ALL ON ALL FUNCTIONS IN SCHEMA public FROM anon';
        EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM anon';
        EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON SEQUENCES FROM anon';
        EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON FUNCTIONS FROM anon';
    END IF;

    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'authenticated') THEN
        EXECUTE 'REVOKE ALL ON SCHEMA public FROM authenticated';
        EXECUTE 'REVOKE ALL ON ALL TABLES IN SCHEMA public FROM authenticated';
        EXECUTE 'REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM authenticated';
        EXECUTE 'REVOKE ALL ON ALL FUNCTIONS IN SCHEMA public FROM authenticated';
        EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM authenticated';
        EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON SEQUENCES FROM authenticated';
        EXECUTE 'ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON FUNCTIONS FROM authenticated';
    END IF;

    -- Optional future hardening: create runtime_user out-of-band with a secret
    -- password, then rerun or reapply equivalent grants. Migration and runtime
    -- roles should be separate before production traffic is moved.
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'runtime_user') THEN
        EXECUTE 'GRANT USAGE ON SCHEMA public TO runtime_user';
        EXECUTE 'GRANT USAGE ON TYPE entry_status, stage_kind, task_type, task_status, auth_provider TO runtime_user';

        FOREACH app_table IN ARRAY app_tables LOOP
            EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.%I TO runtime_user', app_table);
            EXECUTE format('DROP POLICY IF EXISTS %I ON public.%I', app_table || '_runtime_all', app_table);
            EXECUTE format(
                'CREATE POLICY %I ON public.%I FOR ALL TO runtime_user USING (true) WITH CHECK (true)',
                app_table || '_runtime_all',
                app_table
            );
        END LOOP;
    END IF;
END $$;
