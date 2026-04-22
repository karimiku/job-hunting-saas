"use client";

const INK = "#3A3229";
const BODY = "#FCFAF6";
const STROKE = "#3A3229";
const ANTLER = "#8A6A4A";
const ANTLER_SOFT = "#A88A6A";
const BLUSH = "#EDB8A8";
const LEAF = "#8DA889";

export type DeerMood = "happy" | "sleep" | "relax" | "wave" | "worried" | "sparkle";

/**
 * Hand-drawn envelope-deer mascot matching the mockup aesthetic.
 * Soft, rounded, warm — inked with slightly wavy strokes.
 */
export function DeerMascot({
  size = 140,
  mood = "happy",
  tilt = 0,
  className,
}: {
  size?: number;
  mood?: DeerMood;
  tilt?: number;
  className?: string;
}) {
  // Eyes: slightly oval, closer together, with soft highlight
  const eye = (cx: number) => {
    if (mood === "sleep") {
      return (
        <path
          key={cx}
          d={`M${cx - 6} 66 Q${cx} 70 ${cx + 6} 66`}
          stroke={INK}
          strokeWidth="2.6"
          fill="none"
          strokeLinecap="round"
        />
      );
    }
    if (mood === "wave" && cx < 70) {
      return (
        <path
          key={cx}
          d={`M${cx - 6} 62 Q${cx} 56 ${cx + 6} 62`}
          stroke={INK}
          strokeWidth="2.6"
          fill="none"
          strokeLinecap="round"
        />
      );
    }
    if (mood === "worried") {
      return (
        <g key={cx}>
          <ellipse cx={cx} cy="64" rx="3.2" ry="3.6" fill={INK} />
          <circle cx={cx - 0.8} cy="63" r="0.9" fill="#fff" />
        </g>
      );
    }
    return (
      <g key={cx}>
        <ellipse cx={cx} cy="64" rx="3.2" ry="3.8" fill={INK} />
        <circle cx={cx - 0.8} cy="63" r="0.9" fill="#fff" />
      </g>
    );
  };

  const mouth = (() => {
    if (mood === "sleep") {
      return (
        <path
          d="M64 80 Q70 84 76 80"
          stroke={INK}
          strokeWidth="2.2"
          fill="none"
          strokeLinecap="round"
        />
      );
    }
    if (mood === "worried") {
      return (
        <path
          d="M62 82 Q70 78 78 82"
          stroke={INK}
          strokeWidth="2.2"
          fill="none"
          strokeLinecap="round"
        />
      );
    }
    if (mood === "sparkle" || mood === "happy" || mood === "wave") {
      return (
        <path
          d="M62 78 Q70 86 78 78"
          stroke={INK}
          strokeWidth="2.4"
          fill="none"
          strokeLinecap="round"
        />
      );
    }
    return (
      <path
        d="M64 80 Q70 82 76 80"
        stroke={INK}
        strokeWidth="2.2"
        fill="none"
        strokeLinecap="round"
      />
    );
  })();

  return (
    <svg
      width={size}
      height={size * 1.2}
      viewBox="0 0 150 178"
      className={className}
      style={{ transform: `rotate(${tilt}deg)`, overflow: "visible" }}
    >
      {/* Shadow underneath */}
      <ellipse cx="75" cy="170" rx="44" ry="5" fill="rgba(58, 50, 41, 0.1)" />

      {/* Antlers — wavy, branching, hand-drawn feel */}
      <g
        stroke={ANTLER}
        strokeWidth="3.6"
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        {/* Left antler */}
        <path d="M48 40 C 46 28, 42 18, 38 10" />
        <path d="M44 24 C 40 22, 34 22, 28 24" />
        <path d="M42 16 C 46 10, 46 4, 44 0" />
        {/* Right antler */}
        <path d="M102 40 C 104 28, 108 18, 112 10" />
        <path d="M106 24 C 110 22, 116 22, 122 24" />
        <path d="M108 16 C 104 10, 104 4, 106 0" />
      </g>
      {/* Antler tip highlights */}
      <g stroke={ANTLER_SOFT} strokeWidth="1.8" strokeLinecap="round" fill="none" opacity="0.65">
        <path d="M40 14 C 42 10, 42 6, 41 2" />
        <path d="M110 14 C 108 10, 108 6, 109 2" />
      </g>

      {/* Soft ear puffs */}
      <ellipse cx="42" cy="58" rx="9" ry="6" fill="#F5E6D6" stroke={STROKE} strokeWidth="2.2" />
      <ellipse cx="108" cy="58" rx="9" ry="6" fill="#F5E6D6" stroke={STROKE} strokeWidth="2.2" />
      <ellipse cx="42" cy="58" rx="5" ry="3.2" fill={BLUSH} opacity="0.6" />
      <ellipse cx="108" cy="58" rx="5" ry="3.2" fill={BLUSH} opacity="0.6" />

      {/* Envelope body — slightly irregular rounded rectangle */}
      <path
        d="M18 52
           C 18 46, 22 42, 28 42
           L 122 42
           C 128 42, 132 46, 132 52
           L 132 148
           C 132 154, 128 158, 122 158
           L 28 158
           C 22 158, 18 154, 18 148
           Z"
        fill={BODY}
        stroke={STROKE}
        strokeWidth="2.8"
        strokeLinejoin="round"
      />

      {/* Envelope V fold — subtle */}
      <path
        d="M18 58 L 75 104 L 132 58"
        fill="none"
        stroke={STROKE}
        strokeWidth="1.6"
        strokeLinejoin="round"
        opacity="0.25"
      />

      {/* Cheeks */}
      <ellipse cx="40" cy="88" rx="7.5" ry="5" fill={BLUSH} opacity="0.85" />
      <ellipse cx="110" cy="88" rx="7.5" ry="5" fill={BLUSH} opacity="0.85" />

      {/* Eyes */}
      {eye(58)}
      {eye(92)}

      {/* Mouth */}
      {mouth}

      {/* Tiny leaves near one antler — mockup detail */}
      {(mood === "happy" || mood === "wave" || mood === "sparkle") && (
        <g>
          <ellipse cx="30" cy="28" rx="4" ry="2.6" fill={LEAF} transform="rotate(-30 30 28)" />
          <ellipse cx="120" cy="28" rx="4" ry="2.6" fill={LEAF} transform="rotate(30 120 28)" />
        </g>
      )}

      {/* Worried marks */}
      {mood === "worried" && (
        <g>
          <circle cx="132" cy="34" r="3" fill="#E89B8D" />
          <rect x="131" y="18" width="2.2" height="12" rx="1.1" fill="#E89B8D" />
          <circle cx="22" cy="30" r="2.4" fill="#E89B8D" />
          <rect x="21" y="16" width="2" height="10" rx="1" fill="#E89B8D" />
        </g>
      )}

      {/* Sparkles floating around */}
      {mood === "sparkle" && (
        <g fill="#6B9079">
          <path d="M140 52 l2 5 5 2 -5 2 -2 5 -2 -5 -5 -2 5 -2 z" opacity="0.75" />
          <path d="M8 70 l1.4 3.6 3.6 1.4 -3.6 1.4 -1.4 3.6 -1.4 -3.6 -3.6 -1.4 3.6 -1.4 z" opacity="0.6" />
          <circle cx="138" cy="108" r="2" opacity="0.5" />
          <circle cx="12" cy="120" r="1.5" opacity="0.5" />
        </g>
      )}

      {/* Sleep Z's */}
      {mood === "sleep" && (
        <g style={{ transformOrigin: "130px 36px" }} className="lp-zzz">
          <text
            x="126"
            y="38"
            fontSize="14"
            fontWeight="700"
            fill={INK}
            fontFamily="var(--lp-font-serif)"
          >
            z
          </text>
          <text
            x="136"
            y="24"
            fontSize="18"
            fontWeight="700"
            fill={INK}
            fontFamily="var(--lp-font-serif)"
          >
            z
          </text>
        </g>
      )}

      {/* Wave arm */}
      {mood === "wave" && (
        <g transform="translate(120 96) rotate(-20)">
          <path
            d="M0 -12 C -2 -18, 6 -22, 10 -18 L 10 4 C 10 8, 6 10, 2 8 Z"
            fill={BODY}
            stroke={STROKE}
            strokeWidth="2.4"
            strokeLinejoin="round"
          />
        </g>
      )}
    </svg>
  );
}

/** Small mascot used in nav + footer. */
export function MiniMascot({ size = 32 }: { size?: number }) {
  return (
    <svg width={size} height={size * 1.15} viewBox="0 0 150 172">
      <g stroke={ANTLER} strokeWidth="4.4" strokeLinecap="round" fill="none">
        <path d="M48 40 C 46 28, 42 18, 38 10" />
        <path d="M42 16 C 46 10, 46 4, 44 0" />
        <path d="M102 40 C 104 28, 108 18, 112 10" />
        <path d="M108 16 C 104 10, 104 4, 106 0" />
      </g>
      <path
        d="M18 52
           C 18 46, 22 42, 28 42
           L 122 42
           C 128 42, 132 46, 132 52
           L 132 148
           C 132 154, 128 158, 122 158
           L 28 158
           C 22 158, 18 154, 18 148
           Z"
        fill={BODY}
        stroke={STROKE}
        strokeWidth="3.4"
        strokeLinejoin="round"
      />
      <ellipse cx="58" cy="84" rx="5" ry="5.8" fill={INK} />
      <ellipse cx="92" cy="84" rx="5" ry="5.8" fill={INK} />
      <path d="M62 100 Q75 110 88 100" stroke={INK} strokeWidth="3.2" fill="none" strokeLinecap="round" />
      <ellipse cx="42" cy="108" rx="6" ry="4" fill={BLUSH} opacity="0.8" />
      <ellipse cx="108" cy="108" rx="6" ry="4" fill={BLUSH} opacity="0.8" />
    </svg>
  );
}

/* ── Legacy compatibility ───────────────────────────────────── */

export type Mood = DeerMood;

export function Mascot(props: {
  size?: number;
  mood?: Mood;
  tilt?: number;
  accent?: string;
  animate?: boolean;
}) {
  return <DeerMascot size={props.size} mood={props.mood} tilt={props.tilt} />;
}

export function Stamp({
  text = "FREE",
  color = "#6B9079",
  size = 80,
}: {
  text?: string;
  color?: string;
  size?: number;
}) {
  return (
    <svg width={size} height={size} viewBox="0 0 100 100">
      <circle cx="50" cy="50" r="42" fill={color} stroke={INK} strokeWidth="2" />
      <text
        x="50"
        y="46"
        textAnchor="middle"
        fontFamily="var(--lp-font-serif)"
        fontSize="11"
        fontWeight="700"
        fill="#fff"
        letterSpacing="2"
      >
        {text}
      </text>
      <text x="50" y="66" textAnchor="middle" fontSize="18" fill="#fff">
        ✓
      </text>
    </svg>
  );
}
