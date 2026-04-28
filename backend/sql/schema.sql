-- PostgreSQL スキーマ定義
-- job-hunting-saas 全テーブル一括定義

-- ============================================================
-- ENUM 型
-- ============================================================

CREATE TYPE entry_status AS ENUM (
    'in_progress', 'offered', 'accepted', 'rejected', 'withdrawn'
);

CREATE TYPE stage_kind AS ENUM (
    'application', 'document', 'test', 'interview', 'group', 'offer', 'other'
);

CREATE TYPE task_type AS ENUM ('deadline', 'schedule');

CREATE TYPE task_status AS ENUM ('todo', 'done');

CREATE TYPE auth_provider AS ENUM ('google');

-- ============================================================
-- テーブル（FK 依存順）
-- ============================================================

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
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
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
    id          UUID        PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url         TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    source      TEXT        NOT NULL,
    guess       TEXT        NOT NULL DEFAULT '',
    captured_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- インデックス
-- ============================================================

CREATE INDEX idx_companies_user_id ON companies(user_id);

CREATE INDEX idx_entries_user_id ON entries(user_id);

-- ダッシュボード用: オープンなエントリーのみ高速検索
CREATE INDEX idx_entries_user_id_status ON entries(user_id, status)
    WHERE status IN ('in_progress', 'offered');

CREATE INDEX idx_tasks_entry_id ON tasks(entry_id);

-- 締切通知用: 未完了かつ期日ありのタスクのみ
CREATE INDEX idx_tasks_due_date ON tasks(due_date)
    WHERE status = 'todo' AND due_date IS NOT NULL;

CREATE INDEX idx_stage_histories_entry_id ON stage_histories(entry_id);

CREATE INDEX idx_company_aliases_user_company ON company_aliases(user_id, company_id);

-- Inbox 一覧表示用: ユーザの直近クリップから降順で取得
CREATE INDEX idx_inbox_clips_user_captured_at ON inbox_clips(user_id, captured_at DESC);
