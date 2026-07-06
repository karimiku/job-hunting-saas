"use client";

// useActionState で Server Action を呼ぶフォーム。
// pending state は useFormStatus、エラーは Action からの戻り値で表示。

import { useActionState, useState } from "react";
import { useFormStatus } from "react-dom";
import Link from "next/link";
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
  values: {
    companyName: "",
    route: "本選考",
    source: "リクナビ",
    memo: "",
    flowMode: "template",
    customFlowText: "",
  },
};

export function NewEntryForm() {
  const [state, formAction] = useActionState(createNewEntryAction, INITIAL);
  const v = state.values ?? INITIAL.values!;

  return (
    <>
      <header className="mb-4">
        <h1 className="font-serif text-2xl font-extrabold tracking-tight">
          応募先を追加
        </h1>
        <p className="mt-1 text-[12px] leading-relaxed text-ink-3">
          会社名だけ入れれば登録できます。ほかはあとからでOK。
        </p>
      </header>

      <form
        action={formAction}
        className="rounded-xl border border-line bg-surface p-5"
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

        <Field label="応募の種類">
          <RouteInput defaultValue={v.route} />
        </Field>

        <Field label="どこで見つけた？">
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

        <Field label="選考フロー">
          <p className="mb-2 text-[12px] leading-relaxed text-ink-3">
            この会社の選考ステップです。標準のままでOK。あとから変更できます
          </p>
          <fieldset className="mb-2 flex flex-wrap gap-1.5">
            {[
              ["template", "標準"],
              ["custom", "手入力"],
            ].map(([value, label]) => (
              <label
                key={value}
                className="cursor-pointer rounded-full border border-line bg-surface px-3 py-1.5 text-[12px] font-bold text-ink-2 transition-colors has-[:checked]:border-sage has-[:checked]:bg-sage has-[:checked]:text-white"
              >
                <input
                  type="radio"
                  name="flowMode"
                  value={value}
                  defaultChecked={(v.flowMode ?? "template") === value}
                  className="sr-only"
                />
                {label}
              </label>
            ))}
          </fieldset>
          <textarea
            name="customFlowText"
            defaultValue={v.customFlowText}
            placeholder="例: ES提出 → Webテスト → 一次面接 → 最終面接"
            rows={3}
            className="w-full resize-none rounded-lg border border-line bg-cream px-3 py-2 text-sm outline-none transition-colors focus:border-sage"
          />
        </Field>

        {state.error && (
          <p
            role="alert"
            className="mb-3 rounded-md bg-pink/40 px-3 py-2 text-[12px] font-semibold text-ink"
          >
            {state.error}
          </p>
        )}

        <div className="mt-4 flex gap-2">
          <Link
            href="/entry"
            prefetch={false}
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

function RouteInput({ defaultValue }: { defaultValue: string }) {
  const [route, setRoute] = useState(defaultValue || "本選考");

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap gap-1.5">
        {ROUTES.map((r) => (
          <button
            key={r}
            type="button"
            onClick={() => setRoute(r)}
            className={[
              "rounded-full border px-3 py-1.5 text-[12px] font-bold transition-colors",
              route === r
                ? "border-sage bg-sage text-white"
                : "border-line bg-surface text-ink-2 hover:bg-line-2",
            ].join(" ")}
          >
            {r}
          </button>
        ))}
      </div>
      <input
        type="text"
        name="route"
        value={route}
        onChange={(e) => setRoute(e.target.value)}
        placeholder="例: 説明会経由 / 逆求人 / 直接応募"
        className="w-full rounded-lg border border-line bg-cream px-3 py-2 text-sm font-semibold outline-none transition-colors focus:border-sage"
      />
    </div>
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
      {pending ? "保存中…" : "応募先を登録"}
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
      <label className="mb-1.5 block text-[12px] font-bold text-ink-2">
        {label}
        {required && <span className="ml-1 text-pink-deep">*</span>}
      </label>
      {children}
    </div>
  );
}
