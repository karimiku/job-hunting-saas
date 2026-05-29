---
allowed-tools: Read, Bash, Grep, Glob, Agent, AskUserQuestion
argument-hint: [file-path] | [commit-hash] | --full
description: セキュリティ・パフォーマンス・アーキテクチャ観点の包括的コードレビュー（Claude × Codex 議論型）
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

### 7. Claude × Codex 並行レビュー＆議論

Claude と Codex が**独立して**レビューし、結果を突き合わせて**対等に議論**する。

#### 7a. Codex CLI動作確認

```bash
codex exec --full-auto --sandbox read-only --cd /Users/kamiriku/my_projects/job-hunting-saas "Hello, this is a test. Please respond with 'Test successful'."
```

- ✅ 成功 → 7bに進む
- ❌ 失敗 → 最大3回リトライ。3回失敗したらAskUserQuestionで「Codex CLIが応答しません。待機しますか？」と確認。**Codex検証を自己判断でスキップすることは禁止。**

#### 7b. Codex独立レビュー

Codex に**Claudeの指摘を見せずに**独立でレビューさせる:

```bash
codex exec --full-auto --sandbox read-only --cd /Users/kamiriku/my_projects/job-hunting-saas "あなたはシニアGoエンジニアです。以下のファイルをコードレビューしてください。確認や質問は不要です。

【対象ファイル】{レビュー対象ファイルのパス一覧}

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

#### 7c. 指摘の突き合わせ

Claude指摘リストとCodex指摘リストを比較し、3カテゴリに分類:

- **合意**: 両方が同じ箇所・同じ趣旨で指摘 → そのまま採用
- **片方のみ**: 一方だけが出した指摘 → 7dで議論
- **矛盾**: 同じ箇所で異なる見解 → 7dで議論

#### 7d. 議論（納得いくまで）

「片方のみ」「矛盾」の指摘について、Codex と往復で議論する。1ラウンドにつき1回のCodex呼び出し。

```bash
codex exec --full-auto --sandbox read-only --cd /Users/kamiriku/my_projects/job-hunting-saas "以下のコードレビューについて議論しましょう。確認や質問は不要です。具体的に回答してください。

【対象ファイル】{パス}
【対象コード】{コードスニペット}

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

#### 7e. サブエージェント最終チェック

議論で確定した【必須】【推奨】の指摘について、Agent tool（subagent_type: general-purpose）で最終チェック。第三者の視点で妥当性を確認:

```
Claude と Codex が議論して合意した以下のコードレビュー指摘を、第三者の視点で最終チェックしてください:
- 指摘: {合意した指摘内容}
- 対象コード: {コードとコンテキスト}
- 議論経緯: {簡潔な議論サマリー}
見落としや論理的な穴がないか確認し、「✅承認 / ⚠️懸念あり」で回答。懸念がある場合は具体的に指摘。
```

### 8. 最終レポート

全指摘を重要度別に分類して報告:

- 各指摘に「合意」「議論の末採用」「ユーザー判断待ち」のステータスを付記
- Codex独自の指摘で採用されたものも明記
- 良い点も1-3個挙げる

署名: `🤖 Reviewed by Claude Code (Opus 4.6) × Codex CLI — 独立レビュー＆議論済み`

全て日本語で記述（コード・ファイルパスは英語のまま）
