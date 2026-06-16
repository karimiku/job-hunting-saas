import path from "node:path";
import type { NextConfig } from "next";

const firebaseAuthProxyHost = trustedProxyHost({
  envName: "FIREBASE_AUTH_PROXY_HOST",
  fallback: "job-hunting-saas.firebaseapp.com",
  allowedHostsEnvName: "FIREBASE_AUTH_PROXY_ALLOWED_HOSTS",
  defaultAllowedHosts: ["*.firebaseapp.com", "*.web.app"],
});

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
  { key: "Content-Security-Policy", value: "base-uri 'self'; object-src 'none'; frame-ancestors 'self'" },
  ...(process.env.NODE_ENV === "production"
    ? [{ key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains" }]
    : []),
];

const nextConfig: NextConfig = {
  devIndicators: false,
  // monorepo 直下に package.json が無いため Turbopack が workspace root を
  // 誤検出する。明示的に frontend/ をルートに固定する。
  turbopack: {
    root: path.join(__dirname),
  },
  async rewrites() {
    return [
      {
        source: "/__/auth/:path*",
        destination: `https://${firebaseAuthProxyHost}/__/auth/:path*`,
      },
      // ブラウザからの backend 呼び出しを同一 origin に寄せる proxy。
      // cross-origin で直接叩くと POST/PATCH ごとに CORS preflight (OPTIONS) が
      // 1往復余分に乗るため、Client Component の fetch は /backend/* を使う。
      // (COOKIE_DOMAIN=.entre.kamiriku.com なので Set-Cookie もこの経路で有効)
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

function trustedProxyHost({
  envName,
  fallback,
  allowedHostsEnvName,
  defaultAllowedHosts,
}: {
  envName: string;
  fallback: string;
  allowedHostsEnvName: string;
  defaultAllowedHosts: string[];
}): string {
  const raw = process.env[envName]?.trim() || fallback;
  if (!raw || raw.includes("://") || raw.includes("/") || raw.includes("?") || raw.includes("#") || raw.includes("@")) {
    throw new Error(`${envName} must be a bare hostname`);
  }
  const parsed = new URL(`https://${raw}`);
  assertHostAllowed(parsed.hostname, allowedHosts(allowedHostsEnvName, defaultAllowedHosts), envName);
  return parsed.hostname;
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
