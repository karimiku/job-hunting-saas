"use client";

import { ReactNode, useEffect, useRef, useState } from "react";

interface RevealProps {
  children: ReactNode;
  delay?: number;
  className?: string;
}

/** ビューポート進入時にスライドアップで現れるラッパー。 */
export function Reveal({ children, delay = 0, className = "" }: RevealProps) {
  const ref = useRef<HTMLDivElement>(null);
  const [shown, setShown] = useState(false);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    // jsdom や古いブラウザのフォールバック — 次のタスクで表示する。
    // setShown を effect 内で同期に呼ぶと react-hooks/set-state-in-effect になるため。
    if (typeof IntersectionObserver === "undefined") {
      const id = setTimeout(() => setShown(true), 0);
      return () => clearTimeout(id);
    }
    const io = new IntersectionObserver(
      ([e]) => {
        if (e.isIntersecting) {
          setShown(true);
          io.disconnect();
        }
      },
      { threshold: 0.15 },
    );
    io.observe(el);
    return () => io.disconnect();
  }, []);

  return (
    <div
      ref={ref}
      className={`transition-all duration-700 ease-out ${className}`}
      style={{
        opacity: shown ? 1 : 0,
        transform: shown ? "none" : "translateY(24px) scale(0.97)",
        transitionDelay: `${delay}ms`,
      }}
    >
      {children}
    </div>
  );
}
