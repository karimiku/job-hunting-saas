"use client";

import { useEffect, useState } from "react";
import { listTasksByEntry, type TaskResponse } from "@/lib/api/tasks";

export interface UseTasksByEntryResult {
  data: TaskResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

/** entryId 配下のタスク一覧を取得。 */
export function useTasksByEntry(entryId: string | undefined): UseTasksByEntryResult {
  const [data, setData] = useState<TaskResponse[] | undefined>(undefined);
  const [loading, setLoading] = useState(Boolean(entryId));
  const [error, setError] = useState<Error | undefined>(undefined);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!entryId) {
      setLoading(false);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(undefined);
    listTasksByEntry(entryId)
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
  }, [entryId, reloadKey]);

  return { data, loading, error, refetch: () => setReloadKey((n) => n + 1) };
}
