// Server Component。entries を SSR で取得し RoadmapView に渡す。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { listEntriesServer } from "@/lib/api/server-resources";
import { AppShell } from "@/components/entre/AppShell";
import { RoadmapView } from "@/components/entre/RoadmapView";

export default async function RoadmapPage() {
  const [user, entries] = await Promise.all([
    getCurrentUserServer(),
    listEntriesServer().catch(() => [] as never[]),
  ]);
  if (!user) redirect("/login");

  return (
    <AppShell userName={user.name} userSubtitle="○○大学 4年">
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
        <RoadmapView entries={entries} />
      </div>
    </AppShell>
  );
}
