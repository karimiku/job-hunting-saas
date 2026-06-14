-- DDD boundary hardening: enforce repository-level uniqueness in PostgreSQL.

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

ALTER TABLE inbox_clips
    ADD CONSTRAINT inbox_clips_user_id_url_key UNIQUE (user_id, url);

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

ALTER TABLE company_aliases
    ADD CONSTRAINT company_aliases_user_id_company_id_alias_key UNIQUE (user_id, company_id, alias);
