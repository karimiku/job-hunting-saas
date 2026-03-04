# 技術選定 Why Reasons

本ドキュメントは、各技術選定の「なぜ？」をこのプロジェクト固有の理由で整理したものである。

---

## 1. Go

### なぜGo？

- Chrome拡張とWeb UIの2つの入口からAPIを叩く構成。拡張で「1クリック保存」を実現するにはAPIレスポンスが速くないとUXが死ぬ。Goのコンパイル済みバイナリは起動・応答が速い
- source推定ロジック（サイト別にDOM解析パターンが違う）をinterfaceで切り替える設計が必要。Goの暗黙的インターフェース実装なら、サイト別の推定ロジックを追加するとき既存コードを変更せずに済む（開放閉鎖原則）
- 単一バイナリ → Dockerイメージが数十MB。コンテナデプロイが速く、コスト効率もいい
- 継承がなく、interfaceとcompositionでOOP原則を自分で設計する必要がある。フレームワークに隠蔽されず設計力がコードに出る

### Javaじゃダメ？

- Springが多くを隠蔽する。このプロジェクトではDIP・DI・カプセル化を自分の設計判断として実装したい。Goならフレームワークの魔法なしに全部自分で組み立てる

---

## 2. Clean Architecture

### なぜClean Architecture？

- このプロジェクトにはAPIの入口が2つある。Chrome拡張とWeb UI。両方が同じUseCaseを叩く。レイヤーが分かれてないと、拡張用とWeb用でビジネスロジックが重複する

```
Chrome拡張 → Handler A ─┐
                         ├→ UseCase（共通）→ Domain
Web UI     → Handler B ─┘
```

- DB・インフラがまだ未確定の状態でドメイン層から実装を始められる。Entry・Task・Eventの状態遷移ルールや値オブジェクト（EntryStatus, Source等）をDBなしで作ってテストできる
- 通知メール機能が後から入る。メール送信の詳細がUseCase層に漏れてると、通知手段を変えるたびにビジネスロジックを触ることになる。interfaceで分離しておけば通知手段の変更はインフラ層だけで済む

### オーバーエンジニアリングでは？

- CRUDが多いのは事実。ただしInboxの割当推定、source推定、企業名マッチングなど、推定ロジックが複数ある。これらをインフラ層やHandler層に書くと散らばる。UseCase層・Domain層に閉じ込めるためにレイヤー分離が必要

---

## 3. PostgreSQL

### なぜPostgreSQL？

- EntryStatusは5種類（in_progress/offered/accepted/rejected/withdrawn）で固定。PostgreSQLのENUM型を使えば、DBレベルで不正な状態を弾ける。値オブジェクトでの型安全とDB側でも型安全の二重ガード
- ダッシュボードは「進行中のEntryだけ」を表示する画面。PostgreSQLの部分インデックスなら `WHERE status IN ('in_progress', 'offered')` のインデックスを作れる。年間95万Entryの中から進行中だけを高速に引ける。MySQLにはこの機能がない
- Inboxアイテムはchrome拡張から送られる不定形データ（推定成功度がバラバラ）。PostgreSQLのJSONBを使えば、推定できた項目だけ柔軟に保存できる
- 退会時の全データ削除が要件にある。PostgreSQLのCASCADE DELETEとトランザクションの堅牢さで、User削除時に関連データを確実に消せる

### なぜMySQLじゃない？

- 部分インデックスがない → ダッシュボードの「進行中だけ」高速検索に不利
- ENUM型のALTERが面倒 → 状態追加時の運用コストが高い
- JSONB相当の機能で検索性能が劣る → Inboxの柔軟なデータ保存に不利

---

## 4. sqlc（SQLコード生成）

### なぜsqlc？

- SQLを書いて `sqlc generate` を実行すると、Go の関数と構造体が自動生成される。スキーマ定義（schema.sql）とクエリ（query.sql）をコード生成時に照合するため、カラム名のタイポや型の不一致が**コード生成時点**で検出される。アプリを動かす前に弾ける
- PostgreSQL 自体も型安全だが、それは SQL がDBに届いた「実行時」の話。sqlc はその2段階前（コード生成時 → コンパイル時）に型チェックを持ってくる

```
schema.sql（テーブル定義）
    ↓
query.sql（SQLクエリ）
    ↓
sqlc generate  ← ここでカラム名・型の不一致を検出
    ↓
Go コード生成（関数 + 構造体）
    ↓
go build  ← 生成コードを使う側も型チェック
    ↓
実行時（PostgreSQL）← ここまで来る前に弾ける
```

- 生成コードはインフラ層に閉じ込められる。ドメイン層の Entry 構造体に DB のタグが混入しない。ビジネスルールと DB の都合が完全に分離される
- openapi.yaml → oapi-codegen、schema.sql → sqlc。外側を自動生成、中心を手動という設計思想が全レイヤーで一貫する

### なぜGORMじゃない？

- GORM はリフレクションで構造体を読むため、フィールドを公開（大文字）にしないと動かない。このプロジェクトのエンティティは小文字フィールド + getter でカプセル化しており、GORM を使うとこの設計が崩壊する
- `gorm.Model` や `gorm:"column:..."` タグがドメイン構造体に混入し、ドメイン層が DB のカラム名やテーブル構造を知ってしまう。Clean Architecture / DDD の「ドメイン層は外部に依存しない」原則に違反する
- 値オブジェクト（`value.CompanyName` 等）を GORM が理解できないため、`string` に戻すか変換処理が必要になる。型安全なドメインモデルが無意味になる
- SQL が隠蔽されて何が発行されてるか分からない。ダッシュボードのような複合クエリのチューニングが困難

### なぜsqlxじゃない？

- sqlx は SQL を文字列で書いて構造体タグで自動マッピングする。カラム名を間違えてもコンパイルは通り、実行して初めてエラーが分かる
- sqlc なら同じ間違いがコード生成時に検出される。年間95万Entryを扱うSaaSで、型の不一致が本番で発覚するリスクは取りたくない

### なぜ database/sql（標準ライブラリ）じゃない？

- 結果を1カラムずつ `row.Scan(&c.ID, &c.Name, ...)` で手動マッピングする必要がある。カラム数が多いとボイラープレートが膨大で、順番を間違えても実行時まで気づけない
- sqlc はこのマッピングも自動生成するので、手動の Scan ミスが構造的に起きない

---

## 5. Chi（Goルーター）

### なぜChi？

- oapi-codegenが生成するのは `net/http` 互換のインターフェース。Chiは `net/http` 互換だからアダプタなしで直結する
- このプロジェクトは認証ミドルウェアが必須（Chrome拡張とWebの両方でユーザー特定が必要）。標準ライブラリだけだとミドルウェアチェーンが冗長。Chiならミドルウェアの仕組みが揃っていて、かつ `net/http` から逸脱しない

### なぜGin/Echoじゃない？

- Gin/Echoは独自のContext型を使う。Handler層がフレームワークに密結合し、oapi-codegenの生成コードとの間にアダプタが必要になる
- Clean ArchitectureのDIP（依存性逆転）に反する

### なぜ標準ライブラリだけじゃない？

- Go 1.22で改善されたが、認証ミドルウェアやルートグループ化が冗長。このプロジェクトはChrome拡張とWebで認証ミドルウェアの適用パターンが異なるため、Chiのルートグループ化が実用的に必要

---

## 6. OpenAPI + oapi-codegen

### なぜスキーマ駆動？

- Chrome拡張とWeb UIが同じAPIを叩く。API契約が曖昧だと拡張とWebで期待するレスポンスがズレる。OpenAPIで契約を先に固めることで、拡張開発とバックエンド開発を並行できる
- Inboxの保存API（FR-70: 失敗しないUX）は設計が重要。「最低限URL+タイトル+時刻で必ず保存」というフォールバック仕様をOpenAPIに明記しておけば、拡張側もバックエンド側も同じ契約で動ける

### なぜoapi-codegen？

- Chi（net/http）互換のコードを生成する。Gin/Echo専用ジェネレータではない
- 生成されるのはインターフェースとDTO → Handler層に閉じ込められる。ドメイン層には一切影響しない

---

## 7. TDD

### なぜTDD？

- Entryの状態遷移にルールがある。in_progress → offered はOKだが、accepted → in_progress に戻していいのか？こういうビジネスルールはテストで先に振る舞いを確定させてから実装したほうが正確
- Inboxの割当推定は推定ロジック。「この証跡はどのEntryに紐づくか」の候補提示に正解/不正解がある。テストケースを先に書いて期待値を定義することで、推定精度を定量的に管理できる
- Repository interfaceがあるからTDDが成立する。UseCaseのテストでインメモリRepositoryを注入すれば、DBなしで数ミリ秒でテストが回る。Clean ArchitectureがTDDを可能にし、TDDがClean Architectureの正しさを検証する

---

## 8. ローカル開発環境: Docker Compose（Go + PostgreSQL）

### なぜ全部Dockerに入れる？

- 本番はコンテナデプロイ。ローカルだけGoをホストで動かすと、本番と開発で環境が乖離する。Dockerfileを共有することで「ローカルで動く = 本番でも動く」を担保できる
- Docker Composeで Go + PostgreSQL をまとめて起動すれば、新メンバーが `docker compose up` だけで開発環境が立ち上がる。Go のバージョンや依存ツールの差異が出ない
- PostgreSQLのバージョンを本番と揃えるため。`postgres:16` のイメージで統一すれば環境差が出ない

---

## 設計思想の一貫性

### 自動生成戦略

```
外側（自動生成で型安全を担保）
├── openapi.yaml → oapi-codegen → Handler層
├── schema.sql   → sqlc         → Repository層

中心（手動 + TDDで品質を担保）
└── Domain層 + UseCase層
```

- 外側はスキーマから自動生成 → 型の不一致が構造的に起きない
- 中心はビジネスロジック → 人間が考えて書くべき領域
- 自動生成と手動実装の境界がClean Architectureのレイヤー境界と一致している

### 開発フロー

```
1. openapi.yaml を定義（API契約の確定）
2. oapi-codegen で Handler Interface + DTO を生成
3. Domain層の値オブジェクト・エンティティを作る
4. UseCase層のinterfaceとテストを先に書く（Red）
5. UseCase・Domainの中身を実装する（Green）
6. 自動生成されたHandlerとUseCaseを繋ぐ
7. schema.sql → sqlc でRepository実装を生成・ラップ
8. main.go でDI配線
```
