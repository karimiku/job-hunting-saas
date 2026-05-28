"use client";

// Inbox clip 1件を削除する操作ボタン (#98)。
// 削除自体は Server Action (deleteInboxClipAction) が実行し、/inbox を revalidate する。
// ここは送信トリガと pending 表示・エラー表示だけを持つ純粋な操作部分。
// データ取得は行わない（useEffect / client fetch は不使用）。

import { useActionState } from "react";
import { useFormStatus } from "react-dom";
import {
  deleteInboxClipAction,
  type DeleteClipFormState,
} from "@/app/inbox/actions";
import type { InboxClipResponse } from "@/lib/api/inboxClips";

const INITIAL: DeleteClipFormState = {};

export function InboxClipDelete({ clip }: { clip: InboxClipResponse }) {
  const [state, formAction] = useActionState(deleteInboxClipAction, INITIAL);

  return (
    <form action={formAction} className="flex flex-col items-end gap-1">
      <input type="hidden" name="clipId" value={clip.id} />
      <DeleteButton title={clip.title} />
      {state.error && (
        <p
          role="alert"
          className="rounded-md bg-pink/40 px-2 py-1 text-[10px] font-semibold text-ink"
        >
          {state.error}
        </p>
      )}
    </form>
  );
}

function DeleteButton({ title }: { title: string }) {
  const { pending } = useFormStatus();
  return (
    <button
      type="submit"
      disabled={pending}
      aria-label={`クリップ「${title}」を削除`}
      onClick={(e) => {
        if (!window.confirm("このクリップを削除しますか？")) {
          e.preventDefault();
        }
      }}
      className="rounded-lg border border-line bg-surface px-2.5 py-1.5 text-[11px] font-bold text-ink-3 transition-colors enabled:hover:border-pink-deep enabled:hover:text-pink-deep disabled:opacity-60"
    >
      {pending ? "削除中…" : "削除"}
    </button>
  );
}
