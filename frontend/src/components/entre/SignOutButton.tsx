"use client";

// Sign-out だけ Client (Firebase + router を使うため)。
// Dashboard / Profile から共通で使う。

import { useRouter } from "next/navigation";
import { signOut } from "@/lib/auth";

interface Props {
  className?: string;
  children?: React.ReactNode;
}

export function SignOutButton({ className, children }: Props) {
  const router = useRouter();
  return (
    <button
      type="button"
      onClick={async () => {
        await signOut();
        router.push("/login");
      }}
      className={
        className ??
        "rounded-md border border-line bg-surface px-3 py-1.5 text-[11px] font-semibold text-ink-2 transition-colors hover:bg-line-2"
      }
    >
      {children ?? "ログアウト"}
    </button>
  );
}
