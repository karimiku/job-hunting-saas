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

bot既出の指摘は重複して新規投稿しない（ステップ10でbotのコメントに返信する）。以下の観点でレビュー:

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

## 5. マルチエージェント検証（最大3ラウンド）

【必須】【推奨】の指摘について、Codex MCPとサブエージェントが最大3ラウンド議論して妥当性を検証する。

### 5a. ラウンド1: サブエージェントによる反証検証

Agent tool（subagent_type: general-purpose）で【必須】【推奨】の指摘を並列検証。**反証する立場**で検証を依頼:

```
このPRレビュー指摘の妥当性を、Go・Clean Architecture・DDDのベストプラクティスに基づいて反証する立場で検証してください:

【プロジェクト構成】Go + Clean Architecture + DDD
- domain層: entity（エンティティ）、value（値オブジェクト）、repository（インターフェース）
- usecase層: ユースケースごとに1ファイル（Create/Get/List/Update/Delete）
- handler層: oapi-codegen ServerInterface実装、adapter層
- infra層: repository実装（InMemory / PostgreSQL）

【指摘内容】{コメント全文}
【対象コード】{差分とコンテキスト}

検証結果を「✅妥当 / ⚠️要修正 / ❌false positive」で回答。
具体的な理由、Goの公式ドキュメントやEffective Goからの根拠、修正案も記載。
```

### 5b. ラウンド1: Codex MCPによる検証

Codex MCPに指摘とサブエージェントの検証結果を渡し、シニアGoエンジニアの視点で検証:

```
以下のコードレビュー指摘と反証意見の妥当性を、シニアGoエンジニアの視点で検証してください：

【プロジェクト構成】Go + Clean Architecture + DDD
【対象ファイル】{パス}
【変更内容】{差分}
【レビューコメント案】{指摘内容}
【サブエージェントの反証意見】{サブエージェントの検証結果}

【確認観点】
1. 指摘内容はGoのベストプラクティス（Effective Go, Go Proverbs）に照らして正確か
2. Clean Architecture / DDDの原則に沿っているか
3. サブエージェントの反証は妥当か
4. 重要度の判定は適切か
5. より良い代替案はないか

検証結果を「✅妥当 / ⚠️要修正 / ❌false positive」で回答。理由と修正案も記載。
サブエージェントの意見に同意/反対する場合はその理由も明記。
```

### 5c. ラウンド2-3（必要な場合のみ）

ラウンド1でサブエージェントとCodexの結論が異なる場合、最大2回追加ラウンドを実施:

- サブエージェントにCodexの意見を渡して再検証を依頼
- Codexにサブエージェントの再検証結果を渡して最終判定を依頼
- 各ラウンドで相手の意見を引用し、具体的に同意/反論する

議論が収束しない場合は、安全側（指摘を残す方向）で判定する。

### 5d. 結果統合

全ラウンドの議論を踏まえて最終判定:

- 両者合意 ✅ → 採用（高信頼）
- 議論の末に合意 → 採用（議論経緯を付記）
- 片方 ⚠️ → 指摘内容を改訂（改訂理由を付記）
- 両方 ❌ → 指摘を削除
- 収束せず → 安全側で採用し、議論の要約を付記

## 6. ユーザーに報告

検証済みの全指摘をユーザーに一覧提示する。重要度別に分類し、各指摘のファイル・行番号・問題・提案を表示。各指摘にサブエージェント/Codexの検証結果と議論経緯を付記。良い点も1-3個挙げる。

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
- 署名（Codex検証済み）: `🤖 Reviewed by Claude Code (Opus 4.6) ✅ Validated by Codex & Sub-agent (max 3 rounds)`
- 署名（Codex未検証）: `🤖 Reviewed by Claude Code (Opus 4.6) ✅ Validated by Sub-agent`
- 指摘0件なら `LGTM! 🎉` のサマリーのみ投稿
- 全て日本語で記述（コード・ファイルパスは英語のまま）
