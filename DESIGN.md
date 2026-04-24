---
brand: JobHunting SaaS
style: Warm Productivity
themes: [light, dark]
default_theme: light
references: [notion.so, attio.com, amie.so, todoist.com, linear.app]
---

# DESIGN.md

## Design Philosophy

「毎日開きたくなる就活管理ツール」
- 非エンジニアでも迷わない親しみやすさ（← Notion, Airbnb）
- パイプライン管理の情報密度と一覧性（← Attio, Linear）
- 節目だけ感動させるアニメーション（← Amie）
- 日常操作はゼロ摩擦で速く

## Colors

### Light Theme

| Token | Hex | 出典・意図 |
|-------|-----|-----------|
| background | `#F7F6F3` | Notion の warm canvas。純白より目に優しく長時間使える |
| surface | `#FFFFFF` | カード・モーダルの背景 |
| text-primary | `#37352F` | Notion の warm charcoal。黒すぎず読みやすい |
| text-secondary | `#717171` | Airbnb の muted。補足情報・ラベル |
| accent | `#5B8DEF` | 希望・前進を感じるブルー。Amie のブルーを少し柔らかく |
| accent-hover | `#4A7BE0` | accent の hover 状態 |
| success | `#34C759` | ステージ通過・タスク完了 |
| warning | `#FF9500` | 締切が近い |
| danger | `#FF3B30` | 締切超過・お見送り |
| border | `#E8E6E1` | Notion 系の暖かいボーダー |

### Dark Theme (パワーユーザー向け)

| Token | Hex | 出典・意図 |
|-------|-----|-----------|
| background | `#1A1A2E` | Linear の deep background |
| surface | `#232338` | カード背景 |
| text-primary | `#F7F7F8` | Linear |
| text-secondary | `#A1A1AA` | |
| accent | `#7B9EF7` | Light の accent を明度上げ |
| border | `#2E2E45` | |

## Typography

### フォントスタック

```css
/* UI全般 — Inter をベースに日本語フォールバック */
--font-sans: 'Inter', 'Noto Sans JP', 'Hiragino Kaku Gothic ProN',
             'Yu Gothic', sans-serif;

/* 数値・ID・ショートカット表示 */
--font-mono: 'JetBrains Mono', 'SF Mono', monospace;
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
| sm | 4px | バッジ・タグ |
| md | 8px | カード・入力フィールド |
| lg | 12px | モーダル・ドロップダウン |
| xl | 16px | 大きなカード |
| full | 9999px | ピル型ボタン・検索バー（← Airbnb） |

## Shadows

```css
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);   /* カード hover */
--shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.12);   /* モーダル */
```

## Animation Guidelines

**原則: 毎日使うところは速く、節目だけ感動させる**

### 使う場面（← Amie, Todoist）
| 場面 | 種類 | Duration | Easing |
|------|------|----------|--------|
| ステージ変更 | カード移動 | 300ms | ease-out |
| タスク完了 | チェック + フェードアウト | 250ms | ease-in-out |
| 内定獲得 | confetti / celebration | 1000ms | — |
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

### ダッシュボード
- 「今日やること」が最上部
- 締切が近いタスクのハイライト（warning カラー）
- 選考状況のサマリー（ファネル or 数値）

## CSS Variables (まとめ)

```css
:root {
  /* Colors — Light */
  --color-bg: #f7f6f3;
  --color-surface: #ffffff;
  --color-text: #37352f;
  --color-text-muted: #717171;
  --color-accent: #5b8def;
  --color-accent-hover: #4a7be0;
  --color-success: #34c759;
  --color-warning: #ff9500;
  --color-danger: #ff3b30;
  --color-border: #e8e6e1;

  /* Typography */
  --font-sans: 'Inter', 'Noto Sans JP', 'Hiragino Kaku Gothic ProN',
               'Yu Gothic', sans-serif;
  --font-mono: 'JetBrains Mono', 'SF Mono', monospace;

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
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.12);

  /* Transitions */
  --transition-fast: 150ms ease-out;
  --transition-base: 250ms ease-in-out;
  --transition-slow: 400ms ease-out;
}

[data-theme="dark"] {
  --color-bg: #1a1a2e;
  --color-surface: #232338;
  --color-text: #f7f7f8;
  --color-text-muted: #a1a1aa;
  --color-accent: #7b9ef7;
  --color-border: #2e2e45;
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.2);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.4);
}
```
