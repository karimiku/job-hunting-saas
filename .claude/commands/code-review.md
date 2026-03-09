---
allowed-tools: Read, Bash, Grep, Glob, Agent, AskUserQuestion, mcp__codex__codex
argument-hint: [file-path] | [commit-hash] | --full
description: PR投稿前のコード品質レビュー（Go + Clean Architecture + DDD）
---

# Code Quality Review

対象: $ARGUMENTS

**全てのレビュー出力は日本語で記述する。**（コードスニペット・ファイルパスは英語のまま）

## 現在の状態

- Git status: !`git status --porcelain`
- Recent changes: !`git diff --stat HEAD~5`
- Repository info: !`git log --oneline -5`

## プロジェクト構成

Go + Clean Architecture + DDD のプロジェクト。
- domain層: `internal/domain/entity/`（エンティティ）、`internal/domain/value/`（値オブジェクト）、`internal/domain/repository/`（インターフェース）
- usecase層: `internal/usecase/`（ユースケースごとに1ファイル）
- handler層: `internal/handler/`（oapi-codegen ServerInterface実装、adapter層）
- infra層: `internal/infra/`（repository実装）
- middleware: `internal/middleware/`（認証等）
- 自動生成: `internal/gen/openapi/`（oapi-codegenによる生成コード）
- API定義: `api/`（OpenAPIスキーマ）

## レビュー手順

### 1. Go + Clean Architecture 品質チェック

- **依存方向**: domain→usecase→handler/infraの依存方向が守られているか。domain層が外側の層に依存していないか。
- **DIP（依存性逆転）**: usecase層がinfra実装ではなくrepositoryインターフェースに依存しているか。
- **レイヤー責務**: ビジネスロジックがhandler/infraにリークしていないか。handlerはHTTP↔UseCase変換のみか。
- **エンティティのカプセル化**: フィールドはprivate（小文字）で、getterを介してアクセスしているか。状態変更はメソッド経由か。
- **値オブジェクトの不変条件**: コンストラクタでバリデーションしているか。ゼロ値が不正な状態を表さないか。
- **Repository契約**: Save（upsert）の契約、userIDスコープ、存在判定の責務分離が明確か。

### 2. Goベストプラクティス

- **命名**: 目的命名駆動が守られているか。パッケージ名は何を提供するかを示しているか。`util`等の曖昧パッケージがないか。
- **エラーハンドリング**: `_`でのエラー握りつぶし、センチネルエラーの適切な使用、`errors.Is`/`errors.As`の活用。
- **インターフェース**: 必要なメソッドだけ定義しているか（Interface Segregation）。使う側で定義しているか。
- **並行性**: sync.Mutex/RWMutexの適切な使用、goroutineリーク、race condition。
- **コメント**: 自明な処理の復唱コメントがないか。設計意図・注意点のみコメントしているか。

### 3. セキュリティ

- 認証・認可バイパスの可能性
- SQLインジェクション（将来のDB実装に向けて）
- 環境変数・シークレットのハードコード
- 入力バリデーションの漏れ

### 4. パフォーマンス

- N+1クエリの可能性（Repository設計レベル）
- 不要なメモリアロケーション
- 適切なスライス初期化（`make`の容量指定）

### 5. テスト

- テストカバレッジの不足箇所
- エッジケースのテスト漏れ
- mockの設計が適切か

### 6. 指摘の整理

重要度別に分類:
- 【必須】: セキュリティ/バグ/設計違反（Clean Architecture依存方向違反、カプセル化破壊等）
- 【推奨】: 改善推奨（命名改善、エラーハンドリング強化等）
- 【提案】: あれば良い（コメント追加、リファクタリング案等）

ラベル: [architecture], [ddd], [naming], [interface], [error], [security], [perf], [bug], [test], [nits], [typo]

各指摘のフォーマット:
```
[ラベル] 【重要度】 タイトル
ファイル: `パス` L行番号
問題: 説明
提案: 修正案（コードスニペット付き）
```

### 7. マルチエージェント検証（最大3ラウンド）

【必須】【推奨】の指摘について、Codex MCPとサブエージェントが最大3ラウンド議論して妥当性を検証する。

#### 7a. ラウンド1: サブエージェントによる反証検証

Agent tool（subagent_type: general-purpose）で【必須】【推奨】の指摘を並列検証。**反証する立場**で検証を依頼:

```
このコードレビュー指摘の妥当性を、Go・Clean Architecture・DDDのベストプラクティスに基づいて反証する立場で検証してください:

【プロジェクト構成】Go + Clean Architecture + DDD
- domain層: entity（エンティティ）、value（値オブジェクト）、repository（インターフェース）
- usecase層: ユースケースごとに1ファイル（Create/Get/List/Update/Delete）
- handler層: oapi-codegen ServerInterface実装、adapter層
- infra層: repository実装（InMemory / PostgreSQL）

【指摘内容】{コメント全文}
【対象コード】{コードとコンテキスト}

検証結果を「✅妥当 / ⚠️要修正 / ❌false positive」で回答。
具体的な理由、Goの公式ドキュメントやEffective Goからの根拠、修正案も記載。
```

#### 7b. ラウンド1: Codex MCPによる検証

Codex MCPに指摘とサブエージェントの検証結果を渡し、シニアGoエンジニアの視点で検証:

```
以下のコードレビュー指摘と反証意見の妥当性を、シニアGoエンジニアの視点で検証してください：

【プロジェクト構成】Go + Clean Architecture + DDD
【対象ファイル】{パス}
【対象コード】{コードスニペット}
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

#### 7c. ラウンド2-3（必要な場合のみ）

ラウンド1でサブエージェントとCodexの結論が異なる場合、最大2回追加ラウンドを実施:

- サブエージェントにCodexの意見を渡して再検証を依頼
- Codexにサブエージェントの再検証結果を渡して最終判定を依頼
- 各ラウンドで相手の意見を引用し、具体的に同意/反論する

議論が収束しない場合は、安全側（指摘を残す方向）で判定する。

#### 7d. 結果統合

全ラウンドの議論を踏まえて最終判定:

- 両者合意 ✅ → 採用（高信頼）
- 議論の末に合意 → 採用（議論経緯を付記）
- 片方 ⚠️ → 指摘内容を改訂（改訂理由を付記）
- 両方 ❌ → 指摘を削除
- 収束せず → 安全側で採用し、議論の要約を付記

### 8. 最終レポート

検証済みの全指摘を重要度別に分類して報告:

- 各指摘にサブエージェント/Codexの検証結果と議論経緯を付記
- 良い点も1-3個挙げる

署名:
- Codex検証済み: `🤖 Reviewed by Claude Code (Opus 4.6) ✅ Validated by Codex & Sub-agent (max 3 rounds)`
- Codex未検証: `🤖 Reviewed by Claude Code (Opus 4.6) ✅ Validated by Sub-agent`

全て日本語で記述（コード・ファイルパスは英語のまま）
