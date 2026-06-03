"use client";

// Inbox clip 1件を Entry に変換するための開閉フォーム。
// 表示の相対時刻などは Server Component (InboxList) 側に残し、ここは操作部分だけを Client で持つ。
// useState は開閉トグルのみ。データ取得・送信は Server Action (convertInboxClipAction) に委ねる。

import { useActionState, useState } from "react";
import { useFormStatus } from "react-dom";
import type { ReactNode } from "react";
import { ArrowRight, Building2, CheckCircle2, FileText } from "lucide-react";
import {
  convertInboxClipAction,
  type ConvertClipFormState,
} from "@/app/inbox/actions";
import type { InboxClipResponse } from "@/lib/api/inboxClips";
import type { CompanyResponse } from "@/lib/api/companies";

const ROUTES = ["本選考", "インターン", "OB訪問", "その他"] as const;

function buildMemo(clip: InboxClipResponse): string {
  return [clip.title, clip.url].filter(Boolean).join("\n");
}

function normalizeCompanyName(value: string): string {
  return value
    .replace(/[株式会社（）()・\s　]/g, "")
    .toLowerCase();
}

function findCompanyCandidates(
  guess: string,
  companies: CompanyResponse[],
): CompanyResponse[] {
  const normalizedGuess = normalizeCompanyName(guess);
  if (!normalizedGuess) return [];

  return companies
    .map((company) => ({
      company,
      normalized: normalizeCompanyName(company.name),
    }))
    .filter(({ normalized }) =>
      normalized === normalizedGuess ||
      normalized.includes(normalizedGuess) ||
      normalizedGuess.includes(normalized),
    )
    .slice(0, 3)
    .map(({ company }) => company);
}

export function InboxClipConvert({
  clip,
  companies,
}: {
  clip: InboxClipResponse;
  companies: CompanyResponse[];
}) {
  const [open, setOpen] = useState(false);
  const hasCompanyGuess = Boolean(clip.guess?.trim());
  const candidates = findCompanyCandidates(clip.guess?.trim() ?? "", companies);
  const [companyName, setCompanyName] = useState(clip.guess?.trim() ?? "");
  const [existingCompanyId, setExistingCompanyId] = useState(
    candidates[0]?.id ?? "",
  );

  const initial: ConvertClipFormState = {
    values: {
      companyName: clip.guess?.trim() ?? "",
      existingCompanyId: candidates[0]?.id ?? "",
      route: "本選考",
      source: clip.source ?? "",
      sourceUrl: clip.url,
      memo: buildMemo(clip),
    },
  };
  const [state, formAction] = useActionState(convertInboxClipAction, initial);
  const v = state.values ?? initial.values!;
  const selectedCompany = existingCompanyId
    ? companies.find((c) => c.id === existingCompanyId)
    : undefined;

  if (!open) {
    return (
      <div className="flex flex-col items-end gap-1.5">
        {!hasCompanyGuess && (
          <p className="text-[10px] font-semibold text-amber-700">
            会社名を確認してから登録
          </p>
        )}
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="inline-flex items-center gap-1.5 rounded-lg border border-sage bg-sage-wash px-3 py-1.5 text-[11px] font-bold text-sage transition-colors hover:bg-sage hover:text-white focus:outline-none focus:ring-2 focus:ring-sage/40"
        >
          <FileText size={13} aria-hidden />
          Entryとして管理
        </button>
      </div>
    );
  }

  const ids = {
    companyName: `clip-company-${clip.id}`,
    source: `clip-source-${clip.id}`,
    memo: `clip-memo-${clip.id}`,
  };

  return (
    <form
      action={formAction}
      className="rounded-lg border border-sage/40 bg-cream/70 p-3 shadow-card"
    >
      <input type="hidden" name="clipId" value={clip.id} />
      <input type="hidden" name="existingCompanyId" value={existingCompanyId} />
      <input type="hidden" name="sourceUrl" value={clip.url} />

      <div className="mb-3 rounded-md border border-line bg-surface p-2.5">
        <div className="mb-2 flex items-center justify-between gap-2">
          <div>
            <p className="text-[10px] font-black text-sage">
              保存クリップを Entry に変換
            </p>
            <p className="mt-0.5 text-[10px] leading-relaxed text-ink-3">
              会社名と応募経路を確認すると、Entry/Kanban/Task で管理できるようになります。
            </p>
          </div>
          <ArrowRight size={16} className="shrink-0 text-sage" aria-hidden />
        </div>
        <div className="grid grid-cols-[1fr_auto_1fr] items-center gap-2 text-[10px]">
          <PreviewPill
            icon={<Building2 size={12} aria-hidden />}
            label="Company"
            value={selectedCompany?.name || companyName || v.companyName || "会社名を入力"}
          />
          <span className="text-ink-3" aria-hidden>→</span>
          <PreviewPill
            icon={<FileText size={12} aria-hidden />}
            label="Entry"
            value={`${v.route}・${v.source || "ソース未入力"}`}
          />
        </div>
      </div>

      <Field label="会社名" htmlFor={ids.companyName} required>
        <input
          id={ids.companyName}
          name="companyName"
          type="text"
          value={companyName}
          onChange={(e) => {
            setCompanyName(e.target.value);
            setExistingCompanyId("");
          }}
          placeholder="株式会社○○商事"
          required
          autoFocus
          aria-describedby={`${ids.companyName}-hint`}
          className="w-full rounded-md border border-line bg-surface px-2.5 py-1.5 text-[12px] font-semibold outline-none transition-colors focus:border-sage focus:ring-2 focus:ring-sage/20"
        />
        <p id={`${ids.companyName}-hint`} className="mt-1 text-[10px] leading-relaxed text-ink-3">
          {hasCompanyGuess
            ? "求人ページから推定した候補です。正式名称に直してから作成できます。"
            : "自動検出できませんでした。求人ページの会社名を入力してください。"}
        </p>
      </Field>

      {candidates.length > 0 && (
        <fieldset className="mb-2.5 rounded-md border border-line bg-surface p-2.5">
          <legend className="px-1 text-[10px] font-bold text-ink-2">
            既存会社の候補
          </legend>
          <p className="mb-2 text-[10px] leading-relaxed text-ink-3">
            すでに登録済みの会社に紐づけると、Company の重複を防げます。
          </p>
          <div className="flex flex-col gap-1.5">
            {candidates.map((company) => (
              <label
                key={company.id}
                className="flex cursor-pointer items-center gap-2 rounded-md border border-line bg-cream px-2 py-1.5 text-[10px] transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage-wash has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-sage/30"
              >
                <input
                  type="radio"
                  name={`company-choice-${clip.id}`}
                  checked={existingCompanyId === company.id}
                  onChange={() => {
                    setExistingCompanyId(company.id);
                    setCompanyName(company.name);
                  }}
                  className="h-3 w-3 accent-sage"
                />
                <span className="min-w-0 flex-1 truncate font-bold text-ink-2">
                  {company.name}
                </span>
                <span className="shrink-0 text-[9px] font-semibold text-sage">
                  既存を使う
                </span>
              </label>
            ))}
            <label className="flex cursor-pointer items-center gap-2 rounded-md border border-line bg-cream px-2 py-1.5 text-[10px] transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage-wash has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-sage/30">
              <input
                type="radio"
                name={`company-choice-${clip.id}`}
                checked={!existingCompanyId}
                onChange={() => setExistingCompanyId("")}
                className="h-3 w-3 accent-sage"
              />
              <span className="font-bold text-ink-2">新しい会社として作成</span>
            </label>
          </div>
        </fieldset>
      )}

      <fieldset className="mb-2.5">
        <legend className="mb-1 block text-[10px] font-bold text-ink-2">
          応募タイプ
        </legend>
        <div className="flex flex-wrap gap-1.5">
          {ROUTES.map((r) => (
            <label
              key={r}
              className="cursor-pointer rounded-full border border-line bg-surface px-2.5 py-1 text-[10px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-sage/40"
            >
              <input
                type="radio"
                name="route"
                value={r}
                defaultChecked={r === v.route}
                className="sr-only"
              />
              {r}
            </label>
          ))}
        </div>
      </fieldset>

      <Field label="ソース" htmlFor={ids.source} required>
        <input
          id={ids.source}
          name="source"
          type="text"
          defaultValue={v.source}
          placeholder="マイナビ"
          required
          className="w-full rounded-md border border-line bg-surface px-2.5 py-1.5 text-[12px] font-semibold outline-none transition-colors focus:border-sage focus:ring-2 focus:ring-sage/20"
        />
      </Field>

      <Field label="メモに残す内容" htmlFor={ids.memo}>
        <textarea
          id={ids.memo}
          name="memo"
          defaultValue={v.memo}
          rows={3}
          className="w-full resize-none rounded-md border border-line bg-surface px-2.5 py-1.5 text-[11px] outline-none transition-colors focus:border-sage focus:ring-2 focus:ring-sage/20"
        />
      </Field>

      {state.error && (
        <p
          role="alert"
          className="mb-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[10px] font-semibold text-ink"
        >
          {state.error}
        </p>
      )}

      <div className="flex gap-2">
        <button
          type="button"
          onClick={() => setOpen(false)}
          className="flex-1 rounded-md border border-line bg-surface py-1.5 text-center text-[11px] font-bold text-ink-2 transition-colors hover:bg-line-2 focus:outline-none focus:ring-2 focus:ring-sage/30"
        >
          キャンセル
        </button>
        <SubmitButton />
      </div>
    </form>
  );
}

function SubmitButton() {
  const { pending } = useFormStatus();
  return (
    <button
      type="submit"
      disabled={pending}
      className="inline-flex flex-[2] items-center justify-center gap-1.5 rounded-md bg-sage py-1.5 text-[11px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40 disabled:opacity-60"
    >
      {pending ? (
        "作成中…"
      ) : (
        <>
          <CheckCircle2 size={13} aria-hidden />
          Entryを作成して開く
        </>
      )}
    </button>
  );
}

function PreviewPill({
  icon,
  label,
  value,
}: {
  icon: ReactNode;
  label: string;
  value: string;
}) {
  return (
    <div className="min-w-0 rounded-md bg-sage-wash px-2 py-1.5">
      <div className="mb-0.5 flex items-center gap-1 font-mono text-[8px] font-bold text-sage">
        {icon}
        {label}
      </div>
      <div className="truncate text-[10px] font-bold text-ink-2">{value}</div>
    </div>
  );
}

function Field({
  label,
  htmlFor,
  required,
  children,
}: {
  label: string;
  htmlFor?: string;
  required?: boolean;
  children: ReactNode;
}) {
  return (
    <div className="mb-2.5">
      <label
        htmlFor={htmlFor}
        className="mb-1 block text-[10px] font-bold text-ink-2"
      >
        {label}
        {required && <span className="ml-1 text-pink-deep">*</span>}
      </label>
      {children}
    </div>
  );
}
