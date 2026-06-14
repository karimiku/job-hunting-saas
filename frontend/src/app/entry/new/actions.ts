"use server";

// Server Action — フォーム送信から company + entry を一度に作成する。
// 成功時は redirect("/entry") で投げる (Server Action 内 redirect は throw 扱いなので catch しない)。

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import { ApiError } from "@/lib/api/client-types";
import type { EntryResponse } from "@/lib/api/entries";

export interface NewEntryFormState {
  error?: string;
  values?: {
    companyName: string;
    route: string;
    source: string;
    memo: string;
  };
}

const ROUTES = ["本選考", "インターン", "OB訪問", "その他"] as const;
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

export async function createNewEntryAction(
  _prev: NewEntryFormState,
  formData: FormData,
): Promise<NewEntryFormState> {
  const companyName = readField(formData, "companyName").trim();
  const routeRaw = readField(formData, "route", "本選考");
  const sourceRaw = readField(formData, "source", "リクナビ");
  const memo = readField(formData, "memo").trim();

  // enum で受けた値だけ受理 (form の改ざん対策)
  const route = (ROUTES as readonly string[]).includes(routeRaw) ? routeRaw : "本選考";
  const source = (SOURCES as readonly string[]).includes(sourceRaw) ? sourceRaw : "リクナビ";

  const values = { companyName, route, source, memo };
  if (!companyName) {
    return { error: "会社名は必須です", values };
  }

  try {
    await serverFetch<EntryResponse>("/api/v1/entries/with-company", {
      method: "POST",
      body: JSON.stringify({
        companyName,
        route,
        source,
        memo: memo || undefined,
      }),
    });
  } catch (err) {
    return {
      error: err instanceof ApiError ? err.message : "エントリーの保存に失敗しました",
      values,
    };
  }

  // Server Component の /entry を再評価し、redirect (内部で throw)
  revalidatePath("/entry");
  redirect("/entry");
}
