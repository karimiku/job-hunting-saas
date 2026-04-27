# Entré Chrome Extension

マイナビ・リクナビ・ONE CAREER・OfferBox から1クリックで Entré にエントリーを保存する Chrome 拡張。

## 開発

```bash
pnpm install
pnpm dev      # Vite dev server (HMR つき) — popup を直接ブラウザで開けるが Chrome API は使えない
pnpm build    # 本番ビルド → dist/
```

## ロード手順 (Chrome / Chromium 系)

1. `pnpm build` で `dist/` を生成
2. `chrome://extensions` を開く
3. 「デベロッパーモード」を ON
4. 「パッケージ化されていない拡張機能を読み込む」→ `chrome-extension/dist/` を選択
5. ツールバーの 封筒くん アイコンが表示される

## 使い方

1. 対応サイト (マイナビ等) で気になる求人ページを開く
2. ツールバーの Entré アイコンをクリック → ポップアップが開く
3. 検出された会社名・ステータス・メモを確認
4. 「＋ Entré に保存」 → バックエンドに保存（要ログイン状態）

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

## TODO

- [ ] AI 抽出 — 現状は DOM の h1 / og:title をそのまま読んでいる。サーバー側で
  ページ HTML から会社名・職種・締切を抽出する API に切り替える
- [ ] 重複検出 — 同じ URL を保存しようとしたら既存エントリーを表示する
- [ ] Inbox 連携 — 検出済みだが保存していないクリップを Inbox に溜める
- [ ] ストア公開準備 — privacy policy, store listing 画像など
