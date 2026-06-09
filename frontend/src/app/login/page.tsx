import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { LoginClient } from "./LoginClient";

const sessionCookieName = "session";

export default async function LoginPage() {
  const cookieStore = await cookies();

  if (cookieStore.has(sessionCookieName)) {
    const user = await getCurrentUserServer();
    if (user) redirect("/dashboard");
  }

  return <LoginClient />;
}
