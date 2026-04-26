"use client";

import { useEffect, useState } from "react";
import { listInboxClips, type InboxClipResponse } from "@/lib/api/inboxClips";

export interface UseInboxClipsResult {
  data: InboxClipResponse[] | undefined;
  loading: boolean;
  error: Error | undefined;
  refetch: () => void;
}

/** Inbox クリップ一覧を取得するフック。 */
export function useInboxClips(): UseInboxClipsResult {
  const [data, setData] = useState<InboxClipResponse[] | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | undefined>(undefined);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(undefined);
    listInboxClips()
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
  }, [reloadKey]);

  return { data, loading, error, refetch: () => setReloadKey((n) => n + 1) };
}
