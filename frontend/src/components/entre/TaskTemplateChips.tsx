"use client";

// タスク追加フォーム共通の定型テンプレチップ。タップでタイトル・種類を渡す。

export interface TaskTemplate {
  title: string;
  type: "deadline" | "schedule";
}

export const TASK_TEMPLATES: TaskTemplate[] = [
  { title: "ES提出", type: "deadline" },
  { title: "Webテスト", type: "deadline" },
  { title: "適性検査", type: "deadline" },
  { title: "一次面接", type: "schedule" },
  { title: "二次面接", type: "schedule" },
  { title: "最終面接", type: "schedule" },
  { title: "説明会", type: "schedule" },
  { title: "お礼メール", type: "deadline" },
];

interface Props {
  onSelect: (template: TaskTemplate) => void;
}

export function TaskTemplateChips({ onSelect }: Props) {
  return (
    <div className="-mx-1 overflow-x-auto px-1">
      <div className="flex min-w-max gap-1.5 pb-1">
        {TASK_TEMPLATES.map((template) => (
          <button
            key={template.title}
            type="button"
            onClick={() => onSelect(template)}
            className="shrink-0 rounded-full border border-line bg-cream px-2.5 py-1 text-[12px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            {template.title}
          </button>
        ))}
      </div>
    </div>
  );
}
