"use client";

import { useId } from "react";

const INK = "#1F1A16";
const BG = "#FBF6EE";
const SURFACE = "#FFFFFF";
const BLUSH = "#E8C5BC";

export type Mood =
  | "happy"
  | "cheer"
  | "wink"
  | "sleep"
  | "wow"
  | "bow"
  | "shy";

const MOOD_COLORS: Record<Mood, string> = {
  happy: "#FFC83D",
  cheer: "#FF9A6C",
  wink: "#7ECDF5",
  sleep: "#C8B8EC",
  wow: "#8FD9A8",
  bow: "#F5B6C7",
  shy: "#FFD6A8",
};

export function Mascot({
  size = 120,
  mood = "happy" as Mood,
  tilt = 0,
  accent,
  animate = true,
}: {
  size?: number;
  mood?: Mood;
  tilt?: number;
  accent?: string;
  animate?: boolean;
}) {
  const shellColor = accent || MOOD_COLORS[mood];
  const eyeY = mood === "sleep" || mood === "bow" ? 56 : 52;

  const mouth = (() => {
    switch (mood) {
      case "wow":
        return <ellipse cx="60" cy="74" rx="5" ry="7" fill={INK} />;
      case "wink":
        return (
          <path
            d="M44 72 Q60 82 76 72"
            stroke={INK}
            strokeWidth="3.5"
            strokeLinecap="round"
            fill="none"
          />
        );
      case "sleep":
        return (
          <path
            d="M48 74 Q60 74 72 74"
            stroke={INK}
            strokeWidth="3.5"
            strokeLinecap="round"
            fill="none"
          />
        );
      case "cheer":
        return (
          <path
            d="M44 70 Q60 90 76 70 Q60 82 44 70"
            stroke={INK}
            strokeWidth="3"
            fill={INK}
          />
        );
      case "bow":
        return (
          <path
            d="M48 74 Q60 82 72 74"
            stroke={INK}
            strokeWidth="3.5"
            strokeLinecap="round"
            fill="none"
          />
        );
      case "shy":
        return (
          <path
            d="M50 74 Q60 80 70 74"
            stroke={INK}
            strokeWidth="3"
            strokeLinecap="round"
            fill="none"
          />
        );
      default:
        return (
          <path
            d="M44 72 Q60 84 76 72"
            stroke={INK}
            strokeWidth="3.5"
            strokeLinecap="round"
            fill="none"
          />
        );
    }
  })();

  const eyeL =
    mood === "wink" ? (
      <path
        d="M38 52 Q45 56 50 52"
        stroke={INK}
        strokeWidth="3.5"
        fill="none"
        strokeLinecap="round"
      />
    ) : mood === "sleep" || mood === "bow" ? (
      <path
        d="M36 56 L50 56"
        stroke={INK}
        strokeWidth="3"
        strokeLinecap="round"
      />
    ) : (
      <circle cx="44" cy={eyeY} r="4.5" fill={INK} />
    );

  const eyeR =
    mood === "sleep" || mood === "bow" ? (
      <path
        d="M70 56 L84 56"
        stroke={INK}
        strokeWidth="3"
        strokeLinecap="round"
      />
    ) : (
      <circle cx="76" cy={eyeY} r="4.5" fill={INK} />
    );

  const animClass = !animate
    ? ""
    : mood === "bow"
      ? "lp-mascot-bow"
      : mood === "cheer"
        ? "lp-mascot-hop"
        : mood === "shy"
          ? "lp-mascot-wobble"
          : mood === "happy"
            ? "lp-mascot-idle"
            : "";

  return (
    <svg
      width={size}
      height={size * 1.1}
      viewBox="0 0 120 132"
      className={animClass}
      style={{
        transform: `rotate(${tilt}deg)`,
        overflow: "visible",
        transformOrigin: "50% 100%",
      }}
    >
      <rect x="36" y="110" width="10" height="14" rx="3" fill={INK} />
      <rect x="74" y="110" width="10" height="14" rx="3" fill={INK} />
      <rect
        x="8"
        y="18"
        width="104"
        height="96"
        rx="18"
        fill={shellColor}
        stroke={INK}
        strokeWidth="3"
      />
      <path
        d="M8 30 L60 68 L112 30"
        fill="none"
        stroke={INK}
        strokeWidth="3"
        strokeLinejoin="round"
      />
      <rect
        x="50"
        y="8"
        width="20"
        height="14"
        rx="3"
        fill={SURFACE}
        stroke={INK}
        strokeWidth="2.5"
      />
      <rect x="54" y="12" width="12" height="6" rx="2" fill={INK} />
      <circle cx="32" cy="66" r="5" fill={BLUSH} opacity="0.9" />
      <circle cx="88" cy="66" r="5" fill={BLUSH} opacity="0.9" />
      {eyeL}
      {eyeR}
      {mouth}
    </svg>
  );
}

export function MiniMascot({ size = 32 }: { size?: number }) {
  const SUN = "#FFC83D";
  return (
    <svg width={size} height={size} viewBox="0 0 120 120">
      <rect
        x="8"
        y="18"
        width="104"
        height="96"
        rx="18"
        fill={SUN}
        stroke={INK}
        strokeWidth="5"
      />
      <path
        d="M8 30 L60 68 L112 30"
        fill="none"
        stroke={INK}
        strokeWidth="5"
        strokeLinejoin="round"
      />
      <circle cx="44" cy="54" r="6" fill={INK} />
      <circle cx="76" cy="54" r="6" fill={INK} />
      <path
        d="M44 74 Q60 86 76 74"
        stroke={INK}
        strokeWidth="4.5"
        strokeLinecap="round"
        fill="none"
      />
    </svg>
  );
}

export function Stamp({
  text = "FREE",
  color = "#4A6CF7",
  size = 80,
}: {
  text?: string;
  color?: string;
  size?: number;
}) {
  const id = useId();
  return (
    <svg width={size} height={size} viewBox="0 0 100 100">
      <defs>
        <path
          id={id}
          d="M 50,50 m -36,0 a 36,36 0 1,1 72,0 a 36,36 0 1,1 -72,0"
        />
      </defs>
      <circle cx="50" cy="50" r="44" fill={color} stroke={INK} strokeWidth="2" />
      <circle
        cx="50"
        cy="50"
        r="30"
        fill="none"
        stroke={BG}
        strokeWidth="1"
        strokeDasharray="2 3"
      />
      <text
        fontFamily="var(--lp-font-serif)"
        fontSize="11"
        fontWeight="700"
        fill={BG}
        letterSpacing="2"
      >
        <textPath href={`#${id}`} startOffset="0">
          {text} · {text} · {text} · {text} ·{" "}
        </textPath>
      </text>
      <text
        x="50"
        y="56"
        textAnchor="middle"
        fontFamily="var(--lp-font-serif)"
        fontSize="22"
        fontWeight="700"
        fill={BG}
      >
        ✓
      </text>
    </svg>
  );
}
