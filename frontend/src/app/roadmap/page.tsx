"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUser } from "@/lib/use-user";
import { AppShell } from "@/components/entre/AppShell";
import { RoadmapView } from "@/components/entre/RoadmapView";

export default function RoadmapPage() {
  const router = useRouter();
  const state = useUser();

  useEffect(() => {
    if (state.status === "guest") router.replace("/login");
  }, [state.status, router]);

  if (state.status !== "authenticated") {
    return <div className="min-h-screen bg-cream" />;
  }

  return (
    <AppShell userName={state.user.name} userSubtitle="○○大学 4年">
      <div className="mx-auto max-w-[1100px] px-5 py-7 md:px-9 md:py-9">
        <header className="mb-6 md:mb-7">
          <p
            className="font-hand text-[22px] text-sage md:text-[26px]"
            style={{ transform: "rotate(-1.5deg)", display: "inline-block" }}
          >
            your journey,
          </p>
          <h1 className="font-serif text-3xl font-extrabold tracking-tight md:text-[34px]">
            就活ロードマップ
          </h1>
          <p className="mt-1 text-[12px] text-ink-2">
            内定までの道のりを、一緒に歩いていきましょう。
          </p>
        </header>
        <RoadmapView />
      </div>
    </AppShell>
  );
}
