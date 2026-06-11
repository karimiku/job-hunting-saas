import path from "node:path";
import type { NextConfig } from "next";

const firebaseAuthProxyHost =
  process.env.FIREBASE_AUTH_PROXY_HOST ?? "job-hunting-saas.firebaseapp.com";

const backendOrigin =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

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
};

export default nextConfig;
