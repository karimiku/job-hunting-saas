"use client";

import { useEffect, useState } from "react";
import { listEntries, type EntryResponse, type ListEntriesParams } from "@/lib/api/entries";

interface FetchState {
  data: EntryResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
}

export interface UseEntriesResult {
  data: EntryResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

const INITIAL: FetchState = { data: undefined, loading: true, error: undefined };

/** エントリー一覧を取得するフック。loading / error / data の3状態を返す。 */
export function useEntries(params: ListEntriesParams = {}): UseEntriesResult {
  const [state, setState] = useState<FetchState>(INITIAL);
  const [reloadKey, setReloadKey] = useState(0);

  // params をシリアライズしてキー化（オブジェクト同値性で reload しないように）
  const key = JSON.stringify(params);

  useEffect(() => {
    let cancelled = false;
    // setState を effect body で同期実行せず、Promise 解決時のコールバックでのみ呼ぶ
    listEntries(JSON.parse(key) as ListEntriesParams)
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
  }, [key, reloadKey]);

  return {
    data: state.data,
    loading: state.loading,
    error: state.error,
    refetch: () => setReloadKey((n) => n + 1),
  };
}
