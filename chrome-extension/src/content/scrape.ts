// Content script — ページから会社名・職種・締切候補を抽出する。
// MVP では DOM のシグナルだけを集めて popup に渡す。
// 将来は AI 抽出 (server-side) に切り替え予定。

interface ScrapedJob {
  source: string;
  companyName: string;
  jobTitle: string;
  url: string;
  detectedAt: string;
}

const GENERIC_COMPANY_SELECTORS = [
  '[data-testid*="company" i]',
  '[data-test*="company" i]',
  '[class*="companyName" i]',
  '[class*="company-name" i]',
  '[class*="corpName" i]',
  '[class*="corp-name" i]',
  '[class*="company" i] h1',
  '[class*="corp" i] h1',
  "h1",
];

const GENERIC_TITLE_SELECTORS = [
  '[data-testid*="job" i] h1',
  '[data-testid*="title" i]',
  '[class*="jobTitle" i]',
  '[class*="job-title" i]',
  '[class*="title" i] h1',
  'meta[property="og:title"]',
  "h1",
];

const MAX_COMPANY_NAME_LENGTH = 80;

function detectFromMynavi(): Partial<ScrapedJob> | null {
  return {
    source: "マイナビ",
    companyName: pickCompanyName([
      ...jsonLdCompanyNames(),
      textFromSelectors([
        ".corpName",
        ".companyName",
        ".heading1",
        ...GENERIC_COMPANY_SELECTORS,
      ]),
      companyNameFromTitle(document.title),
      companyNameFromTitle(metaContent('meta[property="og:title"]')),
    ]),
    jobTitle: pickTitle(),
  };
}

function detectFromRikunabi(): Partial<ScrapedJob> | null {
  return {
    source: "リクナビ",
    companyName: pickCompanyName([
      ...jsonLdCompanyNames(),
      textFromSelectors([
        '[class*="rnn-company" i]',
        '[class*="companyName" i]',
        ...GENERIC_COMPANY_SELECTORS,
      ]),
      companyNameFromTitle(document.title),
      companyNameFromTitle(metaContent('meta[property="og:title"]')),
    ]),
    jobTitle: pickTitle(),
  };
}

function detectFromOneCareer(): Partial<ScrapedJob> | null {
  return {
    source: "ONE CAREER",
    companyName: pickCompanyName([
      ...jsonLdCompanyNames(),
      textFromSelectors([
        '[class*="company" i]',
        ...GENERIC_COMPANY_SELECTORS,
      ]),
      companyNameFromTitle(document.title),
      companyNameFromTitle(metaContent('meta[property="og:title"]')),
    ]),
    jobTitle: pickTitle(),
  };
}

function detectFromOfferBox(): Partial<ScrapedJob> | null {
  return {
    source: "OfferBox",
    companyName: pickCompanyName([
      ...jsonLdCompanyNames(),
      textFromSelectors([
        '[class*="company" i]',
        ...GENERIC_COMPANY_SELECTORS,
      ]),
      companyNameFromTitle(document.title),
      companyNameFromTitle(metaContent('meta[property="og:title"]')),
    ]),
    jobTitle: pickTitle(),
  };
}

function pickTitle(): string {
  return (
    cleanTitle(textFromSelectors(GENERIC_TITLE_SELECTORS)) ||
    cleanTitle(metaContent('meta[property="og:title"]')) ||
    cleanTitle(document.title)
  );
}

function pickCompanyName(candidates: Array<string | undefined>): string {
  for (const candidate of candidates) {
    const cleaned = cleanCompanyName(candidate);
    if (cleaned) return cleaned;
  }
  return "";
}

function textFromSelectors(selectors: string[]): string {
  for (const selector of selectors) {
    const node = document.querySelector(selector);
    if (!node) continue;
    if (node instanceof HTMLMetaElement) {
      const content = normalizeText(node.content);
      if (content) return content;
    }
    const text = normalizeText(node.textContent ?? "");
    if (text) return text;
  }
  return "";
}

function metaContent(selector: string): string {
  const meta = document.querySelector<HTMLMetaElement>(selector);
  return normalizeText(meta?.content ?? "");
}

function jsonLdCompanyNames(): string[] {
  const names: string[] = [];
  for (const script of document.querySelectorAll<HTMLScriptElement>('script[type="application/ld+json"]')) {
    const raw = script.textContent?.trim();
    if (!raw) continue;
    try {
      collectCompanyNamesFromJson(JSON.parse(raw), names);
    } catch {
      // JSON-LD はサイト側の壊れた断片もあるので、DOM/meta/title の候補にフォールバックする。
    }
  }
  return names;
}

function collectCompanyNamesFromJson(value: unknown, names: string[]): void {
  if (!value) return;
  if (Array.isArray(value)) {
    for (const item of value) collectCompanyNamesFromJson(item, names);
    return;
  }
  if (typeof value !== "object") return;

  const obj = value as Record<string, unknown>;
  const hiringOrganization = obj.hiringOrganization ?? obj.organization;
  if (hiringOrganization) collectCompanyNamesFromJson(hiringOrganization, names);

  const type = obj["@type"];
  const isCompanyLike =
    type === "Organization" ||
    type === "Corporation" ||
    type === "LocalBusiness" ||
    (Array.isArray(type) && type.some((t) => t === "Organization" || t === "Corporation"));
  if (isCompanyLike && typeof obj.name === "string") {
    names.push(obj.name);
  }
}

function companyNameFromTitle(title: string | undefined): string {
  const normalized = normalizeText(title ?? "");
  if (!normalized) return "";
  const parts = normalized
    .split(/[|｜]/)
    .flatMap((part) => part.split(/\s[-–—]\s/))
    .map(cleanCompanyName)
    .filter(Boolean);
  return (
    parts.find((part) => /株式会社|有限会社|合同会社|会社|Corporation|Inc\.?|Co\.?/i.test(part)) ??
    parts.find((part) => !isRecruitingNoise(part)) ??
    ""
  );
}

function cleanTitle(value: string | undefined): string {
  return normalizeText(value ?? "")
    .replace(/\s*[|｜]\s*.*$/, "")
    .trim();
}

function cleanCompanyName(value: string | undefined): string {
  const normalized = normalizeText(value ?? "");
  if (!normalized || isRecruitingNoise(normalized)) return "";

  const cleaned = normalized
    .replace(/\s*[|｜]\s*(マイナビ|リクナビ|ONE CAREER|OfferBox).*$/i, "")
    .replace(/\s*[-–—]\s*(新卒採用|採用情報|求人|募集要項|企業情報).*$/i, "")
    .replace(/^(企業名|会社名|社名)\s*[:：]\s*/, "")
    .replace(/\s*(の)?(新卒採用|採用情報|求人情報|募集要項|インターンシップ情報).*$/i, "")
    .replace(/\s+/g, " ")
    .trim();

  if (!cleaned || isRecruitingNoise(cleaned)) return "";
  if (cleaned.length > MAX_COMPANY_NAME_LENGTH) return "";

  return cleaned;
}

function isRecruitingNoise(value: string): boolean {
  const text = normalizeText(value);
  if (!text) return true;
  if (/^(採用情報|求人情報|募集要項|企業情報|エントリー|ログイン|マイページ)$/i.test(text)) return true;
  return /^(マイナビ|リクナビ|ONE CAREER|OfferBox)\s*\d*$/i.test(text);
}

function normalizeText(value: string): string {
  return value.replace(/\s+/g, " ").trim();
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
    url: window.location.href,
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
