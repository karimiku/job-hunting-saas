"use client";

import { useEffect, useState } from "react";
import { getEntry, type EntryResponse } from "@/lib/api/entries";

export interface UseEntryResult {
  data: EntryResponse | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

/** Entry 1件を取得するフック。id 未指定なら何もしない。 */
export function useEntry(id: string | undefined): UseEntryResult {
  const [data, setData] = useState<EntryResponse | undefined>(undefined);
  const [loading, setLoading] = useState(Boolean(id));
  const [error, setError] = useState<Error | undefined>(undefined);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!id) {
      setLoading(false);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(undefined);
    getEntry(id)
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
  }, [id, reloadKey]);

  return { data, loading, error, refetch: () => setReloadKey((n) => n + 1) };
}
