import { NextResponse, type NextRequest } from "next/server";

export function GET(request: NextRequest) {
  if (process.env.PLAYWRIGHT_E2E_AUTH !== "true") {
    return NextResponse.json({ message: "Not found" }, { status: 404 });
  }

  const redirectTo = safeRedirectPath(request.nextUrl.searchParams.get("redirect"));
  const response = NextResponse.redirect(new URL(redirectTo, request.url));

  response.cookies.set("e2e-auth", "1", {
    path: "/",
    sameSite: "lax",
  });

  return response;
}

function safeRedirectPath(raw: string | null): string {
  if (!raw || !raw.startsWith("/") || raw.startsWith("//")) {
    return "/dashboard";
  }
  const target = new URL(raw, "http://localhost");
  if (target.origin !== "http://localhost") {
    return "/dashboard";
  }
  return `${target.pathname}${target.search}${target.hash}`;
}
