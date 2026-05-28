// Server Component。実データ（面接中・内定・未完了タスク件数）から励ましコピーを組み立てる。
// 固定文言だった旧実装を置き換える。useEffect は使わない。

import Link from "next/link";
import { Mascot } from "./Mascot";

export interface EncouragementCopy {
  hand: string;
  headline: string;
  body: string;
}

// 状況に応じてコピーを出し分ける。データが無ければ汎用の励ましにフォールバック。
export function buildEncouragement(args: {
  interviewing: number;
  offered: number;
  openTasks: number;
}): EncouragementCopy {
  const { interviewing, offered, openTasks } = args;

  if (offered > 0) {
    return {
      hand: "やりましたね！",
      headline: `内定 ${offered} 件、本当におめでとうございます。`,
      body: "ここまでの積み重ねが実を結びました。次の一歩も応援しています。",
    };
  }

  if (interviewing > 0) {
    return {
      hand: "あと少しですね！",
      headline: `面接 ${interviewing} 社、内定まであと一歩。`,
      body:
        openTasks > 0
          ? `未完了のタスクが ${openTasks} 件あります。ひとつずつ片付けていきましょう。`
          : "今日のタスクは片付きました。面接の準備に集中できますね。",
    };
  }

  if (openTasks > 0) {
    return {
      hand: "今日もコツコツ。",
      headline: `未完了のタスクが ${openTasks} 件あります。`,
      body: "まずは1件、できるところから始めてみましょう。応援しています！",
    };
  }

  return {
    hand: "ナイスペース！",
    headline: "今日のタスクはすべて完了しています。",
    body: "新しいエントリーを追加して、就活を一歩前へ進めましょう。",
  };
}

/** 実データから組み立てた励ましカード。 */
export function DashboardEncouragement(props: {
  interviewing: number;
  offered: number;
  openTasks: number;
}) {
  const copy = buildEncouragement(props);

  return (
    <div className="mt-4 flex flex-col items-start gap-4 rounded-xl border-[1.5px] border-line bg-gradient-to-br from-cream-2 to-sage-wash p-5 md:flex-row md:items-center md:p-6">
      <div style={{ animation: "entre-float 3s infinite" }}>
        <Mascot size={64} mood="cheering" />
      </div>
      <div className="flex-1">
        <p className="font-hand text-[18px] text-sage">{copy.hand}</p>
        <p data-testid="encouragement-headline" className="mt-0.5 font-serif text-base font-extrabold">
          {copy.headline}
        </p>
        <p className="mt-1 text-[11px] text-ink-2">{copy.body}</p>
      </div>
      <Link
        href="/roadmap"
        className="rounded-lg bg-sage px-3.5 py-2 text-[11px] font-bold text-white transition-transform hover:-translate-y-0.5"
      >
        ロードマップ →
      </Link>
    </div>
  );
}
