# Supabase migrations

This directory tracks the Supabase database baseline for production migration.
The Go backend remains the only product data API; Supabase Data API is disabled
locally and app tables are protected with RLS plus revoked `anon` /
`authenticated` privileges.

## Source of truth

- `backend/sql/schema.sql` remains the sqlc schema source.
- `backend/sql/migrations/*.sql` documents legacy incremental migrations.
- `supabase/migrations/20260702000000_baseline.sql` is the fresh Supabase
  baseline that combines the current schema and the security posture needed
  before cutover.

Do not delete Firebase/GCP/Terraform assets during migration. They remain the
rollback path until Vercel + Supabase has passed preview and production smoke
tests.

## Local check

```bash
supabase db reset
```

This starts the local Supabase stack and applies migrations to an empty local
database. If Docker is not running, use `psql` against a temporary PostgreSQL
database to validate the SQL syntax instead.

## Remote apply

Applying this to the real Supabase project changes database schema and
privileges. Confirm the target project and backup state first.

1. Verify the remote database is empty or repair migration history deliberately.
2. Confirm `SHOW server_version;` is compatible with the migration.
3. Link the project with `supabase link --project-ref <project-ref>`.
4. Run `supabase db push`.
5. Open Supabase Security Advisor and confirm no critical Data API / RLS issues
   remain for app tables.

Do not commit database passwords, pooler URLs with passwords, service role keys,
or OAuth client secrets.
