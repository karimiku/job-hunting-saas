ALTER TABLE inbox_clips
    ADD COLUMN IF NOT EXISTS content_text TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS selection_flows (
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

CREATE TABLE IF NOT EXISTS selection_stages (
    id            UUID        PRIMARY KEY,
    flow_id       UUID        NOT NULL REFERENCES selection_flows(id) ON DELETE CASCADE,
    position      INTEGER     NOT NULL CHECK (position > 0),
    stage_kind    stage_kind  NOT NULL,
    stage_label   TEXT        NOT NULL,
    evidence_text TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT selection_stages_flow_id_position_key UNIQUE (flow_id, position)
);

CREATE INDEX IF NOT EXISTS idx_selection_flows_entry_id ON selection_flows(entry_id);
CREATE INDEX IF NOT EXISTS idx_selection_stages_flow_id_position ON selection_stages(flow_id, position);
