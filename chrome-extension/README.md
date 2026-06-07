# Entré Chrome Extension

マイナビ・リクナビ・ONE CAREER・OfferBox などから1クリックで Entré の Inbox にページを保存する Chrome 拡張。

## 開発

```bash
pnpm install
pnpm dev      # Vite dev server (HMR つき) — popup を直接ブラウザで開けるが Chrome API は使えない
pnpm build    # 本番ビルド → dist/
```

## 環境変数

接続先は環境変数で切り替えられる（ビルド時に Vite が埋め込む）。未設定時はローカル開発のデフォルトを使う。

| 変数 | 用途 | デフォルト |
| --- | --- | --- |
| `VITE_API_BASE_URL` | Entré バックエンド API のベース URL | `http://localhost:8080` |
| `VITE_WEB_BASE_URL` | Web アプリ（未ログイン時に開くログインページ）のベース URL | `http://localhost:3000` |

```bash
# 例: 本番向けにビルド
VITE_API_BASE_URL=https://api.example.com VITE_WEB_BASE_URL=https://app.example.com pnpm build
```

## ロード手順 (Chrome / Chromium 系)

1. `pnpm build` で `dist/` を生成
2. `chrome://extensions` を開く
3. 「デベロッパーモード」を ON
4. 「パッケージ化されていない拡張機能を読み込む」→ `chrome-extension/dist/` を選択
5. ツールバーの 封筒くん アイコンが表示される

## 使い方

1. 気になる求人ページを開く
2. ツールバーの Entré アイコンをクリック → ポップアップが開く
3. 検出された企業名候補を確認
4. 「＋ Inbox に保存」 → Webアプリの Inbox に保存（要ログイン状態）

## 保存に失敗したとき

popup は失敗しても閉じず、原因別の回復案内を表示する。

| 状況 | 表示 | 回復導線 |
| --- | --- | --- |
| 未ログイン (401) | Web ログインが必要 | 「Web でログインする」ボタンで `VITE_WEB_BASE_URL/login` を開く |
| 権限エラー (403) | 再ログイン/許可設定の確認 | Web を開く |
| 接続エラー (サーバー停止 / CORS / オフライン) | 接続・許可ドメイン設定の確認 | — |
| サーバーエラー (5xx) | 時間をおいて再試行 | — |

## アーキテクチャ

```
src/
├── popup/                  # ツールバーアイコンをクリックしたときの 360×440 popup
│   ├── index.html
│   ├── main.tsx
│   └── Popup.tsx           # 設計の ChromeExtension コンポーネントを React で実装
├── content/
│   └── scrape.ts           # マイナビ等のページから DOM を読み取る content script
├── components/
│   ├── Mascot.tsx          # 封筒くん (frontend と共有実装)
│   └── Confetti.tsx        # 保存成功時の紙吹雪
├── lib/
│   └── api.ts              # Entré バックエンドへの fetch クライアント
└── styles/
    └── popup.css           # Tailwind v4 + A案 デザイントークン
```

## 認証

popup から `credentials: "include"` で Entré API を叩く。
事前に Entré 本体 (`localhost:3000`) でログイン済みであることが前提。
Session Cookie の `SameSite` 属性は現状 `Lax` のため、host_permissions に
バックエンド (`localhost:8080`) を含めて拡張から共有できるようにしている。

## Roadmap

- [ ] AI 抽出 — 現状は DOM の h1 / og:title をそのまま読んでいる。サーバー側で
  ページ HTML から会社名・職種・締切を抽出する API に切り替える
- [ ] 重複検出 — 同じ URL を保存しようとしたら既存エントリーを表示する
- [x] Inbox 連携 — 検出済みページを Inbox に保存する
- [x] Inbox 整理 — 保存したクリップから Web の Inbox で Entry を作成・紐付ける
- [ ] ストア公開準備 — privacy policy, store listing 画像など
