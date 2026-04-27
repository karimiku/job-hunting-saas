import path from "node:path";
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  devIndicators: false,
  // monorepo 直下に package.json が無いため Turbopack が workspace root を
  // 誤検出する。明示的に frontend/ をルートに固定する。
  turbopack: {
    root: path.join(__dirname),
  },
};

export default nextConfig;
