import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { crx } from "@crxjs/vite-plugin";
import manifest from "./manifest.json" with { type: "json" };

const LOCALHOST_HOST_PERMISSIONS = new Set([
  "https://*.localhost/*",
  "http://localhost:8080/*",
]);

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const allowLocalhostPermissions =
    mode !== "production" || process.env.VITE_EXTENSION_ALLOW_LOCALHOST === "true";
  const crxManifest = allowLocalhostPermissions
    ? manifest
    : {
        ...manifest,
        host_permissions: manifest.host_permissions.filter(
          (permission) => !LOCALHOST_HOST_PERMISSIONS.has(permission),
        ),
      };

  return {
    plugins: [react(), tailwindcss(), crx({ manifest: crxManifest })],
    build: {
      outDir: "dist",
      sourcemap: true,
    },
  };
});
