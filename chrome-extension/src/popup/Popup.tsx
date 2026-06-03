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
  openUrl(url);
}

function openWebInbox() {
  const url = `${WEB_BASE}/inbox`;
  openUrl(url);
}

function openWebDashboard() {
  const url = `${WEB_BASE}/dashboard`;
  openUrl(url);
}

function openWebEntries() {
  const url = `${WEB_BASE}/entry`;
  openUrl(url);
}

function openUrl(url: string) {
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
  "i-web.jpn.com": "i-web",
  "supporterz.jp": "サポーターズ",
  "wantedly.com": "Wantedly",
  "hrmos.co": "HRMOS",
  "green-japan.com": "Green",
};

const AUTO_OPEN_INBOX_KEY = "entre.autoOpenInbox";

export function Popup() {
  const [page, setPage] = useState<DetectedPage | null>(null);
  const [companyGuess, setCompanyGuess] = useState("");
  const [saving, setSaving] = useState(false);
  const [confetti, setConfetti] = useState(0);
  const [error, setError] = useState<ErrorState | null>(null);
  const [saved, setSaved] = useState(false);
  const [autoOpenInbox, setAutoOpenInbox] = useState(false);

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

  useEffect(() => {
    if (typeof chrome === "undefined" || !chrome.storage?.local) return;
    chrome.storage.local.get(AUTO_OPEN_INBOX_KEY, (items) => {
      setAutoOpenInbox(Boolean(items[AUTO_OPEN_INBOX_KEY]));
    });
  }, []);

  const handleAutoOpenInboxChange = (checked: boolean) => {
    setAutoOpenInbox(checked);
    if (typeof chrome !== "undefined" && chrome.storage?.local) {
      void chrome.storage.local.set({ [AUTO_OPEN_INBOX_KEY]: checked });
    }
  };

  const handleSave = async () => {
    if (!page || saving) return;
    setError(null);
    setSaved(false);
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
      setSaved(true);
      if (autoOpenInbox) openWebInbox();
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
        <div className="mb-2 grid grid-cols-3 gap-1.5">
          <button
            type="button"
            onClick={openWebDashboard}
            className="rounded-md border border-line bg-surface px-2 py-1.5 text-[10px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            Dashboard
          </button>
          <button
            type="button"
            onClick={openWebInbox}
            className="rounded-md border border-line bg-surface px-2 py-1.5 text-[10px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            Inbox
          </button>
          <button
            type="button"
            onClick={openWebEntries}
            className="rounded-md border border-line bg-surface px-2 py-1.5 text-[10px] font-bold text-ink-2 transition-colors hover:border-sage hover:text-sage"
          >
            Entry
          </button>
        </div>

        {saved && (
          <div className="mb-2 rounded-lg border border-sage bg-sage-wash p-3">
            <div className="flex items-start gap-2">
              <Mascot size={34} mood="happy" />
              <div className="min-w-0 flex-1">
                <p className="text-[12px] font-black text-sage">Inbox に保存しました</p>
                <p className="mt-0.5 text-[10px] leading-relaxed text-ink-2">
                  Web の Inbox で会社名を確認すると、Entry・Kanban・Task で管理できます。
                </p>
                <button
                  type="button"
                  onClick={openWebInbox}
                  className="mt-2 rounded-md bg-sage px-2.5 py-1 text-[10px] font-bold text-white transition-transform hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-sage/40"
                >
                  Web Inbox を開く
                </button>
              </div>
            </div>
          </div>
        )}

        {page ? (
          <DetectedCard page={page} guess={companyGuess} onGuessChange={setCompanyGuess} />
        ) : (
          <div className="rounded-[10px] border border-line bg-surface p-4 text-center">
            <div className="font-serif text-[14px] font-extrabold">保存できるページが見つかりません</div>
            <p className="mt-1 text-[10px] leading-relaxed text-ink-3">
              http/https の求人ページで開き直してから、もう一度 Entré を開いてください。
            </p>
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

        <label className="mt-2 flex cursor-pointer items-center justify-between gap-2 rounded-md border border-line bg-surface px-2.5 py-2">
          <span className="text-[10px] font-bold text-ink-2">保存後に Web Inbox を開く</span>
          <input
            type="checkbox"
            checked={autoOpenInbox}
            onChange={(event) => handleAutoOpenInboxChange(event.target.checked)}
            className="h-3.5 w-3.5 accent-sage"
          />
        </label>
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
          disabled={!page || saving || saved}
          className="flex-[2] rounded-lg bg-sage px-2.5 py-2 text-[11px] font-bold text-white transition-transform enabled:hover:-translate-y-0.5 disabled:opacity-60"
        >
          {saving ? "保存中..." : saved ? "保存済み" : "Inbox に保存"}
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
        <div>
          <div className="text-[9px] text-ink-3">検出されたページ</div>
          <div className="text-[10px] font-bold text-ink-2">保存前に会社名を確認</div>
        </div>
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
        className="mb-1 block h-9 w-full rounded-lg border border-line bg-cream px-2.5 font-sans text-[12px] font-bold text-ink outline-none focus:border-sage focus:ring-2 focus:ring-sage/20"
      />
      <p className="mb-2 text-[9px] leading-relaxed text-ink-3">
        ここで直した会社名候補が Web Inbox の Entry 作成フォームに入ります。
      </p>

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
