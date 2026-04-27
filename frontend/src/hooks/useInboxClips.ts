"use client";

import { useEffect, useState } from "react";
import { listInboxClips, type InboxClipResponse } from "@/lib/api/inboxClips";

interface FetchState {
  data: InboxClipResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
}

export interface UseInboxClipsResult {
  data: InboxClipResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

const INITIAL: FetchState = { data: undefined, loading: true, error: undefined };

/** Inbox クリップ一覧を取得するフック。 */
export function useInboxClips(): UseInboxClipsResult {
  const [state, setState] = useState<FetchState>(INITIAL);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    listInboxClips()
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
  }, [reloadKey]);

  return {
    data: state.data,
    loading: state.loading,
    error: state.error,
    refetch: () => setReloadKey((n) => n + 1),
  };
}
