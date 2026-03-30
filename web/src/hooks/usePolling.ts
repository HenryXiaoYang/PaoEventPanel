import { useState, useEffect, useRef, useCallback } from "react";

export function usePolling<T>(
  fetcher: () => Promise<T>,
  intervalMs: number
): { data: T | null; refresh: () => void } {
  const [data, setData] = useState<T | null>(null);
  const fetcherRef = useRef(fetcher);
  fetcherRef.current = fetcher;

  const refresh = useCallback(async () => {
    try {
      const result = await fetcherRef.current();
      setData(result);
    } catch {
      // silently ignore polling errors
    }
  }, []);

  useEffect(() => {
    let active = true;

    const poll = async () => {
      try {
        const result = await fetcherRef.current();
        if (active) setData(result);
      } catch {
        // ignore
      }
    };

    poll();
    const id = setInterval(poll, intervalMs);
    return () => {
      active = false;
      clearInterval(id);
    };
  }, [intervalMs]);

  return { data, refresh };
}
