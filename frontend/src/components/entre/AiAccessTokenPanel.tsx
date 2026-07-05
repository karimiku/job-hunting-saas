"use client";

import { useMemo, useState } from "react";
import { useActionState } from "react";
import { useFormStatus } from "react-dom";
import { Check, Copy, EyeOff, HelpCircle, KeyRound, ShieldCheck, Trash2, X } from "lucide-react";
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
  loadError,
  tokens,
}: {
  loadError?: string;
  tokens: AiAccessTokenResponse[];
}) {
  const [createState, createAction] = useActionState(
    createAiAccessTokenAction,
    CREATE_INITIAL,
  );
  const [modalOpen, setModalOpen] = useState(false);
  const activeTokens = useMemo(
    () => tokens.filter((token) => !token.revokedAt),
    [tokens],
  );

  return (
    <section className="min-w-0 overflow-hidden rounded-xl border border-line bg-surface p-5">
      <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div className="flex items-center gap-2">
          <span className="grid h-8 w-8 place-items-center rounded-lg bg-sage-soft text-sage">
            <KeyRound size={17} aria-hidden />
          </span>
          <div>
            <h2 className="text-[13px] font-extrabold">AI連携トークン</h2>
            <p className="mt-0.5 text-[12px] text-ink-3">
              Claude / Codex / MCP
            </p>
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          <span className="rounded-md bg-sage-wash px-2 py-1 font-mono text-[12px] font-bold text-sage">
            {activeTokens.length}
          </span>
          <button
            type="button"
            onClick={() => setModalOpen(true)}
            className="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg border border-line bg-white px-3 text-[12px] font-bold text-ink-2 transition-colors hover:bg-line-2"
          >
            <HelpCircle size={14} aria-hidden />
            仕組み
          </button>
        </div>
      </div>

      {loadError ? (
        <p role="alert" className="mb-4 rounded-md bg-pink/40 px-3 py-2 text-[12px] font-semibold text-pink-deep">
          {loadError}
        </p>
      ) : (
        <form action={createAction} className="mb-4 grid min-w-0 gap-2 md:grid-cols-[minmax(0,1fr)_auto]">
          <label className="min-w-0">
            <span className="sr-only">トークン名（任意）</span>
            <input
              name="name"
              defaultValue={createState.values?.name ?? ""}
              placeholder="名前なしでも作成できます"
              className="h-10 w-full rounded-lg border border-line bg-white px-3 text-[12px] font-semibold outline-none transition-colors focus:border-sage"
              maxLength={80}
            />
          </label>
          <CreateButton />
        </form>
      )}

      {createState.error && (
        <p role="alert" className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[12px] font-semibold">
          {createState.error}
        </p>
      )}

      {createState.token && (
        <IssuedToken token={createState.token} name={createState.accessToken?.name} />
      )}

      <div className="mt-4 min-w-0 divide-y divide-line rounded-lg border border-line bg-white">
        {tokens.length === 0 ? (
          <div className="px-3 py-4 text-[12px] font-semibold text-ink-3">
            {loadError ? "読み込み待ちです" : "まだありません"}
          </div>
        ) : (
          tokens.map((token) => <TokenRow key={token.id} token={token} />)
        )}
      </div>

      {modalOpen && <TokenHelpModal onClose={() => setModalOpen(false)} />}
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
  const [hiddenToken, setHiddenToken] = useState<string | null>(null);
  const visible = hiddenToken !== token;

  if (!visible) {
    return (
      <div className="mb-3 flex min-w-0 items-center justify-between gap-3 rounded-lg border border-sage/25 bg-sage-wash px-3 py-2">
        <span className="min-w-0 truncate text-[12px] font-extrabold text-sage">
          トークンをコピーしました
        </span>
        <span className="shrink-0 rounded-md bg-white px-2 py-1 text-[12px] font-bold text-ink-3">
          再表示不可
        </span>
      </div>
    );
  }

  const hideSoon = () => {
    window.setTimeout(() => setHiddenToken(token), 800);
  };

  return (
    <div className="mb-3 min-w-0 overflow-hidden rounded-lg border border-sage/25 bg-sage-wash p-3">
      <div className="mb-2 flex min-w-0 flex-wrap items-start justify-between gap-2">
        <div className="flex min-w-0 items-start gap-2">
          <ShieldCheck className="mt-0.5 shrink-0 text-sage" size={16} aria-hidden />
          <div className="min-w-0">
            <span className="block min-w-0 truncate text-[12px] font-extrabold text-sage">
              {name ?? "作成済み"}
            </span>
            <p className="mt-1 text-[12px] font-bold text-ink-3">
              この値は今だけ表示されます。
            </p>
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-1.5">
          <CopyButton text={token} label="コピーして隠す" onCopied={hideSoon} />
          <button
            type="button"
            onClick={() => setHiddenToken(token)}
            className="inline-flex h-8 items-center gap-1 rounded-md border border-line bg-white px-2 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            <EyeOff size={13} aria-hidden />
            隠す
          </button>
        </div>
      </div>
      <code className="block w-full max-w-full overflow-x-auto whitespace-nowrap rounded-md bg-white px-3 py-2 font-mono text-[12px] font-bold text-ink">
        {token}
      </code>
    </div>
  );
}

function CopyButton({
  label = "コピー",
  onCopied,
  text,
}: {
  label?: string;
  onCopied?: () => void;
  text: string;
}) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      type="button"
      onClick={async () => {
        await navigator.clipboard.writeText(text);
        setCopied(true);
        onCopied?.();
        window.setTimeout(() => setCopied(false), 1400);
      }}
      className="inline-flex h-8 shrink-0 items-center gap-1 whitespace-nowrap rounded-md border border-line bg-white px-2 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
    >
      {copied ? <Check size={13} aria-hidden /> : <Copy size={13} aria-hidden />}
      {copied ? "完了" : label}
    </button>
  );
}

function TokenHelpModal({ onClose }: { onClose: () => void }) {
  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-ink/40 px-4 py-6 backdrop-blur-sm" role="presentation">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="ai-token-help-title"
        className="max-h-[90vh] w-full max-w-[560px] overflow-y-auto rounded-xl border border-line bg-surface shadow-2xl"
      >
        <div className="flex items-start justify-between gap-3 border-b border-line px-5 py-4">
          <div>
            <h3 id="ai-token-help-title" className="text-[15px] font-extrabold">
              AI連携トークンについて
            </h3>
            <p className="mt-1 text-[12px] text-ink-3">設定コマンドではなく、tokenの扱いだけをここで確認できます。</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="grid h-8 w-8 shrink-0 place-items-center rounded-lg border border-line text-ink-2 transition hover:bg-line-2"
            aria-label="閉じる"
          >
            <X size={15} aria-hidden />
          </button>
        </div>
        <div className="space-y-4 px-5 py-4 text-[12px] leading-6 text-ink-2">
          <section>
            <h4 className="mb-1 text-[12px] font-extrabold text-ink">何に使う？</h4>
            <p>CodexなどのMCPクライアントからEntré APIへ接続するための鍵です。</p>
          </section>
          <section>
            <h4 className="mb-1 text-[12px] font-extrabold text-ink">保存されるもの</h4>
            <p>サーバーにはtoken全文ではなくhashだけを保存します。全文は発行直後の画面で一度だけ確認できます。</p>
          </section>
          <section>
            <h4 className="mb-1 text-[12px] font-extrabold text-ink">設定するとき</h4>
            <p>
              MCP wrapper側の環境変数に <code className="rounded bg-line-2 px-1 py-0.5">ENTRE_API_TOKEN</code>{" "}
              として貼り付けます。base URL は <code className="rounded bg-line-2 px-1 py-0.5">ENTRE_API_BASE_URL</code>{" "}
              に設定します。
            </p>
          </section>
          <section>
            <h4 className="mb-1 text-[12px] font-extrabold text-ink">失くした・漏れたとき</h4>
            <p>一覧から失効して、新しいtokenを発行してください。失効済みtokenではAPIに接続できません。</p>
          </section>
        </div>
      </div>
    </div>
  );
}

function TokenRow({ token }: { token: AiAccessTokenResponse }) {
  const [state, formAction] = useActionState(
    revokeAiAccessTokenAction,
    REVOKE_INITIAL,
  );
  const revoked = Boolean(token.revokedAt || state.revokedId === token.id);

  return (
    <div className="grid min-w-0 gap-3 px-3 py-3 md:grid-cols-[minmax(0,1fr)_auto] md:items-center">
      <div className="min-w-0">
        <div className="flex min-w-0 items-center gap-2">
          <span className="truncate text-[12px] font-extrabold">{token.name}</span>
          <span
            className={`rounded-md px-1.5 py-0.5 text-[12px] font-bold ${
              revoked ? "bg-line-2 text-ink-3" : "bg-sage-soft text-sage"
            }`}
          >
            {revoked ? "失効" : "有効"}
          </span>
        </div>
        <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-[12px] font-semibold text-ink-3">
          <span className="font-mono">{token.tokenPrefix}...</span>
          <span>再表示不可</span>
          <span>作成 {formatDate(token.createdAt)}</span>
          <span>利用 {token.lastUsedAt ? formatDate(token.lastUsedAt) : "-"}</span>
        </div>
        {state.error && (
          <p role="alert" className="mt-2 text-[12px] font-bold text-pink-deep">
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
      className="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg border border-line bg-surface px-3 text-[12px] font-bold text-ink-3 transition-colors enabled:hover:border-pink-deep enabled:hover:text-pink-deep disabled:opacity-60"
    >
      <Trash2 size={13} aria-hidden />
      {pending ? "失効中" : "失効"}
    </button>
  );
}

function formatDate(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";

  const parts = new Intl.DateTimeFormat("ja-JP", {
    timeZone: "Asia/Tokyo",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).formatToParts(date);
  const get = (type: Intl.DateTimeFormatPartTypes) =>
    parts.find((part) => part.type === type)?.value ?? "";

  return `${get("month")}/${get("day")} ${get("hour")}:${get("minute")}`;
}
