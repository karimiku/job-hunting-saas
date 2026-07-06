# Vercel + Supabase 移行 ADR

更新日: 2026-07-01

## Status

Accepted for phased migration.

## Context

現状は Next.js frontend、Go backend、PostgreSQL、Firebase Auth、Cloud Run/GCP 系の構成を前提にしている。Vercel の Dockerfile/Container Images と Services により、Next.js と Go API を Vercel 側へ寄せられる余地が出た。Auth と DB を Supabase に寄せると、プロジェクトの主要な外部依存は Vercel と Supabase に集約できる。

## Decision

移行先の基本方針は Vercel + Supabase とする。

- Frontend は Vercel の Next.js deployment を前提にする。
- Go backend は Vercel Services / Container Images を候補にし、Cloud Run と同等の API 互換性を確認する。
- Auth は Firebase Auth から Supabase Auth へ移す。
- Database は Supabase Postgres を移行先にする。
- 既存の Clean Architecture、OpenAPI SSOT、sqlc 方針は維持する。
- Firebase/GCP の削除は最後に行い、Supabase Auth と Vercel backend の parity が取れるまでは併存を許容する。

## Migration Order

1. 移行判断と依存関係を issue / docs に固定する。
2. Supabase project、env、DB 接続方式、破壊的 test guard を整える。
3. Supabase Auth の frontend SSR/session 基盤を追加する。
4. Go backend で Supabase JWT を Bearer token として検証する。
5. `users` / `external_identities` を Supabase subject/provider に対応させる。
6. Frontend server actions / Chrome Extension から Bearer token 付き API 呼び出しへ寄せる。
7. Vercel Services + Container Images で Go backend を PoC deploy する。
8. Routing、scale-down、body size、cold start、timeout の制約を確認する。
9. Cloud Run / Firebase / GCP IaC を削除する。

## Dependency Map

- #213: 親 epic。移行全体の tracking。
- #214: この ADR。移行判断の基準。
- #215: Supabase project/env 整備。
- #216: Frontend Supabase SSR。
- #217: Go backend の Supabase JWT verifier。
- #218: Frontend API 呼び出しの Bearer 化。
- #219: `users` / `external_identities` の subject/provider 移行。
- #220: Vercel Services monorepo routing。
- #221: Go backend の Vercel Container Images PoC。
- #222: `/backend` proxy/path mismatch の解消。
- #223: Vercel runtime 制約の検証。
- #224: Supabase migration baseline。
- #225: Data API / RLS / DB role 設計。
- #226: Supabase pooler / pgxpool 接続設計。
- #227: Remote Supabase DB を壊さない test guard。
- #228: Chrome Extension の auth/API host 移行。
- #229: Backend CD を Vercel に移す。
- #230: GCP/Firebase docs/IaC cleanup。

## Risks

- Vercel Container Images / Services は beta 機能を含むため、制約変更や unsupported capability が起き得る。
- Vercel の scale-down により、Cloud Run と異なる cold start / background work 特性が出る。
- Supabase pooler の transaction mode は prepared statement 非対応なので、pgxpool の接続設定と互換性確認が必要。
- Supabase Auth 移行中は Firebase session と Supabase session の二重運用が発生し得る。
- Supabase の exposed schema では RLS と Data API 権限を明示しないと意図しない公開/非公開が起きる。

## Rollback

移行完了までは Cloud Run / Firebase Auth を rollback path として残す。Vercel backend と Supabase Auth の本番 parity が取れるまで、Firebase/GCP の削除 PR は作らない。

