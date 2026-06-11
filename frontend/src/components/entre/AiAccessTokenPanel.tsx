"use client";

import { useMemo, useState } from "react";
import { useActionState } from "react";
import { useFormStatus } from "react-dom";
import { Check, Copy, KeyRound, Trash2 } from "lucide-react";
import {
  createAiAccessTokenAction,
  revokeAiAccessTokenAction,
  type CreateAiAccessTokenState,
  type RevokeAiAccessTokenState,
} from "@/app/profile/actions";
import type { AiAccessTokenResponse } from "@/lib/api/aiTokens";

const CREATE_INITIAL: CreateAiAccessTokenState = {};
const REVOKE_INITIAL: RevokeAiAccessTokenState = {};

export function AiAccessTokenPanel({
  tokens,
}: {
  tokens: AiAccessTokenResponse[];
}) {
  const [createState, createAction] = useActionState(
    createAiAccessTokenAction,
    CREATE_INITIAL,
  );
  const activeTokens = useMemo(
    () => tokens.filter((token) => !token.revokedAt),
    [tokens],
  );

  return (
    <section className="rounded-xl border border-line bg-surface p-5">
      <div className="mb-4 flex items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <span className="grid h-8 w-8 place-items-center rounded-lg bg-sage-soft text-sage">
            <KeyRound size={17} aria-hidden />
          </span>
          <div>
            <h2 className="text-[13px] font-extrabold">AI連携トークン</h2>
            <p className="mt-0.5 text-[11px] text-ink-3">
              Claude / Codex / MCP
            </p>
          </div>
        </div>
        <span className="rounded-md bg-sage-wash px-2 py-1 font-mono text-[10px] font-bold text-sage">
          {activeTokens.length}
        </span>
      </div>

      <form action={createAction} className="mb-4 grid gap-2 md:grid-cols-[1fr_auto]">
        <label className="min-w-0">
          <span className="sr-only">トークン名</span>
          <input
            name="name"
            defaultValue={createState.values?.name ?? ""}
            placeholder="Claude Desktop"
            className="h-10 w-full rounded-lg border border-line bg-white px-3 text-[12px] font-semibold outline-none transition-colors focus:border-sage"
            maxLength={80}
          />
        </label>
        <CreateButton />
      </form>

      {createState.error && (
        <p role="alert" className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[11px] font-semibold">
          {createState.error}
        </p>
      )}

      {createState.token && (
        <IssuedToken token={createState.token} name={createState.accessToken?.name} />
      )}

      <div className="mt-4 divide-y divide-line rounded-lg border border-line bg-white">
        {tokens.length === 0 ? (
          <div className="px-3 py-4 text-[12px] font-semibold text-ink-3">
            まだありません
          </div>
        ) : (
          tokens.map((token) => <TokenRow key={token.id} token={token} />)
        )}
      </div>
    </section>
  );
}

function CreateButton() {
  const { pending } = useFormStatus();
  return (
    <button
      type="submit"
      disabled={pending}
      className="inline-flex h-10 items-center justify-center gap-1.5 rounded-lg bg-sage px-4 text-[12px] font-extrabold text-white transition-colors enabled:hover:bg-sage-dark disabled:opacity-60"
    >
      <KeyRound size={14} aria-hidden />
      {pending ? "作成中" : "作成"}
    </button>
  );
}

function IssuedToken({ token, name }: { token: string; name?: string }) {
  return (
    <div className="mb-3 rounded-lg border border-sage/25 bg-sage-wash p-3">
      <div className="mb-2 flex items-center justify-between gap-2">
        <span className="text-[11px] font-extrabold text-sage">
          {name ?? "作成済み"}
        </span>
        <CopyTokenButton token={token} />
      </div>
      <code className="block max-w-full overflow-x-auto rounded-md bg-white px-3 py-2 font-mono text-[11px] font-bold text-ink">
        {token}
      </code>
    </div>
  );
}

function CopyTokenButton({ token }: { token: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      type="button"
      onClick={async () => {
        await navigator.clipboard.writeText(token);
        setCopied(true);
        window.setTimeout(() => setCopied(false), 1400);
      }}
      className="inline-flex h-8 items-center gap-1 rounded-md border border-line bg-white px-2 text-[11px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
    >
      {copied ? <Check size={13} aria-hidden /> : <Copy size={13} aria-hidden />}
      {copied ? "完了" : "コピー"}
    </button>
  );
}

function TokenRow({ token }: { token: AiAccessTokenResponse }) {
  const [state, formAction] = useActionState(
    revokeAiAccessTokenAction,
    REVOKE_INITIAL,
  );
  const revoked = Boolean(token.revokedAt || state.revokedId === token.id);

  return (
    <div className="grid gap-3 px-3 py-3 md:grid-cols-[1fr_auto] md:items-center">
      <div className="min-w-0">
        <div className="flex min-w-0 items-center gap-2">
          <span className="truncate text-[12px] font-extrabold">{token.name}</span>
          <span
            className={`rounded-md px-1.5 py-0.5 text-[9px] font-bold ${
              revoked ? "bg-line-2 text-ink-3" : "bg-sage-soft text-sage"
            }`}
          >
            {revoked ? "失効" : "有効"}
          </span>
        </div>
        <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-[10px] font-semibold text-ink-3">
          <span className="font-mono">{token.tokenPrefix}...</span>
          <span>作成 {formatDate(token.createdAt)}</span>
          <span>利用 {token.lastUsedAt ? formatDate(token.lastUsedAt) : "-"}</span>
        </div>
        {state.error && (
          <p role="alert" className="mt-2 text-[10px] font-bold text-pink-deep">
            {state.error}
          </p>
        )}
      </div>
      {!revoked && (
        <form action={formAction}>
          <input type="hidden" name="tokenId" value={token.id} />
          <RevokeButton name={token.name} />
        </form>
      )}
    </div>
  );
}

function RevokeButton({ name }: { name: string }) {
  const { pending } = useFormStatus();
  return (
    <button
      type="submit"
      disabled={pending}
      aria-label={`${name} を失効`}
      onClick={(e) => {
        if (!window.confirm("このトークンを失効しますか？")) {
          e.preventDefault();
        }
      }}
      className="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg border border-line bg-surface px-3 text-[11px] font-bold text-ink-3 transition-colors enabled:hover:border-pink-deep enabled:hover:text-pink-deep disabled:opacity-60"
    >
      <Trash2 size={13} aria-hidden />
      {pending ? "失効中" : "失効"}
    </button>
  );
}

function formatDate(value: string): string {
  return new Intl.DateTimeFormat("ja-JP", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));
}
