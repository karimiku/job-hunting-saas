"use client";

// Inbox clip 1件を Entry に変換するための開閉フォーム。
// 表示の相対時刻などは Server Component (InboxList) 側に残し、ここは操作部分だけを Client で持つ。
// useState は開閉トグルのみ。データ取得・送信は Server Action (convertInboxClipAction) に委ねる。

import { useActionState, useState } from "react";
import { useFormStatus } from "react-dom";
import {
  convertInboxClipAction,
  type ConvertClipFormState,
} from "@/app/inbox/actions";
import type { InboxClipResponse } from "@/lib/api/inboxClips";

const ROUTES = ["本選考", "インターン", "OB訪問", "その他"] as const;

function buildMemo(clip: InboxClipResponse): string {
  return [clip.title, clip.url].filter(Boolean).join("\n");
}

export function InboxClipConvert({ clip }: { clip: InboxClipResponse }) {
  const [open, setOpen] = useState(false);

  const initial: ConvertClipFormState = {
    values: {
      companyName: clip.guess?.trim() ?? "",
      route: "本選考",
      source: clip.source ?? "",
      memo: buildMemo(clip),
    },
  };
  const [state, formAction] = useActionState(convertInboxClipAction, initial);
  const v = state.values ?? initial.values!;

  if (!open) {
    return (
      <div className="flex justify-end">
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="rounded-lg border border-sage bg-sage-wash px-3 py-1.5 text-[11px] font-bold text-sage transition-colors hover:bg-sage hover:text-white"
        >
          ＋ Entry化
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
      className="rounded-lg border border-line bg-cream/60 p-3"
    >
      <input type="hidden" name="clipId" value={clip.id} />

      <Field label="会社名" htmlFor={ids.companyName} required>
        <input
          id={ids.companyName}
          name="companyName"
          type="text"
          defaultValue={v.companyName}
          placeholder="株式会社○○商事"
          required
          autoFocus
          className="w-full rounded-md border border-line bg-surface px-2.5 py-1.5 text-[12px] font-semibold outline-none transition-colors focus:border-sage"
        />
      </Field>

      <Field label="応募経路">
        <fieldset className="flex flex-wrap gap-1.5">
          {ROUTES.map((r) => (
            <label
              key={r}
              className="cursor-pointer rounded-full border border-line bg-surface px-2.5 py-1 text-[10px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
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
        </fieldset>
      </Field>

      <Field label="ソース" htmlFor={ids.source} required>
        <input
          id={ids.source}
          name="source"
          type="text"
          defaultValue={v.source}
          placeholder="マイナビ"
          required
          className="w-full rounded-md border border-line bg-surface px-2.5 py-1.5 text-[12px] font-semibold outline-none transition-colors focus:border-sage"
        />
      </Field>

      <Field label="メモ" htmlFor={ids.memo}>
        <textarea
          id={ids.memo}
          name="memo"
          defaultValue={v.memo}
          rows={2}
          className="w-full resize-none rounded-md border border-line bg-surface px-2.5 py-1.5 text-[11px] outline-none transition-colors focus:border-sage"
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
          className="flex-1 rounded-md border border-line bg-surface py-1.5 text-center text-[11px] font-bold text-ink-2 transition-colors hover:bg-line-2"
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
      className="flex-[2] rounded-md bg-sage py-1.5 text-[11px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
    >
      {pending ? "作成中…" : "Entryを作成"}
    </button>
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
  children: React.ReactNode;
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
