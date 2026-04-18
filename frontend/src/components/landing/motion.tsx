"use client";

import {
  CSSProperties,
  ElementType,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from "react";

export function useScrollY() {
  const [y, setY] = useState(0);
  useEffect(() => {
    let raf = 0;
    const onScroll = () => {
      if (raf) return;
      raf = requestAnimationFrame(() => {
        setY(window.scrollY);
        raf = 0;
      });
    };
    window.addEventListener("scroll", onScroll, { passive: true });
    onScroll();
    return () => window.removeEventListener("scroll", onScroll);
  }, []);
  return y;
}

type RevealDir = "up" | "left" | "right" | "pop";

export function Reveal({
  children,
  delay = 0,
  className = "",
  as: Tag = "div" as ElementType,
  style,
  dir = "up",
}: {
  children: ReactNode;
  delay?: number;
  className?: string;
  as?: ElementType;
  style?: CSSProperties;
  dir?: RevealDir;
}) {
  const ref = useRef<HTMLElement | null>(null);
  const [shown, setShown] = useState(false);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const io = new IntersectionObserver(
      ([e]) => {
        if (e.isIntersecting) {
          setShown(true);
          io.disconnect();
        }
      },
      { threshold: 0.15 }
    );
    io.observe(el);
    return () => io.disconnect();
  }, []);
  const Component = Tag as ElementType;
  return (
    <Component
      ref={ref}
      className={`lp-reveal lp-reveal-${dir} ${shown ? "in" : ""} ${className}`}
      style={{ ...(style || {}), transitionDelay: `${delay}ms` }}
    >
      {children}
    </Component>
  );
}

export function Parallax({
  children,
  speed = 0.2,
  axis = "y",
  className,
  style,
}: {
  children: ReactNode;
  speed?: number;
  axis?: "x" | "y";
  className?: string;
  style?: CSSProperties;
}) {
  const ref = useRef<HTMLDivElement | null>(null);
  const [offset, setOffset] = useState(0);
  useEffect(() => {
    let raf = 0;
    const tick = () => {
      const el = ref.current;
      if (el) {
        const r = el.getBoundingClientRect();
        const mid = r.top + r.height / 2;
        const vh = window.innerHeight || 800;
        setOffset((mid - vh / 2) * -speed);
      }
      raf = requestAnimationFrame(tick);
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [speed]);
  const t =
    axis === "x"
      ? `translate3d(${offset}px,0,0)`
      : `translate3d(0,${offset}px,0)`;
  return (
    <div
      ref={ref}
      className={className}
      style={{ ...(style || {}), transform: t, willChange: "transform" }}
    >
      {children}
    </div>
  );
}

export function CountUp({
  to = 100,
  duration = 1400,
  suffix = "",
  prefix = "",
}: {
  to?: number;
  duration?: number;
  suffix?: string;
  prefix?: string;
}) {
  const ref = useRef<HTMLSpanElement | null>(null);
  const [n, setN] = useState(0);
  const done = useRef(false);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const io = new IntersectionObserver(
      ([e]) => {
        if (e.isIntersecting && !done.current) {
          done.current = true;
          const start = performance.now();
          const step = (t: number) => {
            const p = Math.min(1, (t - start) / duration);
            const eased = 1 - Math.pow(1 - p, 3);
            setN(Math.round(to * eased));
            if (p < 1) requestAnimationFrame(step);
          };
          requestAnimationFrame(step);
        }
      },
      { threshold: 0.3 }
    );
    io.observe(el);
    return () => io.disconnect();
  }, [to, duration]);
  return (
    <span ref={ref}>
      {prefix}
      {n}
      {suffix}
    </span>
  );
}

export function StaggerText({
  text,
  className,
  style,
  delay = 0,
  perChar = 30,
}: {
  text: string;
  className?: string;
  style?: CSSProperties;
  delay?: number;
  perChar?: number;
}) {
  const ref = useRef<HTMLSpanElement | null>(null);
  const [shown, setShown] = useState(false);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const io = new IntersectionObserver(
      ([e]) => {
        if (e.isIntersecting) {
          setShown(true);
          io.disconnect();
        }
      },
      { threshold: 0.3 }
    );
    io.observe(el);
    return () => io.disconnect();
  }, []);
  return (
    <span ref={ref} className={className} style={style}>
      {[...text].map((ch, i) => (
        <span
          key={i}
          style={{
            display: "inline-block",
            opacity: shown ? 1 : 0,
            transform: shown ? "translateY(0)" : "translateY(0.6em)",
            transition: `opacity 500ms cubic-bezier(0.22,1,0.36,1) ${
              delay + i * perChar
            }ms, transform 600ms cubic-bezier(0.22,1,0.36,1) ${
              delay + i * perChar
            }ms`,
            whiteSpace: "pre",
          }}
        >
          {ch}
        </span>
      ))}
    </span>
  );
}
