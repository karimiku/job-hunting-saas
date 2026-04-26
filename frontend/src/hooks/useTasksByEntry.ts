"use client";

import { useEffect, useState } from "react";
import { listTasksByEntry, type TaskResponse } from "@/lib/api/tasks";

interface FetchState {
  data: TaskResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
}

export interface UseTasksByEntryResult {
  data: TaskResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

/** entryId 配下のタスク一覧を取得。 */
export function useTasksByEntry(entryId: string | undefined): UseTasksByEntryResult {
  const [state, setState] = useState<FetchState>({
    data: undefined,
    loading: Boolean(entryId),
    error: undefined,
  });
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!entryId) {
      return;
    }
    let cancelled = false;
    listTasksByEntry(entryId)
      .then((res) => {
        if (!cancelled) {
          setState({ data: res, loading: false, error: undefined });
        }
      })
      .catch((e: unknown) => {
        if (!cancelled) {
          setState({
            data: undefined,
            loading: false,
            error: e instanceof Error ? e : new Error(String(e)),
          });
        }
      });
    return () => {
      cancelled = true;
    };
  }, [entryId, reloadKey]);

  return {
    data: state.data,
    loading: entryId ? state.loading : false,
    error: state.error,
    refetch: () => setReloadKey((n) => n + 1),
  };
}
