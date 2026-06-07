# PR深層コードレビュー

対象PR: #$ARGUMENTS

$ARGUMENTSが空なら「使用方法: /review-pr {PR番号}」と表示して終了。

## 1. PR情報取得

- `gh repo view --json owner,name` でowner/repoを取得
- `gh pr view $ARGUMENTS --json number,title,body,author,baseRefName,headRefName,headRefOid,state,additions,deletions,changedFiles,files` でメタデータ取得
- PRがopenでなければ終了
- `gh pr diff $ARGUMENTS` で差分取得

## 2. 既存コメント確認

- `gh api repos/{owner}/{repo}/pulls/$ARGUMENTS/comments --paginate` でインラインコメント取得
- `gh api repos/{owner}/{repo}/pulls/$ARGUMENTS/reviews --paginate` でレビュー取得
- bot（`devin-ai-integration[bot]`等）のコメントを特定・記録し、重複指摘を避ける
- botのコメントがある場合、各コメントのIDを記録する（後のステップで返信に使用）

## 3. ファイル読み取り

変更された `.go`, `.yaml`, `.yml`, `go.mod`, `go.sum`, `Dockerfile`, `docker-compose*.yml` をReadツールで全文読み取る。自動生成ファイル（`*_gen.go`）は差分のみ確認。ロックファイル・画像はスキップ。

## 4. 深層レビュー

bot既出の指摘は重複して新規投稿しない（ステップ9でbotのコメントに返信する）。以下の観点でレビュー:

### Go + Clean Architecture 観点
- **[architecture]**: Clean Architectureの依存方向違反（domain→infra等の逆依存）、レイヤー責務の混在、DIPの不備
- **[ddd]**: エンティティのカプセル化違反、値オブジェクトの不変条件不足、ドメインロジックのリーク（handler/infraにビジネスロジック）
- **[naming]**: 目的命名駆動の違反（1文字変数、意味不明な略語）、Goの命名慣習違反（パッケージ名、エクスポート名）
- **[interface]**: 不要なインターフェース、インターフェース肥大化、具象型への不要な依存
- **[error]**: エラーハンドリング不備（エラー握りつぶし、センチネルエラー未使用、errors.Is/As未活用）

### 品質観点
- **[security]**: SQLインジェクション、認証・認可バイパス、セッション管理不備、環境変数ハードコード
- **[perf]**: N+1クエリ、不要なメモリアロケーション、goroutineリーク、sync.Mutex誤用
- **[bug]**: nil参照、race condition、エッジケース、型変換の安全性
- **[test]**: テスト不足、テストの壊れやすさ、mock設計の問題
- **[nits]/[typo]**: スタイル・タイプミス・不要コメント

各指摘のフォーマット:
```
[ラベル] 【重要度】 タイトル
ファイル: `パス` L行番号
問題: 説明
提案: 修正案（コードスニペット付き）
```

重要度: 【必須】(セキュリティ/バグ/設計違反) / 【推奨】(改善推奨) / 【提案】(あれば良い)

## 5. Claude × Codex 並行レビュー＆議論

Claude と Codex が**独立して**レビューし、結果を突き合わせて**対等に議論**する。

### 5a. Codex CLI動作確認

```bash
codex exec --full-auto --sandbox read-only --cd /Users/kamiriku/my_projects/job-hunting-saas "Hello, this is a test. Please respond with 'Test successful'."
```

- ✅ 成功 → 5bに進む
- ❌ 失敗 → 最大3回リトライ。3回失敗したらAskUserQuestionで「Codex CLIが応答しません。待機しますか？」と確認。**Codex検証を自己判断でスキップすることは禁止。**

### 5b. Codex独立レビュー

Codex に**Claudeの指摘を見せずに**、PR差分を独立でレビューさせる:

```bash
codex exec --full-auto --sandbox read-only --cd /Users/kamiriku/my_projects/job-hunting-saas "あなたはシニアGoエンジニアです。以下のPR差分をコードレビューしてください。確認や質問は不要です。

【PR】#{PR番号} {タイトル}
【変更ファイル】{ファイル一覧}
【差分】
{gh pr diff の出力}

以下の観点でレビューし、指摘を出してください:
- アーキテクチャ（Clean Architecture依存方向、レイヤー責務、DIP）
- DDD（エンティティカプセル化、値オブジェクト不変条件、ドメインロジックリーク）
- セキュリティ（認証・認可、SQLインジェクション、バリデーション）
- パフォーマンス（N+1クエリ、メモリアロケーション）
- バグ・ロジックエラー（エラーハンドリング、nil参照、race condition）
- コード品質（命名、dead code、インターフェース設計）

各指摘のフォーマット:
[ラベル] 【重要度】 タイトル
ファイル: パス L行番号
問題: 説明
提案: 修正案

重要度: 【必須】(セキュリティ/バグ/設計違反) / 【推奨】(改善推奨) / 【提案】(あれば良い)"
```

### 5c. 指摘の突き合わせ

Claude指摘リストとCodex指摘リストを比較し、3カテゴリに分類:

- **合意**: 両方が同じ箇所・同じ趣旨で指摘 → そのまま採用
- **片方のみ**: 一方だけが出した指摘 → 5dで議論
- **矛盾**: 同じ箇所で異なる見解 → 5dで議論

### 5d. 議論（納得いくまで）

「片方のみ」「矛盾」の指摘について、Codex と往復で議論する。1ラウンドにつき1回のCodex呼び出し。

```bash
codex exec --full-auto --sandbox read-only --cd /Users/kamiriku/my_projects/job-hunting-saas "以下のコードレビューについて議論しましょう。確認や質問は不要です。具体的に回答してください。

【対象ファイル】{パス}
【変更内容】{差分}

【Claudeの見解】
{Claudeの指摘または「この箇所は問題なしと判断」}

【Codexの見解】
{Codexの指摘または「この箇所は問題なしと判断」}

あなたの立場を維持するか、相手の見解に同意するか、判断してください。
- 自分の見解を維持する場合: 具体的な根拠を示してください
- 相手に同意する場合: なぜ同意するか説明してください
- 新しい代替案がある場合: 提示してください

回答フォーマット:
【判定】維持 / 同意 / 代替案あり
【理由】...
【最終的な指摘案】（同意または代替案の場合）..."
```

議論の終了条件:
- **合意に達した** → 合意内容を最終指摘として採用
- **3ラウンド経過しても合意しない** → 両方の見解を併記してユーザーに判断を委ねる

### 5e. サブエージェント最終チェック

議論で確定した【必須】【推奨】の指摘について、Agent tool（subagent_type: general-purpose）で最終チェック。第三者の視点で妥当性を確認:

```
Claude と Codex が議論して合意した以下のPRレビュー指摘を、第三者の視点で最終チェックしてください:
- 指摘: {合意した指摘内容}
- 対象コード: {差分とコンテキスト}
- 議論経緯: {簡潔な議論サマリー}
見落としや論理的な穴がないか確認し、「✅承認 / ⚠️懸念あり」で回答。懸念がある場合は具体的に指摘。
```

## 6. ユーザーに報告

全指摘をユーザーに一覧提示する。重要度別に分類し、各指摘のファイル・行番号・問題・提案を表示。各指摘に「合意」「議論の末採用」「ユーザー判断待ち」のステータスを付記。Codex独自の指摘で採用されたものも明記。良い点も1-3個挙げる。

## 7. 修正許可の確認

AskUserQuestionで確認。選択肢: 「全て修正」「【必須】のみ修正」「修正しない（コメント投稿のみ）」

## 8. コード修正・コミット

承認された指摘をEditツールで修正。修正は最小限に。完了後:
- `git add {修正ファイル}` で個別ステージング
- コミットメッセージ: `fix: PR#$ARGUMENTS code review指摘事項の修正`（Co-Authored-By付き）
- pushはしない

## 9. GitHubレビュー投稿

`gh pr view $ARGUMENTS --json headRefOid --jq .headRefOid` で最新SHA取得後、PRレビューとして一括投稿:

```
gh api repos/{owner}/{repo}/pulls/$ARGUMENTS/reviews --method POST \
  -f commit_id="SHA" -f body="サマリー" -f event="COMMENT" \
  --jsonarray -f comments='[{"path":"ファイル","line":行番号,"side":"RIGHT","body":"コメント"}]'
```

- lineはdiffの行番号（ファイルの絶対行番号ではない）
- 修正済みの指摘には `✅ 修正済み（コミットSHAで修正）` を追記
- コメント上限: 【必須】無制限 / 【推奨】最大5 / 【提案】最大3 / 合計15
- **botコメントへの返信（必須）**: ステップ2で記録したbotの各コメントに対して、`in_reply_to`で必ず返信する:
  - 同意する場合: 「✅ 同意します。{補足理由}」+ 修正済みなら「✅ 修正済み（コミットSHAで修正）」
  - 部分同意: 「⚠️ 部分的に同意します。{理由と代替案}」
  - 不同意: 「❌ この指摘は不要と判断しました。{理由}」
  - 返信のAPIコール例:
    ```
    gh api repos/{owner}/{repo}/pulls/$ARGUMENTS/comments \
      --method POST \
      -f body="返信内容" \
      -F in_reply_to={botのコメントID}
    ```
- サマリーに指摘数・修正状況・良い点を含める
- 署名: `🤖 Reviewed by Claude Code (Opus 4.6) × Codex CLI — 独立レビュー＆議論済み`
- 指摘0件なら `LGTM! 🎉` のサマリーのみ投稿
- 全て日本語で記述（コード・ファイルパスは英語のまま）
