"use client";

// タスク追加フォーム共通の期日クイック選択ボタン。基準日+N日を YYYY-MM-DD で渡す。

export function addDays(base: Date, days: number): string {
  const result = new Date(base);
  result.setDate(result.getDate() + days);
  const year = result.getFullYear();
  const month = `${result.getMonth() + 1}`.padStart(2, "0");
  const day = `${result.getDate()}`.padStart(2, "0");
  return `${year}-${month}-${day}`;
}

interface QuickOption {
  label: string;
  days: number;
}

const QUICK_OPTIONS: QuickOption[] = [
  { label: "今日", days: 0 },
  { label: "明日", days: 1 },
  { label: "+3日", days: 3 },
  { label: "+1週間", days: 7 },
];

interface Props {
  onSelect: (date: string) => void;
  now?: Date;
}

export function DueDateQuickPicker({ onSelect, now = new Date() }: Props) {
  return (
    <div className="flex gap-1.5">
      {QUICK_OPTIONS.map((option) => (
        <button
          key={option.label}
          type="button"
          onClick={() => onSelect(addDays(now, option.days))}
          className="rounded-md border border-line bg-cream px-2 py-1 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
        >
          {option.label}
        </button>
      ))}
    </div>
  );
}
