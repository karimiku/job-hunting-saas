// プロトタイプ用のダミーデータ。バックエンド API 接続時に置き換える。

export type StageKey = "entry" | "doc" | "es" | "interview" | "offer";

export interface Entry {
  id: string;
  co: string;
  logo: string;
  stage: StageKey;
  stageLabel: string;
  stageIdx: number;
  color: string;
  due: string;
  task: string;
  fav: boolean;
}

export const ENTRIES: Entry[] = [
  { id: "e1", co: "○○商事", logo: "○", stage: "interview", stageLabel: "面接", stageIdx: 3, color: "var(--color-stage-interview)", due: "5/29 14:00", task: "一次面接(Web)", fav: true },
  { id: "e2", co: "△△株式会社", logo: "△", stage: "es", stageLabel: "ES提出", stageIdx: 2, color: "var(--color-stage-es)", due: "5/28 23:59", task: "エントリーシート", fav: false },
  { id: "e3", co: "□□□コーポレーション", logo: "□", stage: "offer", stageLabel: "内定", stageIdx: 4, color: "var(--color-stage-offer)", due: "6/02 17:00", task: "内定承諾返事", fav: true },
  { id: "e4", co: "◇◇テック", logo: "◇", stage: "doc", stageLabel: "書類選考", stageIdx: 1, color: "var(--color-stage-doc)", due: "5/30 12:00", task: "適性検査", fav: false },
  { id: "e5", co: "☆☆ホールディングス", logo: "☆", stage: "entry", stageLabel: "エントリー", stageIdx: 0, color: "var(--color-stage-entry)", due: "6/05 23:59", task: "プレエントリー", fav: false },
  { id: "e6", co: "▲▲メディア", logo: "▲", stage: "interview", stageLabel: "面接", stageIdx: 3, color: "var(--color-stage-interview)", due: "6/01 10:00", task: "最終面接", fav: true },
];

export interface StageDef {
  key: StageKey;
  label: string;
  color: string;
}

export const STAGES: StageDef[] = [
  { key: "entry", label: "エントリー", color: "var(--color-stage-entry)" },
  { key: "doc", label: "書類選考", color: "var(--color-stage-doc)" },
  { key: "es", label: "ES提出", color: "var(--color-stage-es)" },
  { key: "interview", label: "面接", color: "var(--color-stage-interview)" },
  { key: "offer", label: "内定", color: "var(--color-stage-offer)" },
];

export const STAGE_COUNTS: Record<StageKey, number> = {
  entry: 8,
  doc: 4,
  es: 5,
  interview: 5,
  offer: 2,
};

export const KANBAN_CARDS: Record<StageKey, { co: string; l: string; d: string }[]> = {
  entry: [
    { co: "☆☆HD", l: "☆", d: "6/05" },
    { co: "××商事", l: "×", d: "6/08" },
  ],
  doc: [
    { co: "◇◇テック", l: "◇", d: "5/30" },
    { co: "※※工業", l: "※", d: "5/31" },
    { co: "#### Inc.", l: "#", d: "6/01" },
  ],
  es: [
    { co: "△△株式会社", l: "△", d: "5/28" },
    { co: "@@商事", l: "@", d: "5/29" },
    { co: "$$ HD", l: "$", d: "5/30" },
    { co: "%%メディア", l: "%", d: "6/01" },
  ],
  interview: [
    { co: "○○商事", l: "○", d: "5/29" },
    { co: "▲▲メディア", l: "▲", d: "5/30" },
    { co: "¥¥銀行", l: "¥", d: "6/02" },
    { co: "!!! Tech", l: "!", d: "6/03" },
    { co: "&&商社", l: "&", d: "6/05" },
  ],
  offer: [
    { co: "□□□コーポ", l: "□", d: "6/02 返答" },
    { co: "~~~~~ Co.", l: "~", d: "保留中" },
  ],
};

export interface Milestone {
  l: string;
  n: number;
  c: string;
  emoji: string;
  done?: boolean;
  current?: boolean;
}

export const MILESTONES: Milestone[] = [
  { l: "エントリー", n: 8, c: "var(--color-stage-entry)", emoji: "📝", done: true },
  { l: "書類選考", n: 4, c: "var(--color-stage-doc)", emoji: "📄", done: true },
  { l: "ES提出", n: 5, c: "var(--color-stage-es)", emoji: "✏️", done: true },
  { l: "面接", n: 5, c: "var(--color-stage-interview)", emoji: "🤝", current: true },
  { l: "内定", n: 2, c: "var(--color-stage-offer)", emoji: "🎉" },
];
