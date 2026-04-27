import { setupServer } from "msw/node";
import { handlers } from "./msw-handlers";

/** Vitest 用の MSW サーバー。`handlers` がデフォルト挙動。 */
export const server = setupServer(...handlers);
