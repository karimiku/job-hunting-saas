"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Mascot, MascotMood } from "@/components/entre/Mascot";

interface Step {
  title: string;
  sub: string;
  mood: MascotMood;
  bg: string; // tailwind bg classes for the surface
}

const STEPS: Step[] = [
  {
    title: "はじめまして！",
    sub: "封筒くんと一緒に、就活はじめましょう。",
    mood: "wink",
    bg: "bg-cream-2",
  },
  {
    title: "バラバラを、ぜんぶ1枚に。",
    sub: "マイナビもONE CAREERも、あなたの就活ぜんぶ。",
    mood: "happy",
    bg: "bg-sage-wash",
  },
  {
    title: "内定までの道のりを、一緒に。",
    sub: "がんばりは、ちゃんと数えています。",
    mood: "cheering",
    bg: "bg-cream",
  },
];

export default function OnboardingPage() {
  const router = useRouter();
  const [step, setStep] = useState(0);
  const s = STEPS[step];

  const next = () => {
    if (step < STEPS.length - 1) {
      setStep((x) => x + 1);
    } else {
      router.push("/dashboard");
    }
  };

  return (
    <div className={`relative flex min-h-screen flex-col items-center justify-between overflow-hidden px-6 py-12 transition-colors duration-700 ${s.bg}`}>
      {/* Decorative sketches */}
      <svg className="absolute left-5 top-20 opacity-40" width="40" height="40" viewBox="0 0 40 40" aria-hidden>
        <path d="M5 20 Q 12 5 20 20 T 35 20" stroke="var(--color-sage)" strokeWidth="1.5" fill="none" />
      </svg>
      <svg className="absolute right-8 top-32 opacity-50" width="30" height="30" viewBox="0 0 30 30" aria-hidden>
        <path
          d="M15 4 L17 13 L26 15 L17 17 L15 26 L13 17 L4 15 L13 13 Z"
          fill="var(--color-pink)"
        />
      </svg>

      {/* Step content */}
      <div
        key={step}
        className="flex flex-1 flex-col items-center justify-center text-center animate-[entre-fade-in_0.5s_both]"
      >
        <div style={{ animation: "entre-float 3s infinite" }} className="mb-5">
          <Mascot size={120} mood={s.mood} />
        </div>
        <p
          className="font-hand text-2xl text-sage mb-1.5"
          style={{ transform: "rotate(-2deg)", display: "inline-block" }}
        >
          step {step + 1} / {STEPS.length}
        </p>
        <h1 className="mb-3 font-serif text-2xl font-extrabold leading-tight tracking-tight">
          {s.title}
        </h1>
        <p className="max-w-[280px] text-xs leading-relaxed text-ink-2">{s.sub}</p>
      </div>

      {/* Footer controls */}
      <div className="w-full max-w-[400px]">
        <div className="mb-5 flex justify-center gap-1.5">
          {STEPS.map((_, i) => (
            <span
              key={i}
              className="h-1.5 rounded-sm transition-all duration-300"
              style={{
                width: i === step ? 22 : 6,
                background: i === step ? "var(--color-sage)" : "var(--color-line)",
              }}
              aria-current={i === step ? "step" : undefined}
            />
          ))}
        </div>
        <button
          type="button"
          onClick={next}
          className="w-full rounded-xl bg-sage py-3.5 text-sm font-bold text-white transition-transform hover:-translate-y-0.5"
        >
          {step < STEPS.length - 1 ? "つぎへ →" : "はじめる ✨"}
        </button>
      </div>
    </div>
  );
}
