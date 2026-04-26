"use client";

import { useEffect, useState } from "react";
import { listEntries, type EntryResponse, type ListEntriesParams } from "@/lib/api/entries";

export interface UseEntriesResult {
  data: EntryResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

/** エントリー一覧を取得するフック。loading / error / data の3状態を返す。 */
export function useEntries(params: ListEntriesParams = {}): UseEntriesResult {
  const [data, setData] = useState<EntryResponse[] | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | undefined>(undefined);
  const [reloadKey, setReloadKey] = useState(0);

  // params をシリアライズしてキー化（オブジェクト同値性で reload しないように）
  const key = JSON.stringify(params);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(undefined);
    listEntries(JSON.parse(key) as ListEntriesParams)
      .then((res) => {
        if (!cancelled) {
          setData(res);
          setLoading(false);
        }
      })
      .catch((e: unknown) => {
        if (!cancelled) {
          setError(e instanceof Error ? e : new Error(String(e)));
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [key, reloadKey]);

  return { data, loading, error, refetch: () => setReloadKey((n) => n + 1) };
}
