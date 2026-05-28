import { useEffect, useState } from "react";
import { Mascot } from "../components/Mascot";
import { Confetti } from "../components/Confetti";
import { ApiRequestError, WEB_BASE, createInboxClip } from "../lib/api";

interface DetectedPage {
  source: string;
  companyGuess: string;
  title: string;
  url: string;
}

// 保存失敗時の表示。message は必須、recovery があれば回復ボタンを出す。
interface ErrorState {
  message: string;
  recovery?: { label: string; onClick: () => void };
}

function openWebLogin() {
  const url = `${WEB_BASE}/login`;
  if (typeof chrome !== "undefined" && chrome.tabs?.create) {
    void chrome.tabs.create({ url });
  } else {
    window.open(url, "_blank", "noopener");
  }
}

// 保存エラーをユーザー向けメッセージ + 回復導線に変換する。
function toErrorState(e: unknown): ErrorState {
  const kind = e instanceof ApiRequestError ? e.kind : "client";
  switch (kind) {
    case "unauthorized":
      return {
        message:
          "Entré にログインしていません。Web アプリでログインしてから、もう一度保存してください。",
        recovery: { label: "Web でログインする", onClick: openWebLogin },
      };
    case "forbidden":
      return {
        message:
          "アクセスが拒否されました。一度ログアウト／再ログインするか、拡張の許可設定を管理者に確認してください。",
        recovery: { label: "Web を開く", onClick: openWebLogin },
      };
    case "network":
      return {
        message:
          "サーバーに接続できませんでした。ネット接続と、Entré サーバー／拡張の許可ドメイン設定を確認してください。",
      };
    case "server":
      return {
        message:
          "サーバーでエラーが発生しました。少し時間をおいて、もう一度お試しください。",
      };
    default:
      return {
        message: "保存に失敗しました。もう一度お試しください。",
      };
  }
}

interface ScrapedPage {
  source?: string;
  companyName?: string;
  jobTitle?: string;
  url?: string;
}

interface ScrapeResponse {
  data?: ScrapedPage | null;
}

const SOURCES: Record<string, string> = {
  "mynavi.jp": "マイナビ",
  "rikunabi.com": "リクナビ",
  "onecareer.jp": "ONE CAREER",
  "offerbox.jp": "OfferBox",
};

export function Popup() {
  const [page, setPage] = useState<DetectedPage | null>(null);
  const [companyGuess, setCompanyGuess] = useState("");
  const [saving, setSaving] = useState(false);
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<ErrorState | null>(null);

  useEffect(() => {
    let mounted = true;
    void detectCurrentPage().then((detected) => {
      if (!mounted) return;
      setPage(detected);
      setCompanyGuess(detected?.companyGuess ?? "");
    });
    return () => {
      mounted = false;
    };
  }, []);

  const handleSave = async () => {
    if (!page || saving) return;
    setError(null);
    setSaving(true);
    try {
      await createInboxClip({
        url: page.url,
        title: page.title,
        source: page.source,
        guess: companyGuess.trim() || undefined,
      });
      setConfetti((n) => n + 1);
      setSaving(false);
      window.setTimeout(() => window.close(), 1500);
    } catch (e) {
      // 失敗時は popup を閉じず、回復導線つきのエラーを表示する。
      setSaving(false);
      setError(toErrorState(e));
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
          <DetectedCard page={page} guess={companyGuess} onGuessChange={setCompanyGuess} />
        ) : (
          <div className="rounded-[10px] border border-line bg-surface p-4 text-center text-[11px] text-ink-3">
            このページは保存できません
          </div>
        )}

        {error && (
          <div
            role="alert"
            className="mt-2 rounded-md bg-pink/40 px-2.5 py-2 text-[10px] font-semibold text-ink"
          >
            <p className="leading-relaxed">{error.message}</p>
            {error.recovery && (
              <button
                type="button"
                onClick={error.recovery.onClick}
                className="mt-1.5 rounded-md bg-sage px-2.5 py-1 text-[10px] font-bold text-white transition-transform hover:-translate-y-0.5"
              >
                {error.recovery.label}
              </button>
            )}
          </div>
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
          {saving ? "保存中..." : confetti ? "✓ 保存しました！" : "＋ Inbox に保存"}
        </button>
      </footer>

      <Confetti trigger={confetti} count={16} />
    </div>
  );
}

function DetectedCard({
  page,
  guess,
  onGuessChange,
}: {
  page: DetectedPage;
  guess: string;
  onGuessChange: (value: string) => void;
}) {
  return (
    <div className="rounded-[10px] border border-line bg-surface p-3">
      <div className="mb-2 flex items-center justify-between gap-2">
        <div className="text-[9px] text-ink-3">検出されたページ</div>
        <span className="rounded-sm bg-cream-2 px-1.5 py-0.5 text-[9px] font-bold text-ink-2">
          {page.source}
        </span>
      </div>

      <label className="mb-1 block text-[10px] font-bold text-ink-2" htmlFor="company-guess">
        企業名候補
      </label>
      <input
        id="company-guess"
        value={guess}
        onChange={(e) => onGuessChange(e.target.value)}
        placeholder="未入力"
        className="mb-2 block h-9 w-full rounded-lg border border-line bg-cream px-2.5 font-sans text-[12px] font-bold text-ink outline-none focus:border-sage"
      />

      <div className="mb-1 text-[10px] font-bold text-ink-2">ページ名</div>
      <div className="line-clamp-2 text-[11px] leading-snug text-ink-2">{page.title}</div>
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
      companyGuess: "株式会社○○商事",
      title: "株式会社○○商事 / 総合職 新卒採用 2026",
      url: "https://job.mynavi.jp/26/pc/search/corp123/outline.html",
    };
  }

  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  if (!tab?.url) return null;

  const url = new URL(tab.url);
  if (!isSavableUrl(url)) return null;

  const scraped = typeof tab.id === "number" ? await scrapeCurrentTab(tab.id) : null;
  const title = normalizeText(scraped?.jobTitle) || normalizeText(tab.title) || url.href;
  const companyGuess =
    normalizeText(scraped?.companyName) || normalizeText(firstTitlePart(title));
  return {
    source: normalizeText(scraped?.source) || detectSource(url.hostname),
    companyGuess,
    title,
    url: normalizeHttpUrl(scraped?.url) ?? tab.url,
  };
}

async function scrapeCurrentTab(tabId: number): Promise<ScrapedPage | null> {
  return new Promise((resolve) => {
    chrome.tabs.sendMessage(
      tabId,
      { type: "ENTRE_SCRAPE_REQUEST" },
      (response: ScrapeResponse | undefined) => {
        if (chrome.runtime.lastError) {
          resolve(null);
          return;
        }
        resolve(response?.data ?? null);
      },
    );
  });
}

function isSavableUrl(url: URL): boolean {
  return url.protocol === "http:" || url.protocol === "https:";
}

function detectSource(hostname: string): string {
  const sourceKey = Object.keys(SOURCES).find(
    (k) => hostname === k || hostname.endsWith(`.${k}`),
  );
  return sourceKey ? SOURCES[sourceKey] : hostname.replace(/^www\./, "");
}

function firstTitlePart(title: string): string {
  return title.split(/[|｜]/)[0]?.trim() ?? "";
}

function normalizeText(value: string | undefined): string {
  return (value ?? "").replace(/\s+/g, " ").trim();
}

function normalizeHttpUrl(value: string | undefined): string | null {
  if (!value) return null;
  try {
    const url = new URL(value);
    return isSavableUrl(url) ? url.href : null;
  } catch {
    return null;
  }
}
