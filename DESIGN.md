---
brand: Entré
style: Quiet Core Workflow
themes: [light]
default_theme: light
references: [notion.so, attio.com, todoist.com, linear.app]
---

# DESIGN.md

## Design Philosophy

「Entry とカンバンを中心に、次の行動だけが分かる就活管理ツール」

- コア導線は Home / Entry / カンバン / タスク / 保存箱に絞る
- Entry は応募先の主語、カンバンは選考状況の俯瞰、タスクは今日やること、保存箱は Entry 候補の置き場
- 非コア機能に見える導線（ロードマップ、通知設定、カレンダー連携、拡張機能訴求）は画面上で主張させない
- 一覧・ボードは情報密度を保ち、フォームや空状態は次の1アクションだけを提示する
- 日常操作は装飾より速度と理解しやすさを優先する

## Colors

### Light Theme

| Token | Hex | 出典・意図 |
|-------|-----|-----------|
| cream | `#FBF8F1` | LP とアプリ共通の暖かいキャンバス |
| cream-2 | `#F6EFD8` | 控えめな面・選択背景 |
| surface | `#FFFFFF` | カード・フォーム・リスト行 |
| ink | `#2B2A26` | 見出し・主文 |
| ink-2 | `#5C5A52` | 補足文・サブ情報 |
| ink-3 | `#8F8D82` | メタ情報・非強調ラベル |
| line | `#E6E1D3` | 境界線 |
| sage | `#4F6E58` | 主CTA・選択状態・コア導線 |
| sage-soft | `#E5ECDE` | 主CTA以外の強調背景 |
| sage-wash | `#F0F3E8` | セクション/アイコン背景 |

## Typography

### フォントスタック

```css
/* UI全般 */
--font-sans: 'Noto Sans JP', 'Hiragino Sans', 'Yu Gothic', system-ui, sans-serif;

/* 見出し */
--font-serif: 'Noto Serif JP', 'Hiragino Mincho ProN', 'Yu Mincho', serif;

/* 数値・ID表示 */
--font-mono: 'JetBrains Mono', ui-monospace, 'SF Mono', Menlo, monospace;
```

### サイズスケール

| Token | Size | Line Height | Weight | 用途 |
|-------|------|-------------|--------|------|
| display | 32px | 40px | 700 | ダッシュボード見出し |
| title | 24px | 32px | 600 | ページタイトル |
| heading | 18px | 26px | 600 | セクション見出し |
| body | 15px | 24px | 400 | 本文・リスト |
| caption | 13px | 18px | 500 | ラベル・補足・日付 |
| micro | 11px | 16px | 500 | バッジ・ステータスタグ |

## Spacing

8px ベースのスケール。Attio の密度と Notion のゆとりの中間。

| Token | Value |
|-------|-------|
| 2xs | 4px |
| xs | 8px |
| sm | 12px |
| md | 16px |
| lg | 24px |
| xl | 32px |
| 2xl | 48px |
| 3xl | 64px |

## Border Radius

| Token | Value | 用途 |
|-------|-------|------|
| sm | 8px | バッジ・タグ |
| md | 12px | カード・入力フィールド |
| lg | 16px | モーダル・ドロップダウン |
| xl | 20px | 大きなカード |
| full | 9999px | ピル型ボタン・検索バー（← Airbnb） |

## Shadows

```css
--shadow-soft: 0 6px 18px -8px rgba(43, 42, 38, 0.12);
--shadow-card: 0 4px 12px -4px rgba(43, 42, 38, 0.08);
--shadow-fab: 0 10px 24px rgba(79, 110, 88, 0.4);
```

## Animation Guidelines

**原則: 毎日使うところは速く、装飾はコア理解を邪魔しない**

### 使う場面（← Amie, Todoist）
| 場面 | 種類 | Duration | Easing |
|------|------|----------|--------|
| ステージ変更 | カード移動 | 300ms | ease-out |
| タスク完了 | チェック + フェードアウト | 250ms | ease-in-out |
| ドラッグ&ドロップ | スムーズ追従 | 150ms | linear |

### 使わない場面（← Linear）
- ページ遷移 → 即時切替（トランジションなし）
- リスト表示・フィルタリング → 即時反映
- 日常的な繰り返し操作 → 速度最優先

### 共通値
```css
--transition-fast: 150ms ease-out;    /* hover, focus */
--transition-base: 250ms ease-in-out; /* 状態変化 */
--transition-slow: 400ms ease-out;    /* 展開・折りたたみ */
```

## Component Patterns

### ステータスバッジ（選考ステージ）
- ピル型（border-radius: full）
- 各ステージに固有の色 + アイコン
- 小さくても視認性を確保（micro サイズ + medium weight）

### パイプラインビュー（← Attio）
- カンバンボード: 選考ステージをカラムに
- ドラッグ&ドロップでステージ変更
- カードには: 企業名、現在ステージ、次の締切、最終更新日

### リストビュー（← Notion）
- テーブル表示: ソート・フィルタ可能
- インラインでステータス変更
- 行 hover で quick actions 表示

### アプリナビゲーション
- Desktop: Home / Entry / カンバン / タスク / 保存箱を左サイドバーに固定
- Mobile: Home / Entry / ボード / タスク / 保存の5タブに絞る
- Profile はアカウント確認とログアウトだけを置き、設定機能のように見せない
- Roadmap はコア導線から外し、認証後は Home に戻す

### ダッシュボード
- 最上部は次の行動を1つだけ示す
- 続けて今日のタスクを表示する
- Entry / 保存箱 / タスクへのCTAは、状況に応じて1つに絞る
- ステータス円グラフ、励ましカード、ロードマップ導線など別機能に見える要素は置かない

## CSS Variables (まとめ)

```css
:root {
  /* Colors — Light */
  --color-cream: #fbf8f1;
  --color-cream-2: #f6efd8;
  --color-surface: #ffffff;
  --color-ink: #2b2a26;
  --color-ink-2: #5c5a52;
  --color-ink-3: #8f8d82;
  --color-line: #e6e1d3;
  --color-sage: #4f6e58;
  --color-sage-soft: #e5ecde;
  --color-sage-wash: #f0f3e8;
  --color-stage-entry: #c9cbb4;
  --color-stage-doc: #a8c0da;
  --color-stage-es: #d4ba82;
  --color-stage-interview: #e9b9b0;
  --color-stage-offer: #9bb58a;

  /* Typography */
  --font-sans: 'Noto Sans JP', 'Hiragino Sans', 'Yu Gothic', system-ui, sans-serif;
  --font-serif: 'Noto Serif JP', 'Hiragino Mincho ProN', 'Yu Mincho', serif;
  --font-mono: 'JetBrains Mono', ui-monospace, 'SF Mono', Menlo, monospace;

  /* Spacing (8px base) */
  --space-2xs: 4px;
  --space-xs: 8px;
  --space-sm: 12px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;
  --space-3xl: 64px;

  /* Radius */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;

  /* Shadows */
  --shadow-soft: 0 6px 18px -8px rgba(43, 42, 38, 0.12);
  --shadow-card: 0 4px 12px -4px rgba(43, 42, 38, 0.08);
  --shadow-fab: 0 10px 24px rgba(79, 110, 88, 0.4);

  /* Transitions */
  --transition-fast: 150ms ease-out;
  --transition-base: 250ms ease-in-out;
  --transition-slow: 400ms ease-out;
}

/* Dark theme tokens are not active in the current beta UI. */
```
