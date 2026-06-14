-- DDD boundary hardening: enforce repository-level uniqueness in PostgreSQL.

-- Normalize legacy InboxClip rows before repository reconstruction starts
-- enforcing InboxClipTitle / InboxClipGuess value object rules.
WITH normalized AS (
    SELECT
        id,
        LEFT(
            COALESCE(
                NULLIF(regexp_replace(title, '^[[:space:]]+|[[:space:]]+$', '', 'g'), ''),
                NULLIF(regexp_replace(url, '^[[:space:]]+|[[:space:]]+$', '', 'g'), ''),
                'Untitled'
            ),
            512
        ) AS title,
        LEFT(regexp_replace(guess, '^[[:space:]]+|[[:space:]]+$', '', 'g'), 256) AS guess
    FROM inbox_clips
)
UPDATE inbox_clips c
SET title = normalized.title,
    guess = normalized.guess
FROM normalized
WHERE c.id = normalized.id
  AND (
      c.title IS DISTINCT FROM normalized.title
      OR c.guess IS DISTINCT FROM normalized.guess
  );

-- InboxClip registration keeps the newest duplicate because the previous
-- lookup policy returned the latest captured clip for (user_id, url).
DELETE FROM inbox_clips a
USING inbox_clips b
WHERE a.user_id = b.user_id
  AND a.url = b.url
  AND (
      a.captured_at < b.captured_at
      OR (a.captured_at = b.captured_at AND a.id::text < b.id::text)
  );

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'inbox_clips'::regclass
          AND conname = 'inbox_clips_user_id_url_key'
    ) THEN
        ALTER TABLE inbox_clips
            ADD CONSTRAINT inbox_clips_user_id_url_key UNIQUE (user_id, url);
    END IF;
END $$;

-- CompanyAlias is immutable; keep the oldest alias row for duplicate
-- (user_id, company_id, alias) groups.
DELETE FROM company_aliases a
USING company_aliases b
WHERE a.user_id = b.user_id
  AND a.company_id = b.company_id
  AND a.alias = b.alias
  AND (
      a.created_at > b.created_at
      OR (a.created_at = b.created_at AND a.id::text > b.id::text)
  );

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conrelid = 'company_aliases'::regclass
          AND conname = 'company_aliases_user_id_company_id_alias_key'
    ) THEN
        ALTER TABLE company_aliases
            ADD CONSTRAINT company_aliases_user_id_company_id_alias_key UNIQUE (user_id, company_id, alias);
    END IF;
END $$;
