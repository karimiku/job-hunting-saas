"use client";

import { useEffect, useState } from "react";

interface CountUpProps {
  end: number;
  duration?: number;
  suffix?: string;
  prefix?: string;
}

/** 0 から `end` までイージング付きでカウントアップする数字。 */
export function CountUp({ end, duration = 1000, suffix = "", prefix = "" }: CountUpProps) {
  const [n, setN] = useState(0);

  useEffect(() => {
    let raf: number;
    let start: number | undefined;
    const tick = (t: number) => {
      if (start === undefined) start = t;
      const p = Math.min((t - start) / duration, 1);
      const eased = 1 - Math.pow(1 - p, 3);
      setN(Math.round(end * eased));
      if (p < 1) raf = requestAnimationFrame(tick);
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [end, duration]);

  return (
    <>
      {prefix}
      {n}
      {suffix}
    </>
  );
}
