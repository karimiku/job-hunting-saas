"use server";

// Server Action — Inbox clip を Company + Entry に変換し、成功後に clip を削除する。
// 新規会社の場合は backend の atomic endpoint で部分作成を防ぎ、失敗時は clip を削除しない。
// 成功時は revalidate して作成した Entry 詳細へ redirect する（redirect は throw 扱いなので catch しない）。

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { serverFetch } from "@/lib/api/server";
import { ApiError } from "@/lib/api/client-types";
import type { CompanyResponse } from "@/lib/api/companies";
import type { EntryResponse } from "@/lib/api/entries";
import type { SelectionFlowResponse } from "@/lib/selection-flow";
import { flowPayloadFromForm } from "@/lib/selection-flow";

export interface ConvertClipFormState {
  error?: string;
  values?: {
    companyName: string;
    existingCompanyId: string;
    route: string;
    source: string;
    sourceUrl: string;
    memo: string;
    flowMode?: string;
    customFlowText?: string;
  };
}

const ROUTES = ["本選考", "インターン", "OB訪問", "その他"] as const;

function readField(form: FormData, name: string, fallback = ""): string {
  const v = form.get(name);
  return typeof v === "string" ? v : fallback;
}

export async function convertInboxClipAction(
  _prev: ConvertClipFormState,
  formData: FormData,
): Promise<ConvertClipFormState> {
  const clipId = readField(formData, "clipId").trim();
  const existingCompanyId = readField(formData, "existingCompanyId").trim();
  const companyName = readField(formData, "companyName").trim();
  const routeRaw = readField(formData, "route", "本選考");
  const source = readField(formData, "source").trim();
  const sourceUrl = readField(formData, "sourceUrl").trim();
  const memo = readField(formData, "memo").trim();
  const flowMode = readField(formData, "flowMode", "template");
  const customFlowText = readField(formData, "customFlowText").trim();

  // enum で受けた値だけ受理 (form の改ざん対策)。source は backend が自由入力を許すのでそのまま。
  const route = (ROUTES as readonly string[]).includes(routeRaw) ? routeRaw : "本選考";
  const values = {
    companyName,
    existingCompanyId,
    route,
    source,
    sourceUrl,
    memo,
    flowMode,
    customFlowText,
  };

  if (!clipId) {
    return { error: "クリップの指定が不正です", values };
  }
  if (!companyName) {
    return { error: "会社名は必須です", values };
  }
  if (!source) {
    return { error: "ソースは必須です", values };
  }

  // 1. 既存 Company を使うか、新しい Company と Entry をまとめて作成する。
  //    Inbox 変換では同じ会社を何度も作ると Entry/Kanban が読みにくくなるため、
  //    フォームから渡された既存会社 ID を優先する。
  let entry: EntryResponse;
  if (existingCompanyId) {
    let company: CompanyResponse;
    try {
      company = await serverFetch<CompanyResponse>(
        `/api/v1/companies/${existingCompanyId}`,
      );
    } catch (err) {
      return {
        error:
          err instanceof ApiError
            ? err.message
            : "既存会社の取得に失敗しました",
        values,
      };
    }

    try {
      entry = await serverFetch<EntryResponse>("/api/v1/entries", {
        method: "POST",
        body: JSON.stringify({
          companyId: company.id,
          route,
          source,
          sourceUrl: sourceUrl || undefined,
          memo: memo || undefined,
        }),
      });
    } catch (err) {
      return {
        error: err instanceof ApiError ? err.message : "エントリーの保存に失敗しました",
        values,
      };
    }
  } else {
    try {
      entry = await serverFetch<EntryResponse>("/api/v1/entries/with-company", {
        method: "POST",
        body: JSON.stringify({
          companyName,
          route,
          source,
          sourceUrl: sourceUrl || undefined,
          memo: memo || undefined,
        }),
      });
    } catch (err) {
      return {
        error: err instanceof ApiError ? err.message : "エントリーの保存に失敗しました",
        values,
      };
    }
  }

  try {
    await serverFetch<SelectionFlowResponse>(
      `/api/v1/entries/${entry.id}/selection-flow`,
      {
        method: "PUT",
        body: JSON.stringify(flowPayloadFromForm(flowMode, customFlowText)),
      },
    );
  } catch (err) {
    return {
      error: err instanceof ApiError ? err.message : "選考フローの保存に失敗しました",
      values,
    };
  }

  // 2. Entry 作成成功後に clip を削除する。
  //    ここまで来れば Entry は作成済みなので、削除失敗は致命的でない（残った clip は #98 の削除UIで後始末できる）。
  try {
    await serverFetch<void>(`/api/v1/inbox/clips/${clipId}`, {
      method: "DELETE",
    });
  } catch {
    // entry は作成済み。clip 削除失敗は握りつぶして変換成功として扱う。
  }

  // 3. 関連画面を再検証して作成した Entry 詳細へ遷移する。
  revalidatePath("/inbox");
  revalidatePath("/entry");
  revalidatePath("/dashboard");
  revalidatePath("/kanban");
  redirect(`/entry/${entry.id}`);
}

export interface DeleteClipFormState {
  error?: string;
}

// Server Action — Inbox clip を削除し、/inbox を再検証する (#98)。
// 所有権チェックは backend (DELETE /api/v1/inbox/clips/{clipId} が userID で絞り込む) に委ねる。
export async function deleteInboxClipAction(
  _prev: DeleteClipFormState,
  formData: FormData,
): Promise<DeleteClipFormState> {
  const clipId = readField(formData, "clipId").trim();
  if (!clipId) {
    return { error: "クリップの指定が不正です" };
  }

  try {
    await serverFetch<void>(`/api/v1/inbox/clips/${clipId}`, {
      method: "DELETE",
    });
  } catch (err) {
    return {
      error: err instanceof ApiError ? err.message : "クリップの削除に失敗しました",
    };
  }

  revalidatePath("/inbox");
  return {};
}
