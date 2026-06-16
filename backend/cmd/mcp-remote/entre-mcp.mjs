#!/usr/bin/env node
import { existsSync, appendFileSync } from "node:fs";
import { readdir } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import util from "node:util";

const REQUEST_TIMEOUT_MS = 20_000;

function debug(message, ...args) {
  const logPath = process.env.ENTRE_MCP_DEBUG_LOG?.trim();
  if (!logPath) return;
  const line = `${new Date().toISOString()} ${util.format(message, ...args)}\n`;
  try {
    appendFileSync(logPath, line, { mode: 0o600 });
  } catch {
    // stdout is reserved for MCP frames, so logging failures are intentionally ignored.
  }
}

async function loadMcpSdk() {
  const scriptDir = path.dirname(fileURLToPath(import.meta.url));
  const repoRoot = path.resolve(scriptDir, "../../..");
  const pnpmDir = path.join(repoRoot, "frontend", "node_modules", ".pnpm");

  const entries = await readdir(pnpmDir, { withFileTypes: true });
  const sdkDirs = entries
    .filter((entry) => entry.isDirectory() && entry.name.startsWith("@modelcontextprotocol+sdk@"))
    .map((entry) => entry.name)
    .sort();
  const sdkDir = sdkDirs[sdkDirs.length - 1];

  if (!sdkDir) {
    throw new Error(`@modelcontextprotocol/sdk was not found under ${pnpmDir}`);
  }

  const packageRoot = path.join(pnpmDir, sdkDir, "node_modules");
  const mcpPath = path.join(packageRoot, "@modelcontextprotocol", "sdk", "dist", "esm", "server", "mcp.js");
  const stdioPath = path.join(packageRoot, "@modelcontextprotocol", "sdk", "dist", "esm", "server", "stdio.js");
  const zodPath = path.join(packageRoot, "zod", "v4", "index.js");

  for (const requiredPath of [mcpPath, stdioPath, zodPath]) {
    if (!existsSync(requiredPath)) {
      throw new Error(`MCP dependency file was not found: ${requiredPath}`);
    }
  }

  const [{ McpServer }, { StdioServerTransport }, z] = await Promise.all([
    import(pathToFileURL(mcpPath).href),
    import(pathToFileURL(stdioPath).href),
    import(pathToFileURL(zodPath).href),
  ]);
  return { McpServer, StdioServerTransport, z };
}

function configFromEnv() {
  const baseURL = process.env.ENTRE_API_BASE_URL?.trim().replace(/\/+$/, "");
  const token = process.env.ENTRE_API_TOKEN?.trim();
  if (!baseURL) throw new Error("ENTRE_API_BASE_URL is required");
  if (!token) throw new Error("ENTRE_API_TOKEN is required");
  try {
    new URL(baseURL);
  } catch (error) {
    throw new Error(`ENTRE_API_BASE_URL is invalid: ${error.message}`);
  }
  return { baseURL, token };
}

function toText(value, isError = false) {
  return {
    content: [
      {
        type: "text",
        text: JSON.stringify(value, null, 2),
      },
    ],
    ...(isError ? { isError: true } : {}),
  };
}

function registerTool(server, name, options, handler) {
  server.registerTool(name, options, async (args) => {
    debug("tool call: %s", name);
    try {
      return toText(await handler(args ?? {}));
    } catch (error) {
      debug("tool error: %s: %s", name, error?.stack ?? error);
      return toText({ error: error?.message ?? String(error) }, true);
    }
  });
}

function encodePathSegment(value) {
  return encodeURIComponent(String(value ?? "").trim());
}

function stringOrNull(value) {
  if (value === null || value === undefined || value === "") return null;
  return String(value);
}

function normalizeDueDate(value) {
  const raw = String(value ?? "").trim();
  if (!raw) return null;
  if (/^\d{4}-\d{2}-\d{2}$/.test(raw)) return `${raw}T00:00:00Z`;

  const parsed = new Date(raw);
  if (Number.isNaN(parsed.getTime())) {
    throw new Error(`invalid dueDate ${JSON.stringify(raw)}: use YYYY-MM-DD or RFC3339`);
  }
  return parsed.toISOString();
}

function exposeInternalFields() {
  return process.env.ENTRE_MCP_EXPOSE_INTERNAL_IDS === "1";
}

function publicEntry(entry, companyName, ref) {
  const out = {
    ref,
    company: companyName || "unknown",
    route: entry.route,
    source: entry.source,
    status: entry.status,
    stageKind: entry.stageKind,
    stageLabel: entry.stageLabel,
  };
  if (exposeInternalFields()) {
    out.id = entry.id;
    out.companyId = entry.companyId;
    out.sourceUrl = entry.sourceUrl;
    out.memo = entry.memo;
    out.createdAt = stringOrNull(entry.createdAt);
    out.updatedAt = stringOrNull(entry.updatedAt);
  }
  return out;
}

function publicEntryDetail(entry, companyName, ref) {
  return {
    ...publicEntry(entry, companyName, ref),
    memo: entry.memo,
  };
}

function publicTask(task, companyName, entryRef, taskRef = null) {
  const out = {
    ...(taskRef ? { ref: taskRef } : {}),
    entryRef,
    company: companyName || null,
    title: task.title,
    type: task.type,
    dueDate: stringOrNull(task.dueDate),
    status: task.status,
    notify: Boolean(task.notify),
    memo: task.memo,
  };
  if (exposeInternalFields()) {
    out.id = task.id;
    out.entryId = task.entryId;
    out.createdAt = stringOrNull(task.createdAt);
    out.updatedAt = stringOrNull(task.updatedAt);
  }
  return out;
}

function publicESMemo(memo, companyName, entryRef = null) {
  const out = {
    ...(entryRef ? { entryRef } : {}),
    company: companyName || null,
    category: memo.category,
    title: memo.title,
    content: memo.content,
    source: memo.source,
  };
  if (exposeInternalFields()) {
    out.id = memo.id;
    out.entryId = memo.entryId;
    out.createdAt = stringOrNull(memo.createdAt);
    out.updatedAt = stringOrNull(memo.updatedAt);
  }
  return out;
}

class EntreClient {
  constructor({ baseURL, token }) {
    this.baseURL = baseURL;
    this.token = token;
    this.entryRefToId = new Map();
    this.entryIdToRef = new Map();
    this.taskRefToId = new Map();
    this.taskIdToRef = new Map();
    this.clipRefToId = new Map();
    this.clipIdToRef = new Map();
  }

  async get(pathname) {
    return this.request("GET", pathname);
  }

  async post(pathname, body) {
    return this.request("POST", pathname, body);
  }

  async patch(pathname, body) {
    return this.request("PATCH", pathname, body);
  }

  async put(pathname, body) {
    return this.request("PUT", pathname, body);
  }

  async request(method, pathname, body) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
    try {
      const headers = {
        Accept: "application/json",
        Authorization: `Bearer ${this.token}`,
      };
      const init = { method, headers, signal: controller.signal };
      if (body !== undefined) {
        headers["Content-Type"] = "application/json";
        init.body = JSON.stringify(body);
      }

      const response = await fetch(`${this.baseURL}${pathname}`, init);
      const responseText = await response.text();
      const data = responseText ? JSON.parse(responseText) : null;

      if (!response.ok) {
        const message = data?.message || data?.error || responseText || response.statusText;
        throw new Error(`${method} ${pathname} failed: ${message}`);
      }
      return data;
    } catch (error) {
      if (error?.name === "AbortError") {
        throw new Error(`${method} ${pathname} timed out after ${REQUEST_TIMEOUT_MS}ms`);
      }
      throw error;
    } finally {
      clearTimeout(timeout);
    }
  }

  async companyMap() {
    const data = await this.get("/api/v1/companies");
    return new Map((data?.companies ?? []).map((company) => [company.id, company.name]));
  }

  assignRef(kind, id) {
    const rawId = String(id ?? "").trim();
    if (!rawId) return null;

    const refToId = this[`${kind}RefToId`];
    const idToRef = this[`${kind}IdToRef`];
    const existing = idToRef.get(rawId);
    if (existing) return existing;

    const ref = `${kind}-${idToRef.size + 1}`;
    refToId.set(ref, rawId);
    idToRef.set(rawId, ref);
    return ref;
  }

  assignEntryRef(id) {
    return this.assignRef("entry", id);
  }

  assignTaskRef(id) {
    return this.assignRef("task", id);
  }

  assignClipRef(id) {
    return this.assignRef("clip", id);
  }

  async resolveClipRef(refOrId) {
    const raw = String(refOrId ?? "").trim();
    if (!raw) return null;
    if (this.clipRefToId.has(raw)) return this.clipRefToId.get(raw);
    if (/^\d+$/.test(raw) && this.clipRefToId.has(`clip-${raw}`)) {
      return this.clipRefToId.get(`clip-${raw}`);
    }
    if (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(raw)) {
      this.assignClipRef(raw);
      return raw;
    }

    const clips = await this.listInboxClips();
    const normalizedRef = /^\d+$/.test(raw) ? `clip-${raw}` : raw;
    const matched = clips.filter((clip) => clip.ref === normalizedRef || clip.title === raw);
    if (matched.length === 1) return this.clipRefToId.get(matched[0].ref);
    if (matched.length > 1) throw new Error(`inboxClipRef ${JSON.stringify(raw)} is ambiguous; use list_inbox_clips ref`);
    throw new Error(`unknown inboxClipRef ${JSON.stringify(raw)}; run list_inbox_clips first`);
  }

  async resolveEntryRef(refOrId) {
    const raw = String(refOrId ?? "").trim();
    if (!raw) throw new Error("entryRef is required");
    if (this.entryRefToId.has(raw)) return this.entryRefToId.get(raw);
    if (/^\d+$/.test(raw) && this.entryRefToId.has(`entry-${raw}`)) {
      return this.entryRefToId.get(`entry-${raw}`);
    }
    if (/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(raw)) {
      this.assignEntryRef(raw);
      return raw;
    }

    const { entries } = await this.entriesWithCompanies();
    const normalizedRef = /^\d+$/.test(raw) ? `entry-${raw}` : raw;
    const matched = entries.filter((entry) => entry.company === raw || entry.ref === normalizedRef);
    if (matched.length === 1) return this.entryRefToId.get(matched[0].ref);
    if (matched.length > 1) throw new Error(`entryRef ${JSON.stringify(raw)} is ambiguous; use list_entries ref`);
    throw new Error(`unknown entryRef ${JSON.stringify(raw)}; run list_entries first`);
  }

  async entriesWithCompanies() {
    const [entriesData, companies] = await Promise.all([
      this.get("/api/v1/entries"),
      this.companyMap(),
    ]);
    const entries = (entriesData?.entries ?? []).map((entry) => {
      const ref = this.assignEntryRef(entry.id);
      return publicEntry(entry, companies.get(entry.companyId), ref);
    });
    return { entries, companies, rawEntries: entriesData?.entries ?? [] };
  }

  async getCompany(companyId) {
    const id = String(companyId ?? "").trim();
    if (!id) throw new Error("companyId is required");
    return this.get(`/api/v1/companies/${encodePathSegment(id)}`);
  }

  async entryCompanyMap() {
    const [entriesData, companies] = await Promise.all([
      this.get("/api/v1/entries"),
      this.companyMap(),
    ]);
    return new Map(
      (entriesData?.entries ?? []).map((entry) => [entry.id, companies.get(entry.companyId) || entry.companyId]),
    );
  }

  async listEntries() {
    const { entries } = await this.entriesWithCompanies();
    return entries;
  }

  async getEntryContext(entryRef) {
    const id = await this.resolveEntryRef(entryRef);
    const ref = this.assignEntryRef(id);

    const entry = await this.get(`/api/v1/entries/${encodePathSegment(id)}`);
    const [company, tasksData] = await Promise.all([
      this.getCompany(entry.companyId),
      this.get(`/api/v1/entries/${encodePathSegment(id)}/tasks`),
    ]);
    return {
      entry: publicEntryDetail(entry, company.name, ref),
      tasks: (tasksData?.tasks ?? []).map((task) =>
        publicTask(task, company.name, this.assignEntryRef(task.entryId), this.assignTaskRef(task.id)),
      ),
    };
  }

  async listOpenTasks() {
    const [tasksData, entryCompanies] = await Promise.all([
      this.get("/api/v1/tasks"),
      this.entryCompanyMap(),
    ]);
    return (tasksData?.tasks ?? [])
      .filter((task) => task.status === "todo")
      .map((task) =>
        publicTask(task, entryCompanies.get(task.entryId), this.assignEntryRef(task.entryId), this.assignTaskRef(task.id)),
      );
  }

  async listInboxClips() {
    const data = await this.get("/api/v1/inbox/clips");
    return (data?.clips ?? []).map((clip) => {
      const out = {
        ref: this.assignClipRef(clip.id),
        title: clip.title,
        source: clip.source,
        guess: clip.guess,
      };
      if (exposeInternalFields()) {
        out.id = clip.id;
        out.url = clip.url;
        out.capturedAt = stringOrNull(clip.capturedAt);
      }
      return out;
    });
  }

  async listESMemos(limit) {
    const n = Number(limit ?? 20);
    const path = Number.isFinite(n) && n > 0 ? `/api/v1/es-memos?limit=${encodeURIComponent(Math.floor(n))}` : "/api/v1/es-memos";
    const [data, entryCompanies] = await Promise.all([
      this.get(path),
      this.entryCompanyMap().catch(() => new Map()),
    ]);
    return (data?.memos ?? []).map((memo) =>
      publicESMemo(
        memo,
        memo.entryId ? entryCompanies.get(memo.entryId) : null,
        memo.entryId ? this.assignEntryRef(memo.entryId) : null,
      ),
    );
  }

  async appendESMemo(input) {
    const title = String(input.title ?? "").trim();
    const content = String(input.content ?? "").trim();
    if (!title || !content) throw new Error("title and content are required");

    const entryRef = stringOrNull(input.entryRef ?? input.entryId);
    const entryId = entryRef ? await this.resolveEntryRef(entryRef) : null;
    const preview = {
      confirmationRequired: !input.confirm,
      action: "append_es_memo",
      memo: {
        title,
        content,
        category: String(input.category ?? "general").trim() || "general",
        entryRef: entryId ? this.assignEntryRef(entryId) : null,
        source: String(input.source ?? "mcp").trim() || "mcp",
      },
    };
    if (!input.confirm) return preview;

    const body = {
      ...(entryId ? { entryId } : {}),
      category: preview.memo.category,
      title,
      content,
      source: preview.memo.source,
    };
    const created = await this.post("/api/v1/es-memos", body);
    const entryCompanies = await this.entryCompanyMap().catch(() => new Map());
    return {
      created: true,
      memo: publicESMemo(
        created,
        created.entryId ? entryCompanies.get(created.entryId) : null,
        created.entryId ? this.assignEntryRef(created.entryId) : null,
      ),
    };
  }

  async createTask(input) {
    const entryRef = String(input.entryRef ?? input.entryId ?? "").trim();
    const entryId = await this.resolveEntryRef(entryRef);
    const title = String(input.title ?? "").trim();
    if (!title) throw new Error("title is required");

    const type = String(input.type ?? "deadline").trim() || "deadline";
    const dueDate = normalizeDueDate(input.dueDate);
    const entryContext = await this.getEntryContext(entryRef);
    const preview = {
      confirmationRequired: !input.confirm,
      action: "create_task",
      task: {
        entryRef: entryContext.entry.ref,
        company: entryContext.entry.company,
        title,
        type,
        dueDate,
        memo: input.memo ?? "",
        notify: Boolean(input.notify),
      },
    };
    if (!input.confirm) return preview;

    const body = {
      title,
      type,
      memo: input.memo ?? "",
      ...(dueDate ? { dueDate } : {}),
    };
    let created = await this.post(`/api/v1/entries/${encodePathSegment(entryId)}/tasks`, body);
    if (input.notify && !created.notify) {
      created = await this.patch(`/api/v1/tasks/${encodePathSegment(created.id)}`, { notify: true });
    }
    return {
      created: true,
      task: publicTask(created, entryContext.entry.company, entryContext.entry.ref, this.assignTaskRef(created.id)),
    };
  }

  async upsertEntrySelectionFlow(input) {
    const entryRef = String(input.entryRef ?? input.entryId ?? "").trim();
    const entryId = await this.resolveEntryRef(entryRef);
    const body = await this.selectionFlowRequest(input, String(input.source ?? "ai_paste"));
    const preview = {
      confirmationRequired: !input.confirm,
      action: "upsert_entry_selection_flow",
      entryRef: this.assignEntryRef(entryId),
      selectionFlow: this.publicSelectionFlowPreview(body),
    };
    if (!input.confirm) return preview;

    const flow = await this.put(`/api/v1/entries/${encodePathSegment(entryId)}/selection-flow`, body);
    return {
      updated: true,
      selectionFlow: this.publicSelectionFlow(flow),
    };
  }

  async createEntryFromJobPosting(input) {
    const companyName = String(input.companyName ?? "").trim();
    if (!companyName) throw new Error("companyName is required");

    const entryBody = {
      companyName,
      route: String(input.route ?? "本選考").trim() || "本選考",
      source: String(input.source ?? "求人ページ").trim() || "求人ページ",
      ...(stringOrNull(input.sourceUrl) ? { sourceUrl: String(input.sourceUrl).trim() } : {}),
      ...(stringOrNull(input.memo) ? { memo: String(input.memo).trim() } : {}),
    };
    const flowBody = await this.selectionFlowRequest(
      {
        ...input,
        source: input.flowSource ?? input.selectionFlowSource ?? "ai_paste",
      },
      "ai_paste",
    );
    const preview = {
      confirmationRequired: !input.confirm,
      action: "create_entry_from_job_posting",
      entry: entryBody,
      selectionFlow: this.publicSelectionFlowPreview(flowBody),
    };
    if (!input.confirm) return preview;

    const entry = await this.post("/api/v1/entries/with-company", entryBody);
    const entryRef = this.assignEntryRef(entry.id);
    const flow = await this.put(`/api/v1/entries/${encodePathSegment(entry.id)}/selection-flow`, flowBody);
    return {
      created: true,
      entry: publicEntryDetail(entry, companyName, entryRef),
      selectionFlow: this.publicSelectionFlow(flow),
    };
  }

  async selectionFlowRequest(input, defaultSource) {
    const stages = this.selectionStages(input.stages);
    const inboxClipRaw = input.inboxClipRef ?? input.inboxClipId;
    const inboxClipId = inboxClipRaw ? await this.resolveClipRef(inboxClipRaw) : null;
    const currentStagePosition = Number(input.currentStagePosition ?? 0);
    const confidence =
      input.confidence === null || input.confidence === undefined || input.confidence === ""
        ? null
        : Number(input.confidence);

    if (currentStagePosition && (!Number.isInteger(currentStagePosition) || currentStagePosition < 1)) {
      throw new Error("currentStagePosition must be a positive integer");
    }
    if (confidence !== null && (!Number.isInteger(confidence) || confidence < 0 || confidence > 100)) {
      throw new Error("confidence must be an integer between 0 and 100");
    }

    return {
      source: String(input.source ?? defaultSource).trim() || defaultSource,
      ...(currentStagePosition ? { currentStagePosition } : {}),
      ...(confidence !== null ? { confidence } : {}),
      ...(inboxClipId ? { inboxClipId } : {}),
      stages,
    };
  }

  selectionStages(rawStages) {
    const stages = Array.isArray(rawStages) ? rawStages : [];
    if (stages.length < 1) throw new Error("stages must contain at least one stage");
    if (stages.length > 20) throw new Error("stages must contain at most 20 stages");
    return stages.map((stage, index) => {
      const stageKind = String(stage?.stageKind ?? "").trim();
      const stageLabel = String(stage?.stageLabel ?? "").trim();
      if (!stageKind || !stageLabel) {
        throw new Error(`stages[${index}] requires stageKind and stageLabel`);
      }
      return {
        stageKind,
        stageLabel,
        ...(stringOrNull(stage.evidenceText) ? { evidenceText: String(stage.evidenceText).trim() } : {}),
      };
    });
  }

  publicSelectionFlowPreview(flow) {
    return {
      source: flow.source,
      currentStagePosition: flow.currentStagePosition ?? 1,
      confidence: flow.confidence ?? null,
      inboxClipRef: flow.inboxClipId ? this.assignClipRef(flow.inboxClipId) : null,
      stages: flow.stages,
    };
  }

  publicSelectionFlow(flow) {
    const out = {
      entryRef: this.assignEntryRef(flow.entryId),
      source: flow.source,
      currentStagePosition: flow.currentStagePosition,
      confidence: flow.confidence ?? null,
      inboxClipRef: flow.inboxClipId ? this.assignClipRef(flow.inboxClipId) : null,
      stages: (flow.stages ?? []).map((stage) => ({
        position: stage.position,
        stageKind: stage.stageKind,
        stageLabel: stage.stageLabel,
        evidenceText: stage.evidenceText,
      })),
    };
    if (exposeInternalFields()) {
      out.id = flow.id;
      out.entryId = flow.entryId;
      out.inboxClipId = flow.inboxClipId;
      out.createdAt = stringOrNull(flow.createdAt);
      out.updatedAt = stringOrNull(flow.updatedAt);
      out.stages = (flow.stages ?? []).map((stage) => ({
        id: stage.id,
        position: stage.position,
        stageKind: stage.stageKind,
        stageLabel: stage.stageLabel,
        evidenceText: stage.evidenceText,
      }));
    }
    return out;
  }
}

function captureJobEmail(input) {
  const text = String(input.text ?? "").trim();
  if (!text) throw new Error("text is required");

  const subject = String(input.subject ?? "").trim();
  const combined = `${subject}\n${text}`;
  const companyName = String(input.companyName ?? "").trim() || null;
  const stageKind = /面接|interview/i.test(combined)
    ? "interview"
    : /説明会|seminar|session/i.test(combined)
      ? "event"
      : /締切|期限|deadline/i.test(combined)
        ? "deadline"
        : "unknown";
  const dateMatch = combined.match(/(\d{4})[/-](\d{1,2})[/-](\d{1,2})|(\d{4})年(\d{1,2})月(\d{1,2})日/);
  const dueDate = dateMatch
    ? `${dateMatch[1] ?? dateMatch[4]}-${String(dateMatch[2] ?? dateMatch[5]).padStart(2, "0")}-${String(
        dateMatch[3] ?? dateMatch[6],
      ).padStart(2, "0")}`
    : null;

  return {
    subject: subject || null,
    companyName,
    detectedStageKind: stageKind,
    dueDate,
    taskCandidates: dueDate
      ? [
          {
            title: subject || "選考メール対応",
            type: stageKind === "deadline" ? "deadline" : "schedule",
            dueDate,
            memo: text.slice(0, 500),
          },
        ]
      : [],
    note: "LLM APIは呼ばず、ローカルの簡易抽出だけを返します。",
  };
}

async function main() {
  debug("node stdio server starting");
  const { McpServer, StdioServerTransport, z } = await loadMcpSdk();
  const client = new EntreClient(configFromEnv());
  const server = new McpServer({
    name: "entre-remote-mcp",
    version: "0.1.0",
  });
  const stageSchema = z.object({
    stageKind: z.enum(["application", "document", "test", "interview", "group", "offer", "other"]),
    stageLabel: z.string(),
    evidenceText: z.string().optional(),
  });
  const selectionFlowSchema = {
    source: z.enum(["template", "manual", "ai_inbox", "ai_paste"]),
    currentStagePosition: z.number().int().positive().optional(),
    confidence: z.number().int().min(0).max(100).optional(),
    inboxClipRef: z.string().optional().describe("list_inbox_clips が返す ref。例: clip-1"),
    inboxClipId: z.string().optional().describe("内部IDを明示したい場合のみ指定"),
    stages: z.array(stageSchema).min(1).max(20),
    confirm: z.boolean().optional(),
  };

  registerTool(
    server,
    "list_entries",
    {
      description: "応募先一覧を本番APIから取得します。公開しやすいよう内部UUIDやtimestampは返しません。",
      inputSchema: {},
    },
    () => client.listEntries(),
  );

  registerTool(
    server,
    "get_entry_context",
    {
      description: "応募先1件と紐づくTaskを本番APIから取得します。entryRef は list_entries の ref を指定します。",
      inputSchema: {
        entryRef: z.string().describe("list_entries が返す ref。例: entry-1"),
      },
    },
    ({ entryRef, entryId }) => client.getEntryContext(entryRef ?? entryId),
  );

  registerTool(
    server,
    "list_open_tasks",
    {
      description: "未完了Task一覧を本番APIから取得します。内部UUIDやtimestampは返しません。",
      inputSchema: {},
    },
    () => client.listOpenTasks(),
  );

  registerTool(
    server,
    "list_inbox_clips",
    {
      description: "Inbox clip一覧を本番APIから取得します。URLや内部IDはデフォルトでは返しません。",
      inputSchema: {},
    },
    () => client.listInboxClips(),
  );

  registerTool(
    server,
    "list_es_memos",
    {
      description: "ES/自己PR/面接ネタ用メモ一覧を本番APIから取得します。内部UUIDやtimestampは返しません。",
      inputSchema: {
        limit: z.number().int().positive().optional(),
      },
    },
    ({ limit }) => client.listESMemos(limit),
  );

  registerTool(
    server,
    "append_es_memo",
    {
      description: "ES/自己PR/ガクチカ/面接ネタ用メモの保存候補を返します。confirm=true のときだけ本番APIへ保存します。",
      inputSchema: {
        title: z.string(),
        content: z.string(),
        category: z.string().optional(),
        entryRef: z.string().optional(),
        source: z.string().optional(),
        confirm: z.boolean().optional(),
      },
    },
    (input) => client.appendESMemo(input),
  );

  registerTool(
    server,
    "create_task",
    {
      description: "Entryに紐づくTaskを作成します。entryRef は list_entries の ref を指定します。confirm=true のときだけ本番APIへ保存します。",
      inputSchema: {
        entryRef: z.string(),
        title: z.string(),
        type: z.enum(["deadline", "schedule"]).optional(),
        dueDate: z.string().optional().describe("YYYY-MM-DD または RFC3339"),
        memo: z.string().optional(),
        notify: z.boolean().optional(),
        confirm: z.boolean().optional(),
      },
    },
    (input) => client.createTask(input),
  );

  registerTool(
    server,
    "capture_job_email",
    {
      description: "選考メール本文を簡易抽出し、Entry更新候補とTask作成候補を返します。LLM APIは呼びません。",
      inputSchema: {
        text: z.string(),
        subject: z.string().optional(),
        companyName: z.string().optional(),
      },
    },
    captureJobEmail,
  );

  registerTool(
    server,
    "upsert_entry_selection_flow",
    {
      description:
        "既存Entryに会社ごとの可変選考フローを保存します。entryRef は list_entries の ref を指定します。confirm=true のときだけ本番APIへ保存します。",
      inputSchema: {
        entryRef: z.string().describe("list_entries が返す ref。例: entry-1"),
        ...selectionFlowSchema,
      },
    },
    (input) => client.upsertEntrySelectionFlow(input),
  );

  registerTool(
    server,
    "create_entry_from_job_posting",
    {
      description:
        "求人ページ本文をLLMが読んで作ったEntry候補と可変選考フローを新規保存します。confirm=true のときだけ本番APIへ保存します。",
      inputSchema: {
        companyName: z.string(),
        route: z.string().optional().describe("本選考、インターンなど。省略時は本選考。"),
        source: z.string().optional().describe("応募媒体名。例: サポーターズ、外資就活、求人ページ"),
        sourceUrl: z.string().optional().describe("求人URL"),
        memo: z.string().optional(),
        flowSource: z.enum(["template", "manual", "ai_inbox", "ai_paste"]).optional().describe("選考フロー由来。省略時は ai_paste。"),
        currentStagePosition: z.number().int().positive().optional(),
        confidence: z.number().int().min(0).max(100).optional(),
        inboxClipRef: z.string().optional().describe("list_inbox_clips が返す ref。例: clip-1"),
        inboxClipId: z.string().optional().describe("内部IDを明示したい場合のみ指定"),
        stages: z.array(stageSchema).min(1).max(20),
        confirm: z.boolean().optional(),
      },
    },
    (input) => client.createEntryFromJobPosting(input),
  );

  const transport = new StdioServerTransport();
  await server.connect(transport);
  debug("node stdio server connected");
}

main().catch((error) => {
  debug("fatal error: %s", error?.stack ?? error);
  process.stderr.write(`entre MCP server error: ${error?.message ?? String(error)}\n`);
  process.exit(1);
});
