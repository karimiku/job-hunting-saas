"use server";

// Server Action — フォーム送信から company + entry を一度に作成する。
// 成功時は作成した応募先の詳細へ redirect する (Server Action 内 redirect は throw 扱いなので catch しない)。

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import { ApiError } from "@/lib/api/client-types";
import type { EntryResponse } from "@/lib/api/entries";
import type { SelectionFlowResponse } from "@/lib/selection-flow";
import { flowPayloadFromForm } from "@/lib/selection-flow";

export interface NewEntryFormState {
  error?: string;
  values?: {
    companyName: string;
    route: string;
    source: string;
    memo: string;
    flowMode?: string;
    customFlowText?: string;
  };
}

const DEFAULT_ROUTE = "本選考";
const SOURCES = [
  "リクナビ",
  "マイナビ",
  "ONE CAREER",
  "OfferBox",
  "企業HP",
  "i-web",
  "ワンキャリ",
  "サポーターズ",
  "その他",
] as const;

function readField(form: FormData, name: string, fallback = ""): string {
  const v = form.get(name);
  return typeof v === "string" ? v : fallback;
}

function normalizeRoute(raw: string): string {
  return raw.trim() || DEFAULT_ROUTE;
}

export async function createNewEntryAction(
  _prev: NewEntryFormState,
  formData: FormData,
): Promise<NewEntryFormState> {
  const companyName = readField(formData, "companyName").trim();
  const routeRaw = readField(formData, "route", DEFAULT_ROUTE);
  const sourceRaw = readField(formData, "source", "リクナビ");
  const memo = readField(formData, "memo").trim();
  const flowMode = readField(formData, "flowMode", "template");
  const customFlowText = readField(formData, "customFlowText").trim();

  const route = normalizeRoute(routeRaw);
  const source = (SOURCES as readonly string[]).includes(sourceRaw) ? sourceRaw : "リクナビ";

  const values = { companyName, route, source, memo, flowMode, customFlowText };
  if (!companyName) {
    return { error: "会社名は必須です", values };
  }

  let entry: EntryResponse;
  try {
    entry = await serverFetch<EntryResponse>("/api/v1/entries/with-company", {
      method: "POST",
      body: JSON.stringify({
        companyName,
        route,
        source,
        memo: memo || undefined,
      }),
    });
    await serverFetch<SelectionFlowResponse>(
      `/api/v1/entries/${entry.id}/selection-flow`,
      {
        method: "PUT",
        body: JSON.stringify(flowPayloadFromForm(flowMode, customFlowText)),
      },
    );
  } catch (err) {
    return {
      error: err instanceof ApiError ? err.message : "応募先の登録に失敗しました",
      values,
    };
  }

  // Server Component の /entry を再評価し、作成した応募先の詳細へ redirect (内部で throw)
  revalidatePath("/entry");
  redirect(`/entry/${entry.id}`);
}
