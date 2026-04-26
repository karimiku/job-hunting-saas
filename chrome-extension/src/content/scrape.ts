// Content script — ページから会社名・職種・締切候補を抽出する。
// MVP では DOM のシグナルだけを集めて popup に渡す。
// 将来は AI 抽出 (server-side) に切り替え予定。

interface ScrapedJob {
  source: string;
  companyName: string;
  jobTitle: string;
  detectedAt: string;
}

function detectFromMynavi(): Partial<ScrapedJob> | null {
  // job.mynavi.jp 系の構造 — h1 や og:title から会社名を取る
  const h1 = document.querySelector("h1");
  const og = document.querySelector<HTMLMetaElement>('meta[property="og:title"]');
  return {
    source: "マイナビ",
    companyName: h1?.textContent?.trim() || og?.content?.split("|")[0]?.trim() || "",
    jobTitle: document.title || "",
  };
}

function detectFromRikunabi(): Partial<ScrapedJob> | null {
  return {
    source: "リクナビ",
    companyName: document.querySelector("h1")?.textContent?.trim() ?? "",
    jobTitle: document.title || "",
  };
}

function detectFromOneCareer(): Partial<ScrapedJob> | null {
  return {
    source: "ONE CAREER",
    companyName: document.querySelector("h1")?.textContent?.trim() ?? "",
    jobTitle: document.title || "",
  };
}

function detectFromOfferBox(): Partial<ScrapedJob> | null {
  return {
    source: "OfferBox",
    companyName: document.querySelector("h1")?.textContent?.trim() ?? "",
    jobTitle: document.title || "",
  };
}

/** host が target ドメイン (またはそのサブドメイン) にマッチするか厳密にチェックする。
 *  `String.includes` だと "evil.com/mynavi.jp/..." のような擬装にマッチしてしまうため、
 *  hostname の末尾一致 + ドット境界で評価する。 */
function hostMatches(host: string, target: string): boolean {
  return host === target || host.endsWith(`.${target}`);
}

function detectCurrentPage(): ScrapedJob | null {
  const host = window.location.hostname;
  let partial: Partial<ScrapedJob> | null = null;

  if (hostMatches(host, "mynavi.jp")) partial = detectFromMynavi();
  else if (hostMatches(host, "rikunabi.com")) partial = detectFromRikunabi();
  else if (hostMatches(host, "onecareer.jp")) partial = detectFromOneCareer();
  else if (hostMatches(host, "offerbox.jp")) partial = detectFromOfferBox();

  if (!partial?.companyName) return null;

  return {
    source: partial.source ?? "",
    companyName: partial.companyName,
    jobTitle: partial.jobTitle ?? "",
    detectedAt: new Date().toISOString(),
  };
}

// popup 側からのリクエストに応答
chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (message?.type === "ENTRE_SCRAPE_REQUEST") {
    sendResponse({ data: detectCurrentPage() });
    return true;
  }
  return undefined;
});

// 初回ロード時に自動検出してストレージに保存（バッジ表示など将来用）
const detected = detectCurrentPage();
if (detected) {
  void chrome.storage.local.set({ "entre:lastDetected": detected });
}
