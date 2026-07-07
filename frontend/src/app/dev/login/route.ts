import { NextResponse, type NextRequest } from "next/server";

const defaultBackendOrigin = "http://localhost:8080";

export async function GET(request: NextRequest) {
  if (process.env.NODE_ENV === "production") {
    return NextResponse.json({ message: "Not found" }, { status: 404 });
  }

  const backendOrigin = backendApiOrigin();
  const res = await fetch(`${backendOrigin}/dev/session`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Origin: request.nextUrl.origin,
    },
    body: JSON.stringify({
      email: request.nextUrl.searchParams.get("email") ?? undefined,
      name: request.nextUrl.searchParams.get("name") ?? undefined,
    }),
    cache: "no-store",
  });

  if (!res.ok) {
    const message = await res.text().catch(() => "");
    return NextResponse.json(
      { message: message.trim() || `dev login failed: ${res.status}` },
      { status: res.status },
    );
  }

  const redirectTo = safeRedirectPath(request.nextUrl.searchParams.get("redirect"));
  const response = new NextResponse(devLoginHtml(redirectTo), {
    status: 200,
    headers: {
      "Content-Type": "text/html; charset=utf-8",
      "Cache-Control": "no-store",
    },
  });
  for (const cookie of responseCookies(res.headers)) {
    response.headers.append("Set-Cookie", cookie);
  }
  return response;
}

function devLoginHtml(redirectTo: string): string {
  const target = JSON.stringify(redirectTo);
  return `<!doctype html>
<meta charset="utf-8">
<meta name="robots" content="noindex">
<title>Dev login</title>
<script>location.replace(${target});</script>
<body style="font-family: system-ui, sans-serif; padding: 24px;">Redirecting...</body>`;
}

function backendApiOrigin(): string {
  const raw =
    process.env.BACKEND_API_BASE_URL ??
    process.env.NEXT_PUBLIC_API_BASE_URL ??
    defaultBackendOrigin;
  const parsed = new URL(raw);
  return parsed.origin;
}

function responseCookies(headers: Headers): string[] {
  const withGetSetCookie = headers as Headers & { getSetCookie?: () => string[] };
  const cookies = withGetSetCookie.getSetCookie?.();
  if (cookies && cookies.length > 0) return cookies;
  const cookie = headers.get("set-cookie");
  return cookie ? [cookie] : [];
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
