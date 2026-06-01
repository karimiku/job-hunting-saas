import { NextResponse, type NextRequest } from "next/server";

export function GET(request: NextRequest) {
  if (process.env.PLAYWRIGHT_E2E_AUTH !== "true") {
    return NextResponse.json({ message: "Not found" }, { status: 404 });
  }

  const redirectParam = request.nextUrl.searchParams.get("redirect") ?? "/dashboard";
  const redirectTo = redirectParam.startsWith("/") ? redirectParam : "/dashboard";
  const response = NextResponse.redirect(new URL(redirectTo, request.url));

  response.cookies.set("e2e-auth", "1", {
    path: "/",
    sameSite: "lax",
  });

  return response;
}
