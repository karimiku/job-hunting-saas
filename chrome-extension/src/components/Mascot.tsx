import { CSSProperties } from "react";

export type MascotMood = "happy" | "cheering" | "wink" | "thinking" | "sleeping";

interface MascotProps {
  size?: number;
  mood?: MascotMood;
  className?: string;
  style?: CSSProperties;
}

/** 封筒くん — Entré のマスコット。frontend と同じ実装。 */
export function Mascot({
  size = 64,
  mood = "happy",
  className,
  style,
}: MascotProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 100 100"
      className={className}
      style={{ display: "inline-block", ...style }}
    >
      <rect x="10" y="28" width="80" height="60" rx="8" fill="#FEF8E6" stroke="#2B2A26" strokeWidth="2" />
      <path d="M10 32 L50 60 L90 32" fill="none" stroke="#2B2A26" strokeWidth="2" strokeLinejoin="round" />
      <path d="M45 22 Q50 12 58 18 Q52 24 45 22 Z" fill="#9BB58A" stroke="#2B2A26" strokeWidth="1.5" />

      {mood === "happy" && (
        <>
          <ellipse cx="40" cy="68" rx="2.4" ry="1.2" fill="#2B2A26" />
          <ellipse cx="62" cy="68" rx="2.4" ry="1.2" fill="#2B2A26" />
          <path d="M44 76 Q51 80 58 76" stroke="#2B2A26" strokeWidth="1.5" fill="none" strokeLinecap="round" />
        </>
      )}
      {mood === "cheering" && (
        <>
          <path d="M37 66 L43 70 M43 66 L37 70" stroke="#2B2A26" strokeWidth="1.6" strokeLinecap="round" />
          <path d="M59 66 L65 70 M65 66 L59 70" stroke="#2B2A26" strokeWidth="1.6" strokeLinecap="round" />
          <path d="M42 75 Q51 84 60 75" stroke="#2B2A26" strokeWidth="1.6" fill="#fff" strokeLinecap="round" />
          <path d="M5 38 Q2 28 8 22" stroke="#2B2A26" strokeWidth="2" fill="none" strokeLinecap="round" />
          <path d="M95 38 Q98 28 92 22" stroke="#2B2A26" strokeWidth="2" fill="none" strokeLinecap="round" />
        </>
      )}
      {mood === "wink" && (
        <>
          <path d="M37 68 L43 68" stroke="#2B2A26" strokeWidth="1.6" strokeLinecap="round" />
          <ellipse cx="62" cy="68" rx="2.4" ry="1.2" fill="#2B2A26" />
          <path d="M44 76 Q51 80 58 76" stroke="#2B2A26" strokeWidth="1.5" fill="none" strokeLinecap="round" />
        </>
      )}
      {mood === "thinking" && (
        <>
          <ellipse cx="40" cy="68" rx="2.4" ry="1.2" fill="#2B2A26" />
          <ellipse cx="62" cy="68" rx="2.4" ry="1.2" fill="#2B2A26" />
          <path d="M44 78 Q51 76 58 78" stroke="#2B2A26" strokeWidth="1.5" fill="none" strokeLinecap="round" />
        </>
      )}
      {mood === "sleeping" && (
        <>
          <path d="M36 68 Q40 70 44 68" stroke="#2B2A26" strokeWidth="1.5" fill="none" strokeLinecap="round" />
          <path d="M58 68 Q62 70 66 68" stroke="#2B2A26" strokeWidth="1.5" fill="none" strokeLinecap="round" />
        </>
      )}

      <circle cx="34" cy="74" r="2.5" fill="#E9B9B0" opacity="0.7" />
      <circle cx="68" cy="74" r="2.5" fill="#E9B9B0" opacity="0.7" />
    </svg>
  );
}
