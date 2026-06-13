"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import { Check, Copy, HelpCircle, KeyRound, Loader2, ShieldCheck, Trash2, X } from "lucide-react";
import {
  type AIAccessTokenItem,
  createAIAccessToken,
  listAIAccessTokens,
  revokeAIAccessToken,
} from "@/lib/api/aiAccessTokens";

interface CreatedToken {
  token: string;
  item: AIAccessTokenItem;
}

export function AIAccessTokenPanel() {
  const [tokens, setTokens] = useState<AIAccessTokenItem[]>([]);
  const [name, setName] = useState("");
  const [created, setCreated] = useState<CreatedToken | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const activeTokens = useMemo(() => tokens.filter((token) => !token.revokedAt), [tokens]);

  useEffect(() => {
    let alive = true;
    listAIAccessTokens()
      .then((items) => {
        if (alive) setTokens(items);
      })
      .catch(() => {
        if (alive) setError("トークン一覧を取得できませんでした");
      })
      .finally(() => {
        if (alive) setLoading(false);
      });
    return () => {
      alive = false;
    };
  }, []);

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setCopied(false);
    try {
      const res = await createAIAccessToken(name);
      setCreated(res);
      setTokens((prev) => [res.item, ...prev]);
      setName("");
    } catch {
      setError("トークンを発行できませんでした");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleCopy() {
    if (!created) return;
    await navigator.clipboard.writeText(created.token);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1800);
  }

  async function handleRevoke(tokenId: string) {
    setError(null);
    try {
      await revokeAIAccessToken(tokenId);
      const now = new Date().toISOString();
      setTokens((prev) =>
        prev.map((token) => (token.id === tokenId ? { ...token, revokedAt: now } : token)),
      );
    } catch {
      setError("トークンを失効できませんでした");
    }
  }

  return (
    <section className="rounded-xl border border-line bg-surface p-5">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div className="min-w-0">
          <div className="mb-2 flex items-center gap-2">
            <div className="grid h-8 w-8 place-items-center rounded-lg bg-sage-soft text-sage">
              <KeyRound size={17} aria-hidden />
            </div>
            <h2 className="text-[13px] font-extrabold">AI連携トークン</h2>
          </div>
          <p className="max-w-[560px] text-[12px] leading-6 text-ink-2">
            MCP wrapper などから Entré API に接続するための token です。発行時だけ全文を表示します。
          </p>
        </div>
        <button
          type="button"
          onClick={() => setModalOpen(true)}
          className="inline-flex h-9 shrink-0 items-center justify-center gap-1.5 rounded-lg border border-line bg-white px-3 text-[11px] font-bold text-ink-2 transition-colors hover:bg-line-2"
        >
          <HelpCircle size={14} aria-hidden />
          仕組み
        </button>
      </div>

      <form onSubmit={handleCreate} className="mt-5 flex flex-col gap-3 rounded-lg border border-line bg-white p-3">
        <label className="text-[11px] font-bold text-ink-2" htmlFor="ai-token-name">
          名前（任意）
        </label>
        <div className="flex flex-col gap-2 sm:flex-row">
          <input
            id="ai-token-name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            maxLength={80}
            placeholder="名前なしでも作成できます"
            className="min-h-10 flex-1 rounded-lg border border-line bg-surface px-3 text-[13px] outline-none transition focus:border-sage"
          />
          <button
            type="submit"
            disabled={submitting}
            className="inline-flex h-10 shrink-0 items-center justify-center gap-2 rounded-lg bg-sage px-4 text-[12px] font-extrabold text-white transition hover:bg-sage-mid disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitting ? <Loader2 size={15} className="animate-spin" aria-hidden /> : <KeyRound size={15} aria-hidden />}
            発行
          </button>
        </div>
      </form>

      {created && (
        <div className="mt-4 rounded-lg border border-sage/30 bg-sage-soft/70 p-3">
          <div className="flex items-start gap-2">
            <ShieldCheck className="mt-0.5 shrink-0 text-sage" size={17} aria-hidden />
            <div className="min-w-0 flex-1">
              <p className="text-[12px] font-extrabold">このtokenは今だけ表示されます</p>
              <p className="mt-1 text-[11px] leading-5 text-ink-2">
                閉じたり再読み込みすると全文は確認できません。必要なら失効して作り直してください。
              </p>
            </div>
          </div>
          <div className="mt-3 flex flex-col gap-2 sm:flex-row">
            <code className="min-w-0 flex-1 overflow-x-auto rounded-md border border-line bg-white px-3 py-2 font-mono text-[11px] text-ink">
              {created.token}
            </code>
            <button
              type="button"
              onClick={handleCopy}
              className="inline-flex h-9 shrink-0 items-center justify-center gap-1.5 rounded-lg bg-ink px-3 text-[11px] font-bold text-white transition hover:bg-ink/85"
            >
              {copied ? <Check size={14} aria-hidden /> : <Copy size={14} aria-hidden />}
              {copied ? "コピー済み" : "コピー"}
            </button>
          </div>
        </div>
      )}

      {error && <p className="mt-3 rounded-lg bg-red-50 px-3 py-2 text-[12px] font-bold text-red-700">{error}</p>}

      <div className="mt-5">
        <div className="mb-2 flex items-center justify-between gap-3">
          <p className="text-[11px] font-bold text-ink-2">発行済み</p>
          <span className="rounded-md bg-line-2 px-2 py-1 text-[10px] font-bold text-ink-3">
            active {activeTokens.length}
          </span>
        </div>
        {loading ? (
          <div className="flex h-20 items-center justify-center rounded-lg border border-dashed border-line text-ink-3">
            <Loader2 size={18} className="animate-spin" aria-hidden />
          </div>
        ) : tokens.length === 0 ? (
          <p className="rounded-lg border border-dashed border-line px-3 py-5 text-center text-[12px] text-ink-3">
            まだ token はありません。
          </p>
        ) : (
          <ul className="flex flex-col gap-2">
            {tokens.map((token) => {
              const revoked = Boolean(token.revokedAt);
              return (
                <li
                  key={token.id}
                  className="flex flex-col gap-2 rounded-lg border border-line bg-white px-3 py-3 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="max-w-[220px] truncate text-[12px] font-extrabold">{token.name}</span>
                      <code className="rounded bg-line-2 px-1.5 py-0.5 font-mono text-[10px] text-ink-3">
                        {token.tokenPreview}
                      </code>
                      <span
                        className={`rounded px-1.5 py-0.5 text-[9px] font-bold ${
                          revoked ? "bg-line-2 text-ink-3" : "bg-sage-soft text-sage"
                        }`}
                      >
                        {revoked ? "revoked" : "active"}
                      </span>
                    </div>
                    <p className="mt-1 text-[10px] text-ink-3">
                      作成 {formatDate(token.createdAt)}
                      {token.lastUsedAt ? ` / 最終利用 ${formatDate(token.lastUsedAt)}` : ""}
                    </p>
                  </div>
                  {!revoked && (
                    <button
                      type="button"
                      onClick={() => handleRevoke(token.id)}
                      className="inline-flex h-8 shrink-0 items-center justify-center gap-1.5 rounded-lg border border-line px-2.5 text-[11px] font-bold text-ink-2 transition hover:bg-line-2"
                    >
                      <Trash2 size={13} aria-hidden />
                      失効
                    </button>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </div>

      {modalOpen && <TokenHelpModal onClose={() => setModalOpen(false)} />}
    </section>
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
            <p className="mt-1 text-[11px] text-ink-3">設定コマンドではなく、tokenの扱いだけをここで確認できます。</p>
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

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("ja-JP", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}
