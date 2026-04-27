"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useEffect, useState } from "react";
import Link from "next/link";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { Mascot } from "@/components/entre/Mascot";
import { createCompany } from "@/lib/api/companies";
import { createEntry } from "@/lib/api/entries";

const ROUTES = ["本選考", "インターン", "OB訪問", "その他"] as const;
const SOURCES = [
  "リクナビ",
  "マイナビ",
  "ONE CAREER",
  "OfferBox",
  "企業HP",
  "i-web",
  "ワンキャリ",
  "サポーターズ",
  "その他",
] as const;

export default function NewEntryPage() {
  const router = useRouter();
  const state = useUser();

  const [companyName, setCompanyName] = useState("");
  const [route, setRoute] = useState<(typeof ROUTES)[number]>("本選考");
  const [source, setSource] = useState<(typeof SOURCES)[number]>("リクナビ");
  const [memo, setMemo] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!companyName.trim() || submitting) return;
    setError(null);
    setSubmitting(true);
    try {
      const company = await createCompany({ name: companyName.trim() });
      await createEntry({
        companyId: company.id,
        route,
        source,
        memo: memo.trim() || undefined,
      });
      router.push("/entry");
    } catch (err) {
      setSubmitting(false);
      setError(err instanceof Error ? err.message : "保存に失敗しました");
    }
  }

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[640px] px-5 py-6 md:px-8 md:py-7">
        <Link
          href="/entry"
          className="mb-3 inline-flex items-center gap-1 text-[11px] font-semibold text-ink-3 hover:text-sage"
        >
          ‹ Entry 一覧
        </Link>

        <header className="mb-5 flex items-center gap-3">
          <Mascot mood={submitting ? "cheering" : "thinking"} size={56} />
          <div>
            <p
              className="font-hand text-[20px] text-sage"
              style={{ transform: "rotate(-1.5deg)", display: "inline-block" }}
            >
              new entry,
            </p>
            <h1 className="font-serif text-2xl font-extrabold tracking-tight">
              新しいエントリー
            </h1>
          </div>
        </header>

        <form
          onSubmit={handleSubmit}
          className="rounded-xl border border-line bg-surface p-5 animate-[entre-fly-in_0.6s_cubic-bezier(0.2,0.8,0.4,1)_both]"
        >
          {/* 会社名 */}
          <Field label="会社名" required>
            <input
              type="text"
              value={companyName}
              onChange={(e) => setCompanyName(e.target.value)}
              placeholder="株式会社○○商事"
              required
              autoFocus
              className="w-full rounded-lg border border-line bg-cream px-3 py-2 text-sm font-semibold outline-none transition-colors focus:border-sage"
            />
          </Field>

          {/* 応募経路 */}
          <Field label="応募経路">
            <div className="flex gap-1.5 flex-wrap">
              {ROUTES.map((r) => (
                <button
                  key={r}
                  type="button"
                  onClick={() => setRoute(r)}
                  className={`rounded-full border px-3 py-1.5 text-[11px] font-bold transition-colors ${
                    route === r
                      ? "border-sage bg-sage text-white"
                      : "border-line bg-surface text-ink-2"
                  }`}
                  aria-pressed={route === r}
                >
                  {r}
                </button>
              ))}
            </div>
          </Field>

          {/* ソース */}
          <Field label="ソース (応募媒体)">
            <select
              value={source}
              onChange={(e) => setSource(e.target.value as (typeof SOURCES)[number])}
              className="w-full rounded-lg border border-line bg-cream px-3 py-2 text-sm font-semibold outline-none transition-colors focus:border-sage"
            >
              {SOURCES.map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </select>
          </Field>

          {/* メモ */}
          <Field label="メモ (任意)">
            <textarea
              value={memo}
              onChange={(e) => setMemo(e.target.value)}
              placeholder="気になるポイントなど"
              rows={3}
              className="w-full resize-none rounded-lg border border-line bg-cream px-3 py-2 text-sm outline-none transition-colors focus:border-sage"
            />
          </Field>

          {error && (
            <p
              role="alert"
              className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[11px] font-semibold text-ink"
            >
              {error}
            </p>
          )}

          <div className="mt-4 flex gap-2">
            <Link
              href="/entry"
              className="flex-1 rounded-lg border border-line bg-surface py-2.5 text-center text-sm font-bold text-ink-2 transition-colors hover:bg-line-2"
            >
              キャンセル
            </Link>
            <button
              type="submit"
              disabled={!companyName.trim() || submitting}
              className="flex-[2] rounded-lg bg-sage py-2.5 text-sm font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
            >
              {submitting ? "保存中…" : "＋ Entré に保存"}
            </button>
          </div>
        </form>
      </div>
    </AppShell>
  );
}

function Field({
  label,
  required,
  children,
}: {
  label: string;
  required?: boolean;
  children: React.ReactNode;
}) {
  return (
    <div className="mb-4">
      <label className="mb-1.5 block text-[10px] font-bold text-ink-2">
        {label}
        {required && <span className="ml-1 text-pink-deep">*</span>}
      </label>
      {children}
    </div>
  );
}
