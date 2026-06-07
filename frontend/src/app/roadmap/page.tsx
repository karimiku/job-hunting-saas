// Roadmap はコア導線から外したため、認証後はホームへ戻す。

import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";

export default async function RoadmapPage() {
  const user = await getCurrentUserServer();
  if (!user) redirect("/login");

  redirect("/dashboard");
}
