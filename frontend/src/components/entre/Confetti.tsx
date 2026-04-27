"use client";

import { CSSProperties, useState } from "react";

interface ConfettiProps {
  /** 値が変わるたびに紙吹雪を再発火する。0 の間は描画しない。 */
  trigger: number;
  count?: number;
}

const COLORS = ["#4F6E58", "#E9B9B0", "#D4BA82", "#A8C0DA", "#9BB58A", "#F6EFD8"];

interface Piece {
  id: number;
  cx: number;
  cy: number;
  cr: number;
  color: string;
  delay: number;
  shape: 0 | 1 | 2;
}

/** タスク完了・内定獲得時の紙吹雪エフェクト。 */
export function Confetti({ trigger, count = 18 }: ConfettiProps) {
  if (!trigger) return null;
  // trigger ごとに ConfettiBurst を再マウントすることで pieces を再生成する。
  return <ConfettiBurst key={trigger} count={count} />;
}

/** trigger 値ごとにマウントされ、useState 初期化子で乱数値を1度だけ生成する。 */
function ConfettiBurst({ count }: { count: number }) {
  const [pieces] = useState<Piece[]>(() =>
    Array.from({ length: count }, (_, i) => ({
      id: i,
      cx: (Math.random() * 2 - 1) * 80,
      cy: 40 + Math.random() * 60,
      cr: (Math.random() * 2 - 1) * 540,
      color: COLORS[i % COLORS.length],
      delay: Math.random() * 0.1,
      shape: (i % 3) as 0 | 1 | 2,
    })),
  );

  return (
    <div className="pointer-events-none absolute left-1/2 top-[40%] z-50">
      {pieces.map((p) => {
        const style: CSSProperties & Record<string, string | number> = {
          position: "absolute",
          left: 0,
          top: 0,
          width: p.shape === 0 ? 7 : p.shape === 1 ? 5 : 9,
          height: p.shape === 0 ? 7 : p.shape === 1 ? 12 : 9,
          background: p.color,
          borderRadius: p.shape === 2 ? "50%" : 2,
          animation: `entre-confetti 1.1s cubic-bezier(.2,.8,.4,1) ${p.delay}s both`,
          "--cx": `${p.cx}px`,
          "--cy": `${p.cy}px`,
          "--cr": `${p.cr}deg`,
        };
        return <div key={p.id} style={style} />;
      })}
    </div>
  );
}
