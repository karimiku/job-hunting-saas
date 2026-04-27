"use client";

// useActionState で Server Action を呼ぶフォーム。
// pending state は useFormStatus、エラーは Action からの戻り値で表示。

import { useActionState } from "react";
import { useFormStatus } from "react-dom";
import Link from "next/link";
import { Mascot } from "@/components/entre/Mascot";
import { createNewEntryAction, type NewEntryFormState } from "./actions";

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

const INITIAL: NewEntryFormState = {
  values: { companyName: "", route: "本選考", source: "リクナビ", memo: "" },
};

export function NewEntryForm() {
  const [state, formAction] = useActionState(createNewEntryAction, INITIAL);
  const v = state.values ?? INITIAL.values!;

  return (
    <>
      <header className="mb-5 flex items-center gap-3">
        <Mascot mood="thinking" size={56} />
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
        action={formAction}
        className="rounded-xl border border-line bg-surface p-5 animate-[entre-fly-in_0.6s_cubic-bezier(0.2,0.8,0.4,1)_both]"
      >
        <Field label="会社名" required>
          <input
            type="text"
            name="companyName"
            defaultValue={v.companyName}
            placeholder="株式会社○○商事"
            required
            autoFocus
            className="w-full rounded-lg border border-line bg-cream px-3 py-2 text-sm font-semibold outline-none transition-colors focus:border-sage"
          />
        </Field>

        <Field label="応募経路">
          <RouteRadio defaultValue={v.route} />
        </Field>

        <Field label="ソース (応募媒体)">
          <select
            name="source"
            defaultValue={v.source}
            className="w-full rounded-lg border border-line bg-cream px-3 py-2 text-sm font-semibold outline-none transition-colors focus:border-sage"
          >
            {SOURCES.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
        </Field>

        <Field label="メモ (任意)">
          <textarea
            name="memo"
            defaultValue={v.memo}
            placeholder="気になるポイントなど"
            rows={3}
            className="w-full resize-none rounded-lg border border-line bg-cream px-3 py-2 text-sm outline-none transition-colors focus:border-sage"
          />
        </Field>

        {state.error && (
          <p
            role="alert"
            className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[11px] font-semibold text-ink"
          >
            {state.error}
          </p>
        )}

        <div className="mt-4 flex gap-2">
          <Link
            href="/entry"
            className="flex-1 rounded-lg border border-line bg-surface py-2.5 text-center text-sm font-bold text-ink-2 transition-colors hover:bg-line-2"
          >
            キャンセル
          </Link>
          <SubmitButton />
        </div>
      </form>
    </>
  );
}

/** route はラジオではなくボタンチップ。state は input[type=hidden] で送る。 */
function RouteRadio({ defaultValue }: { defaultValue: string }) {
  return (
    <fieldset className="flex flex-wrap gap-1.5">
      {ROUTES.map((r) => (
        <label
          key={r}
          className="cursor-pointer rounded-full border border-line bg-surface px-3 py-1.5 text-[11px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
        >
          <input
            type="radio"
            name="route"
            value={r}
            defaultChecked={r === defaultValue}
            className="sr-only"
          />
          {r}
        </label>
      ))}
    </fieldset>
  );
}

function SubmitButton() {
  const { pending } = useFormStatus();
  return (
    <button
      type="submit"
      disabled={pending}
      className="flex-[2] rounded-lg bg-sage py-2.5 text-sm font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
    >
      {pending ? "保存中…" : "＋ Entré に保存"}
    </button>
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
