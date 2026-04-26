"use client";

import { useEffect, useState } from "react";
import { getEntry, type EntryResponse } from "@/lib/api/entries";

interface FetchState {
  data: EntryResponse | undefined;
  loading: boolean;
  error: Error | undefined;
}

export interface UseEntryResult {
  data: EntryResponse | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

/** Entry 1件を取得するフック。id 未指定なら何もしない。 */
export function useEntry(id: string | undefined): UseEntryResult {
  // 初期 loading は id があれば true。id が undefined のときは loading=false で確定。
  const [state, setState] = useState<FetchState>({
    data: undefined,
    loading: Boolean(id),
    error: undefined,
  });
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!id) {
      // 何もしない (loading は初期値で false / 既に false)
      return;
    }
    let cancelled = false;
    getEntry(id)
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
  }, [id, reloadKey]);

  return {
    data: state.data,
    // id が消えたら loading は false 扱い (effect では state を触らない)
    loading: id ? state.loading : false,
    error: state.error,
    refetch: () => setReloadKey((n) => n + 1),
  };
}
