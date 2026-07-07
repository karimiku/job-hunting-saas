import path from "node:path";
import type { NextConfig } from "next";

const backendOrigin = trustedProxyOrigin({
  envName: "BACKEND_API_BASE_URL",
  legacyEnvNames: ["NEXT_PUBLIC_API_BASE_URL"],
  fallback: "http://localhost:8080",
  allowedHostsEnvName: "BACKEND_API_ALLOWED_HOSTS",
  defaultAllowedHosts: [
    "localhost",
    "127.0.0.1",
    "api.entre.kamiriku.com",
    "entre-backend-gfsd4pzoxq-an.a.run.app",
  ],
});

const securityHeaders = [
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "X-Frame-Options", value: "SAMEORIGIN" },
  { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
  { key: "Permissions-Policy", value: "camera=(), microphone=(), geolocation=()" },
  { key: "Content-Security-Policy", value: buildCSP() },
  ...(process.env.NODE_ENV === "production"
    ? [{ key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains" }]
    : []),
];

// XSS の第二防御線として script/style/connect 系のポリシーを追加する。
// nonce 方式 (strict-dynamic) が最も厳格だが、proxy.ts でのリクエスト単位の
// nonce 生成・全ページ dynamic rendering 化が必要になり影響範囲が大きいため、
// まずは Next.js 公式ドキュメント (node_modules/next/dist/docs/01-app/02-guides/
// content-security-policy.md の "Without Nonces" 節) が示す構成を採用する。
//
// script-src / style-src に 'unsafe-inline' が必要な理由:
//   - Next.js は RSC のストリーミング hydration 用に nonce なしの inline
//     `<script>self.__next_f.push(...)</script>` を各ページに埋め込む。
//     nonce/strict-dynamic を使わない限り 'unsafe-inline' なしでは
//     hydration が CSP 違反でブロックされ、アプリが動作しなくなる。
//   - style-src も同様に、開発時のエラーオーバーレイ等が inline style を使う。
// 'unsafe-eval' は開発時のみ許可 (React が dev only でスタックトレース復元に
// eval を使うため。本番の React/Next はビルド成果物に eval を含まない)。
function buildCSP(): string {
  const isDev = process.env.NODE_ENV !== "production";
  const supabaseOrigin = supabaseOriginForCSP();
  const directives = [
    "default-src 'self'",
    `script-src 'self' 'unsafe-inline'${isDev ? " 'unsafe-eval'" : ""}`,
    "style-src 'self' 'unsafe-inline'",
    "img-src 'self' data:",
    "font-src 'self'",
    `connect-src 'self'${supabaseOrigin ? ` ${supabaseOrigin}` : ""}`,
    "object-src 'none'",
    "base-uri 'self'",
    "form-action 'self'",
    "frame-ancestors 'self'",
  ];
  return directives.join("; ");
}

// ブラウザから Supabase Auth (NEXT_PUBLIC_SUPABASE_URL) へ直接 fetch するため、
// connect-src に許可する。値が壊れていても CSP 生成自体は落とさず、単に追加しない。
function supabaseOriginForCSP(): string | undefined {
  const raw = process.env.NEXT_PUBLIC_SUPABASE_URL?.trim();
  if (!raw) return undefined;
  try {
    return new URL(raw).origin;
  } catch {
    return undefined;
  }
}

const nextConfig: NextConfig = {
  devIndicators: false,
  // monorepo 直下に package.json が無いため Turbopack が workspace root を
  // 誤検出する。明示的に frontend/ をルートに固定する。
  turbopack: {
    root: path.join(__dirname),
  },
  async rewrites() {
    return [
      // ブラウザからの backend 呼び出しを同一 origin に寄せる proxy。
      // cross-origin で直接叩くと POST/PATCH ごとに CORS preflight (OPTIONS) が
      // 1往復余分に乗るため、Client Component の fetch は /backend/* を使う。
      // Supabase Auth は Authorization header を主経路にし、dev/legacy cookie も同一originで扱える。
      {
        source: "/backend/:path*",
        destination: `${backendOrigin}/:path*`,
      },
    ];
  },
  async headers() {
    return [
      {
        source: "/:path*",
        headers: securityHeaders,
      },
    ];
  },
};

export default nextConfig;

function trustedProxyOrigin({
  envName,
  legacyEnvNames = [],
  fallback,
  allowedHostsEnvName,
  defaultAllowedHosts,
}: {
  envName: string;
  legacyEnvNames?: string[];
  fallback: string;
  allowedHostsEnvName: string;
  defaultAllowedHosts: string[];
}): string {
  const raw = firstEnvValue([envName, ...legacyEnvNames]) || fallback;
  let parsed: URL;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error(`${envName} must be an absolute http(s) URL`);
  }
  if (!["http:", "https:"].includes(parsed.protocol)) {
    throw new Error(`${envName} must use http or https`);
  }
  if (parsed.username || parsed.password || parsed.pathname !== "/" || parsed.search || parsed.hash) {
    throw new Error(`${envName} must be an origin without credentials, path, query, or hash`);
  }
  assertHostAllowed(parsed.hostname, allowedHosts(allowedHostsEnvName, defaultAllowedHosts), envName);
  return parsed.origin;
}

function firstEnvValue(envNames: string[]): string {
  for (const envName of envNames) {
    const value = process.env[envName]?.trim();
    if (value) return value;
  }
  return "";
}

function allowedHosts(envName: string, defaults: string[]): string[] {
  const raw = process.env[envName];
  if (!raw) return defaults;
  return raw
    .split(",")
    .map((host) => host.trim().toLowerCase())
    .filter(Boolean);
}

function assertHostAllowed(hostname: string, allowed: string[], envName: string) {
  const normalized = hostname.toLowerCase();
  if (allowed.some((allowedHost) => hostMatches(normalized, allowedHost))) return;
  throw new Error(`${envName} host ${hostname} is not in its allowed host list`);
}

function hostMatches(hostname: string, allowedHost: string): boolean {
  if (allowedHost.startsWith("*.")) {
    const suffix = allowedHost.slice(1);
    return hostname.endsWith(suffix) && hostname.length > suffix.length;
  }
  return hostname === allowedHost;
}
