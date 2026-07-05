import { redirect } from "next/navigation";
import { getCurrentUserServer } from "@/lib/auth-server";
import { LoginClient } from "./LoginClient";

export default async function LoginPage() {
  const user = await getCurrentUserServer();
  if (user) redirect("/dashboard");

  return <LoginClient devLoginEnabled={process.env.NODE_ENV !== "production"} />;
}
