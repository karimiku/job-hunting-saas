import assert from "node:assert/strict";
import http from "node:http";
import test from "node:test";

import { EntreClient } from "./entre-mcp.mjs";

const entry = {
  id: "11111111-1111-4111-8111-111111111111",
  companyId: "22222222-2222-4222-8222-222222222222",
  route: "本選考",
  source: "求人ページ",
  sourceUrl: "https://example.test/jobs",
  status: "in_progress",
  stageKind: "document",
  stageLabel: "書類選考",
  memo: "memo",
};

const company = {
  id: entry.companyId,
  name: "テスト株式会社",
};

const task = {
  id: "33333333-3333-4333-8333-333333333333",
  entryId: entry.id,
  title: "ES提出",
  type: "deadline",
  dueDate: "2026-06-30T00:00:00Z",
  status: "todo",
  notify: true,
  memo: "task memo",
};

test("deleteEntry previews the deletion without calling DELETE", async (t) => {
  const fixture = await startFixtureServer();
  t.after(() => fixture.close());
  const client = new EntreClient({ baseURL: fixture.baseURL, token: "test-token" });

  const result = await client.deleteEntry({ entryRef: "entry-1" });

  assert.deepEqual(result, {
    confirmationRequired: true,
    action: "delete_entry",
    entry: publicEntryDetail(),
    relatedTaskCount: 1,
  });
  assert.equal(fixture.requests.some((request) => request.method === "DELETE"), false);
});

test("deleteEntry resolves entryRef and deletes after confirmation", async (t) => {
  const fixture = await startFixtureServer();
  t.after(() => fixture.close());
  const client = new EntreClient({ baseURL: fixture.baseURL, token: "test-token" });

  const result = await client.deleteEntry({ entryRef: "entry-1", confirm: true });

  assert.deepEqual(result, {
    deleted: true,
    entry: publicEntryDetail(),
    relatedTaskCount: 1,
  });
  assert.deepEqual(
    fixture.requests
      .filter((request) => request.method === "DELETE")
      .map((request) => request.pathname),
    [`/api/v1/entries/${entry.id}`],
  );
  assert.equal(
    fixture.requests.every((request) => request.authorization === "Bearer test-token"),
    true,
  );
});

function publicEntryDetail() {
  return {
    ref: "entry-1",
    company: company.name,
    route: entry.route,
    source: entry.source,
    status: entry.status,
    stageKind: entry.stageKind,
    stageLabel: entry.stageLabel,
    memo: entry.memo,
  };
}

async function startFixtureServer() {
  const requests = [];
  const server = http.createServer((request, response) => {
    const url = new URL(request.url, "http://127.0.0.1");
    requests.push({
      method: request.method,
      pathname: decodeURIComponent(url.pathname),
      authorization: request.headers.authorization,
    });

    if (request.method === "GET" && url.pathname === "/api/v1/entries") {
      return json(response, 200, { entries: [entry] });
    }
    if (request.method === "GET" && url.pathname === "/api/v1/companies") {
      return json(response, 200, { companies: [company] });
    }
    if (request.method === "GET" && url.pathname === `/api/v1/entries/${entry.id}`) {
      return json(response, 200, entry);
    }
    if (request.method === "GET" && url.pathname === `/api/v1/companies/${company.id}`) {
      return json(response, 200, company);
    }
    if (request.method === "GET" && url.pathname === `/api/v1/entries/${entry.id}/tasks`) {
      return json(response, 200, { tasks: [task] });
    }
    if (request.method === "DELETE" && url.pathname === `/api/v1/entries/${entry.id}`) {
      response.writeHead(204);
      return response.end();
    }
    return json(response, 404, { message: `unexpected ${request.method} ${url.pathname}` });
  });

  await new Promise((resolve, reject) => {
    server.once("error", reject);
    server.listen(0, "127.0.0.1", resolve);
  });
  const { port } = server.address();
  return {
    baseURL: `http://127.0.0.1:${port}`,
    requests,
    close: () => new Promise((resolve, reject) => server.close((error) => (error ? reject(error) : resolve()))),
  };
}

function json(response, statusCode, body) {
  response.writeHead(statusCode, { "Content-Type": "application/json" });
  response.end(JSON.stringify(body));
}
