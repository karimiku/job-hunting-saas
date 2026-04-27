"use client";

import { useEffect, useState } from "react";
import { getEntry, type EntryResponse } from "@/lib/api/entries";

interface FetchState {
  // 取得を要求している id。state と現 prop id が一致しない間は stale 扱いにする。
  requestedId: string | undefined;
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

/** Entry 1件を取得するフック。
 *  id が変わったら即時に loading=true / data=undefined を返し、stale data を出さない。
 *  state.requestedId と現 prop id を比較して derived state で判定する (同期 setState は使わない)。 */
export function useEntry(id: string | undefined): UseEntryResult {
  const [state, setState] = useState<FetchState>(() => ({
    requestedId: id,
    data: undefined,
    loading: Boolean(id),
    error: undefined,
  }));
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    if (!id) return;
    let cancelled = false;
    getEntry(id)
      .then((res) => {
        if (!cancelled) {
          setState({ requestedId: id, data: res, loading: false, error: undefined });
        }
      })
      .catch((e: unknown) => {
        if (!cancelled) {
          setState({
            requestedId: id,
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

  const refetch = () => setReloadKey((n) => n + 1);

  // id が無い → 何も読み込まない
  if (!id) {
    return { data: undefined, loading: false, error: undefined, refetch };
  }
  // id が変わった直後 (state は前回の id 用) → 即時 loading 表示
  if (state.requestedId !== id) {
    return { data: undefined, loading: true, error: undefined, refetch };
  }
  return {
    data: state.data,
    loading: state.loading,
    error: state.error,
    refetch,
  };
}
