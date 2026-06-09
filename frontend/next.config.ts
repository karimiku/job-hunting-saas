import path from "node:path";
import type { NextConfig } from "next";

const firebaseAuthProxyHost =
  process.env.FIREBASE_AUTH_PROXY_HOST ?? "job-hunting-saas.firebaseapp.com";

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
    ];
  },
};

export default nextConfig;
