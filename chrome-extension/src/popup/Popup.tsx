import { useEffect, useState } from "react";
import { Mascot } from "../components/Mascot";
import { Confetti } from "../components/Confetti";
import { saveDetectedJob } from "../lib/api";

interface DetectedPage {
  source: string;
  companyName: string;
  jobTitle: string;
  url: string;
}

const STAGES = [
  { key: "entry", label: "エントリー", color: "#C9CBB4" },
  { key: "doc", label: "書類選考", color: "#A8C0DA" },
  { key: "es", label: "ES提出", color: "#D4BA82" },
  { key: "interview", label: "面接", color: "#E9B9B0" },
  { key: "offer", label: "内定", color: "#9BB58A" },
] as const;

const SOURCES: Record<string, string> = {
  "mynavi.jp": "マイナビ",
  "rikunabi.com": "リクナビ",
  "onecareer.jp": "ONE CAREER",
  "offerbox.jp": "OfferBox",
};

export function Popup() {
  const [page, setPage] = useState<DetectedPage | null>(null);
  const [stageIdx, setStageIdx] = useState(0);
  const [memo, setMemo] = useState("");
  const [saving, setSaving] = useState(false);
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    void detectCurrentPage().then(setPage);
  }, []);

  const handleSave = async () => {
    if (!page || saving) return;
    setError(null);
    setSaving(true);
    try {
      await saveDetectedJob({
        companyName: page.companyName,
        route: "本選考",
        source: page.source,
        memo,
      });
      setConfetti((n) => n + 1);
      setSaving(false);
      window.setTimeout(() => window.close(), 1500);
    } catch (e) {
      setSaving(false);
      setError(e instanceof Error ? e.message : "保存に失敗しました");
    }
  };

  return (
    <div className="relative flex h-[440px] w-[360px] flex-col overflow-hidden border border-line bg-cream font-sans text-ink shadow-lg">
      {/* Header */}
      <header className="flex items-center gap-2 border-b border-line bg-surface px-3.5 py-3">
        <span className="font-serif text-[18px] font-black italic">Entré</span>
        <span className="flex-1 text-[10px] text-ink-3">このページから保存</span>
        <Mascot size={26} />
      </header>

      {/* Body */}
      <div className="flex-1 overflow-auto p-3.5">
        {page ? (
          <DetectedCard page={page} />
        ) : (
          <div className="rounded-[10px] border border-line bg-surface p-4 text-center text-[11px] text-ink-3">
            このページからはエントリーを検出できません
          </div>
        )}

        <div className="mt-3 mb-1.5 text-[10px] font-bold text-ink-2">ステータス</div>
        <div className="mb-3 flex gap-1">
          {STAGES.map((s, i) => (
            <button
              key={s.key}
              type="button"
              onClick={() => setStageIdx(i)}
              className="flex-1 rounded-md py-1 text-[8px] font-bold transition-colors"
              style={{
                background: i === stageIdx ? s.color : "var(--color-surface)",
                color: i === stageIdx ? "#fff" : "var(--color-ink-2)",
                border: `1px solid ${i === stageIdx ? s.color : "var(--color-line)"}`,
              }}
              aria-pressed={i === stageIdx}
            >
              {s.label}
            </button>
          ))}
        </div>

        <div className="mb-1.5 text-[10px] font-bold text-ink-2">メモ</div>
        <textarea
          value={memo}
          onChange={(e) => setMemo(e.target.value)}
          placeholder="気になるポイントなど"
          className="block min-h-[50px] w-full resize-none rounded-lg border border-line bg-surface px-2.5 py-2 font-sans text-[10px] text-ink-2 outline-none focus:border-sage"
        />

        {error && (
          <p className="mt-2 rounded-md bg-pink/40 px-2.5 py-1.5 text-[10px] font-semibold text-ink">
            {error}
          </p>
        )}
      </div>

      {/* Footer */}
      <footer className="flex gap-2 border-t border-line bg-surface px-3.5 py-2.5">
        <button
          type="button"
          onClick={() => window.close()}
          className="flex-1 rounded-lg border border-line bg-transparent px-2.5 py-2 text-[11px] font-bold text-ink-2 transition-colors hover:bg-line"
        >
          キャンセル
        </button>
        <button
          type="button"
          onClick={handleSave}
          disabled={!page || saving}
          className="flex-[2] rounded-lg bg-sage px-2.5 py-2 text-[11px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
        >
          {saving ? "保存中..." : confetti ? "✓ 保存しました！" : "＋ Entré に保存"}
        </button>
      </footer>

      <Confetti trigger={confetti} count={16} />
    </div>
  );
}

function DetectedCard({ page }: { page: DetectedPage }) {
  return (
    <div className="rounded-[10px] border border-line bg-surface p-3">
      <div className="mb-1 text-[9px] text-ink-3">
        検出されたページ · {page.source}
      </div>
      <div className="mb-1.5 text-[13px] font-extrabold leading-tight">
        {page.companyName}
      </div>
      <div className="text-[11px] text-ink-2">{page.jobTitle}</div>
      <div className="mt-2 truncate font-mono text-[9px] text-ink-3">{page.url}</div>
    </div>
  );
}

/** 現在のタブから会社名/ソースを推定する。 */
async function detectCurrentPage(): Promise<DetectedPage | null> {
  if (typeof chrome === "undefined" || !chrome.tabs?.query) {
    // 開発用フォールバック (vite dev server で popup を直接開いたとき)
    return {
      source: "マイナビ",
      companyName: "株式会社○○商事 / 総合職 新卒採用 2026",
      jobTitle: "総合職・新卒",
      url: "https://job.mynavi.jp/26/pc/search/corp123/outline.html",
    };
  }

  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  if (!tab?.url) return null;

  const url = new URL(tab.url);
  // host の末尾一致 + ドット境界で厳密マッチ (incomplete-url-substring-sanitization 対策)
  const sourceKey = Object.keys(SOURCES).find(
    (k) => url.hostname === k || url.hostname.endsWith(`.${k}`),
  );
  if (!sourceKey) return null;

  return {
    source: SOURCES[sourceKey],
    companyName: tab.title?.split("|")[0]?.trim() ?? "（タイトル不明）",
    jobTitle: tab.title ?? "",
    url: tab.url,
  };
}
