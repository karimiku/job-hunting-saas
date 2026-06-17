import http from "node:http";
import { spawn } from "node:child_process";

const mockApiPort = process.env.PLAYWRIGHT_MOCK_API_PORT ?? "18080";
const frontendPort = process.env.PLAYWRIGHT_PORT ?? process.env.PORT ?? "3100";
const mockApiBase = `http://127.0.0.1:${mockApiPort}`;
const frontendBase = `http://127.0.0.1:${frontendPort}`;
const defaultLoginRedirect = `${frontendBase}/dashboard`;
const user = {
  id: "e2e-user",
  email: "e2e@example.com",
  name: "E2E Student",
};
const state = {
  companies: [],
  entries: [],
  clips: [],
  tasks: [],
  selectionFlows: [],
};

const now = () => new Date().toISOString();
const nextId = (prefix) =>
  `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;

const json = (res, status, body) => {
  res.writeHead(status, {
    "content-type": "application/json",
    "cache-control": "no-store",
  });
  res.end(JSON.stringify(body));
};

const isAuthed = (req) => /\be2e-auth=1\b/.test(req.headers.cookie ?? "");

const safeLoginRedirect = (rawRedirect) => {
  if (!rawRedirect) return defaultLoginRedirect;
  try {
    const target = new URL(rawRedirect, frontendBase);
    if (target.origin === frontendBase) {
      return target.toString();
    }
  } catch {
    // Fall through to the local default.
  }
  return defaultLoginRedirect;
};

const requireAuth = (req, res) => {
  if (isAuthed(req)) return true;
  json(res, 401, { message: "Unauthorized" });
  return false;
};

const readJson = async (req) => {
  const chunks = [];
  for await (const chunk of req) chunks.push(chunk);
  const raw = Buffer.concat(chunks).toString("utf8");
  return raw ? JSON.parse(raw) : {};
};

const entryWithDefaults = (input) => ({
  id: nextId("entry"),
  companyId: input.companyId,
  route: input.route,
  source: input.source,
  sourceUrl: input.sourceUrl ?? "",
  status: "active",
  stageKind: "application",
  stageLabel: "エントリー",
  memo: input.memo ?? "",
  createdAt: now(),
  updatedAt: now(),
});

const entryResponse = (entry) => ({
  ...entry,
  companyName: state.companies.find((company) => company.id === entry.companyId)?.name,
});

const taskWithDefaults = (entryId, input) => ({
  id: nextId("task"),
  entryId,
  title: input.title,
  type: input.type ?? "deadline",
  status: "todo",
  dueDate: input.dueDate ?? null,
  memo: input.memo ?? "",
  createdAt: now(),
  updatedAt: now(),
});

const selectionFlowWithDefaults = (entryId, input) => ({
  id: nextId("flow"),
  entryId,
  source: input.source ?? "template",
  currentStagePosition: input.currentStagePosition ?? 1,
  confidence: input.confidence ?? null,
  inboxClipId: input.inboxClipId ?? null,
  stages: (input.stages ?? []).map((stage, index) => ({
    id: nextId("stage"),
    position: index + 1,
    stageKind: stage.stageKind,
    stageLabel: stage.stageLabel,
    evidenceText: stage.evidenceText ?? "",
  })),
  createdAt: now(),
  updatedAt: now(),
});

const mockApi = http.createServer((req, res) => {
  const url = new URL(req.url ?? "/", mockApiBase);

  if (req.method === "OPTIONS") {
    res.writeHead(204);
    res.end();
    return;
  }

  if (url.pathname === "/health") {
    res.writeHead(200, { "content-type": "text/plain" });
    res.end("ok");
    return;
  }

  if (url.pathname === "/e2e/login") {
    const redirectTo = safeLoginRedirect(url.searchParams.get("redirect"));
    res.writeHead(302, {
      location: redirectTo,
      "set-cookie": "e2e-auth=1; Path=/; SameSite=Lax",
    });
    res.end();
    return;
  }

  if (url.pathname === "/auth/me") {
    if (!requireAuth(req, res)) return;
    json(res, 200, user);
    return;
  }

  if (url.pathname.startsWith("/api/") && !requireAuth(req, res)) return;

  Promise.resolve()
    .then(async () => {
      if (url.pathname === "/api/v1/page-data/task" && req.method === "GET") {
        json(res, 200, {
          user,
          entries: state.entries.map(entryResponse),
          tasks: state.tasks,
        });
        return;
      }

      if (url.pathname === "/api/v1/page-data/app" && req.method === "GET") {
        json(res, 200, {
          user,
          entries: state.entries.map(entryResponse),
          tasks: state.tasks,
          clips: state.clips,
          companies: state.companies,
        });
        return;
      }

      if (url.pathname === "/api/v1/entries" && req.method === "GET") {
        json(res, 200, { entries: state.entries.map(entryResponse) });
        return;
      }

      if (url.pathname === "/api/v1/entries" && req.method === "POST") {
        const body = await readJson(req);
        const entry = entryWithDefaults(body);
        state.entries.push(entry);
        json(res, 201, entryResponse(entry));
        return;
      }

      if (url.pathname === "/api/v1/entries/with-company" && req.method === "POST") {
        const body = await readJson(req);
        const company = {
          id: nextId("company"),
          name: body.companyName,
          memo: "",
          createdAt: now(),
          updatedAt: now(),
        };
        const entry = entryWithDefaults({
          companyId: company.id,
          route: body.route,
          source: body.source,
          sourceUrl: body.sourceUrl,
          memo: body.memo,
        });
        state.companies.push(company);
        state.entries.push(entry);
        json(res, 201, entryResponse(entry));
        return;
      }

      const entryMatch = url.pathname.match(/^\/api\/v1\/entries\/([^/]+)$/);
      if (entryMatch && req.method === "GET") {
        const entry = state.entries.find((item) => item.id === entryMatch[1]);
        json(res, entry ? 200 : 404, entry ? entryResponse(entry) : { message: "not found" });
        return;
      }

      if (entryMatch && req.method === "PATCH") {
        const body = await readJson(req);
        const entry = state.entries.find((item) => item.id === entryMatch[1]);
        if (!entry) {
          json(res, 404, { message: "not found" });
          return;
        }
        Object.assign(entry, body, { updatedAt: now() });
        json(res, 200, entryResponse(entry));
        return;
      }

      const selectionFlowMatch = url.pathname.match(
        /^\/api\/v1\/entries\/([^/]+)\/selection-flow$/,
      );
      if (selectionFlowMatch && req.method === "GET") {
        const flow = state.selectionFlows.find(
          (item) => item.entryId === selectionFlowMatch[1],
        );
        json(res, flow ? 200 : 404, flow ?? { message: "not found" });
        return;
      }

      if (selectionFlowMatch && req.method === "PUT") {
        const body = await readJson(req);
        const entry = state.entries.find((item) => item.id === selectionFlowMatch[1]);
        if (!entry) {
          json(res, 404, { message: "not found" });
          return;
        }
        const flow = selectionFlowWithDefaults(entry.id, body);
        state.selectionFlows = state.selectionFlows.filter(
          (item) => item.entryId !== entry.id,
        );
        state.selectionFlows.push(flow);
        const current = flow.stages.find(
          (stage) => stage.position === flow.currentStagePosition,
        );
        if (current) {
          entry.stageKind = current.stageKind;
          entry.stageLabel = current.stageLabel;
          entry.status = current.stageKind === "offer" ? "offered" : "in_progress";
          entry.updatedAt = now();
        }
        json(res, 200, flow);
        return;
      }

      const selectionFlowCurrentMatch = url.pathname.match(
        /^\/api\/v1\/entries\/([^/]+)\/selection-flow\/current-stage$/,
      );
      if (selectionFlowCurrentMatch && req.method === "PATCH") {
        const body = await readJson(req);
        const flow = state.selectionFlows.find(
          (item) => item.entryId === selectionFlowCurrentMatch[1],
        );
        const entry = state.entries.find(
          (item) => item.id === selectionFlowCurrentMatch[1],
        );
        if (!flow || !entry) {
          json(res, 404, { message: "not found" });
          return;
        }
        const current = flow.stages.find((stage) => stage.position === body.position);
        if (!current) {
          json(res, 400, { message: "currentStagePosition must be within stages" });
          return;
        }
        flow.currentStagePosition = body.position;
        flow.updatedAt = now();
        entry.stageKind = current.stageKind;
        entry.stageLabel = current.stageLabel;
        entry.status = current.stageKind === "offer" ? "offered" : "in_progress";
        entry.updatedAt = now();
        json(res, 200, flow);
        return;
      }

      if (url.pathname === "/api/v1/companies" && req.method === "GET") {
        json(res, 200, { companies: state.companies });
        return;
      }

      if (url.pathname === "/api/v1/companies" && req.method === "POST") {
        const body = await readJson(req);
        const company = {
          id: nextId("company"),
          name: body.name,
          memo: body.memo ?? "",
          createdAt: now(),
          updatedAt: now(),
        };
        state.companies.push(company);
        json(res, 201, company);
        return;
      }

      const companyMatch = url.pathname.match(/^\/api\/v1\/companies\/([^/]+)$/);
      if (companyMatch && req.method === "GET") {
        const company = state.companies.find((item) => item.id === companyMatch[1]);
        json(res, company ? 200 : 404, company ?? { message: "not found" });
        return;
      }

      if (companyMatch && req.method === "DELETE") {
        state.companies = state.companies.filter((item) => item.id !== companyMatch[1]);
        res.writeHead(204);
        res.end();
        return;
      }

      if (url.pathname === "/api/v1/inbox/clips" && req.method === "GET") {
        json(res, 200, { clips: state.clips });
        return;
      }

      if (url.pathname === "/api/v1/inbox/clips" && req.method === "POST") {
        const body = await readJson(req);
        const clip = {
          id: nextId("clip"),
          url: body.url,
          title: body.title,
          source: body.source,
          guess: body.guess ?? "",
          contentText: body.contentText ?? "",
          capturedAt: now(),
        };
        state.clips = state.clips.filter((item) => item.url !== clip.url);
        state.clips.push(clip);
        json(res, 201, clip);
        return;
      }

      const clipMatch = url.pathname.match(/^\/api\/v1\/inbox\/clips\/([^/]+)$/);
      if (clipMatch && req.method === "DELETE") {
        state.clips = state.clips.filter((item) => item.id !== clipMatch[1]);
        res.writeHead(204);
        res.end();
        return;
      }

      const entryTasksMatch = url.pathname.match(
        /^\/api\/v1\/entries\/([^/]+)\/tasks$/,
      );
      if (entryTasksMatch && req.method === "GET") {
        json(res, 200, {
          tasks: state.tasks.filter((task) => task.entryId === entryTasksMatch[1]),
        });
        return;
      }

      if (entryTasksMatch && req.method === "POST") {
        const body = await readJson(req);
        const task = taskWithDefaults(entryTasksMatch[1], body);
        state.tasks.push(task);
        json(res, 201, task);
        return;
      }

      if (url.pathname === "/api/v1/tasks" && req.method === "GET") {
        json(res, 200, { tasks: state.tasks });
        return;
      }

      const taskMatch = url.pathname.match(/^\/api\/v1\/tasks\/([^/]+)$/);
      if (taskMatch && req.method === "GET") {
        const task = state.tasks.find((item) => item.id === taskMatch[1]);
        json(res, task ? 200 : 404, task ?? { message: "not found" });
        return;
      }

      if (taskMatch && req.method === "PATCH") {
        const body = await readJson(req);
        const task = state.tasks.find((item) => item.id === taskMatch[1]);
        if (!task) {
          json(res, 404, { message: "not found" });
          return;
        }
        Object.assign(task, body, { updatedAt: now() });
        json(res, 200, task);
        return;
      }

      if (taskMatch && req.method === "DELETE") {
        state.tasks = state.tasks.filter((item) => item.id !== taskMatch[1]);
        res.writeHead(204);
        res.end();
        return;
      }

      json(res, 404, { message: "Not found" });
    })
    .catch((error) => {
      json(res, 500, { message: error instanceof Error ? error.message : "error" });
    });
});

await new Promise((resolve, reject) => {
  mockApi.once("error", reject);
  mockApi.listen(Number(mockApiPort), "127.0.0.1", resolve);
});

const frontend = spawn("pnpm", ["dev"], {
  stdio: "inherit",
  env: {
    ...process.env,
    PORT: frontendPort,
    BACKEND_API_BASE_URL: mockApiBase,
    PLAYWRIGHT_E2E_AUTH: "true",
  },
});

const shutdown = (signal) => {
  frontend.kill(signal);
  mockApi.close(() => process.exit(signal === "SIGINT" ? 130 : 143));
};

process.on("SIGINT", () => shutdown("SIGINT"));
process.on("SIGTERM", () => shutdown("SIGTERM"));

frontend.on("exit", (code, signal) => {
  mockApi.close(() => {
    if (signal) process.kill(process.pid, signal);
    process.exit(code ?? 0);
  });
});
