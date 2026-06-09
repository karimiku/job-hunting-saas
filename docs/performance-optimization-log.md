# Performance Optimization Log

Status: active for beta
Created: 2026-06-09

## Purpose

公開ベータ前後で確認したSSR / RSCのパフォーマンス問題と、実施した改善を記録する。

このログは、単なる結果メモではなく、どの構造が遅さを生んでいたか、どの改善がどれくらい効いたかを後から説明できる状態にするためのもの。

## Context

FrontendはVercel上のNext.js App Routerで動作している。主要画面はRSC / SSRでHTMLまたはRSC payloadをサーバ側生成し、その過程でGo APIを呼び出す。

そのため、ブラウザのHARで見える `https://entre.kamiriku.com/*?_rsc=...` の遅さは、主に次の合計になる。

```text
Browser
  -> Vercel Next Server
    -> Cloud Run Go API
      -> Supabase PostgreSQL
    -> RSC payload生成
  -> Browserへ返却
```

HARには `api.entre.kamiriku.com` への直接リクエストは出ない。API呼び出しはNext Server側で行われる。

## Initial Problem

最初に大きく遅かったのは `/task` だった。

初期HARの代表値:

```text
/task       max 1425ms
/inbox      max 607ms
/entry      max 585ms
/kanban     max 373ms
/dashboard  max 379ms
```

主因は、Task画面のSSRでEntryごとにTaskを取得していたこと。

```text
Task page RSC
  -> entries一覧を取得
  -> entriesをfor each
      -> /api/v1/entries/{entryId}/tasks
  -> 全タスクを集約
```

これは典型的なN+1。Entry数が増えるほど、Next ServerからGo APIへのserver-side fetchが増える構造だった。

## Fix 1: 全タスク取得API

対応:

```text
GET /api/v1/tasks
```

を追加し、ログインユーザーの全Taskを1回で取得するようにした。

変更後の構造:

```text
Task page RSC
  -> entries一覧を取得
  -> /api/v1/tasks を1回取得
  -> entryIdで会社名をjoin
```

この時点で `/task` の最大値は大きく改善した。

```text
/task max 1425ms -> 約594ms
```

## Fix 2: Task RSCの逐次await削減

全タスクAPI追加後も、Task画面では `entries` の取得完了を待ってから `tasks` を取得する逐次awaitが残っていた。

改善前:

```text
entries = await listEntriesWithCompanyNamesServer()
tasks, clips = await Promise.all([
  listAllTasksServer(entries),
  listInboxClipsServer()
])
```

改善後:

```text
entries, tasks, clips = await Promise.all([
  listEntriesWithCompanyNamesServer(),
  listTasksServer(),
  listInboxClipsServer()
])

tasksWithCompanyName = attachCompanyNamesToTasks(tasks, entries)
```

取得とjoinを分離し、独立したserver-side fetchを並列化した。

## Fix 3: navCounts用の重複fetch削減

次に残っていた問題は、各ページで `getNavCountsServer()` が `entries / tasks / inbox` を再取得していたこと。

例:

```text
Entry page RSC
  -> listEntriesWithCompanyNamesServer()
  -> getNavCountsServer()
       -> listEntriesServer()
       -> listTasksServer()
       -> listInboxClipsServer()
```

Entry画面では表示用にすでにentriesを取得しているのに、navCounts用にentriesを再取得していた。

対応として `buildNavCounts()` を追加し、ページ内で取得済みのデータからサイドバー件数を組み立てるようにした。

```text
entries, tasks, clips = await Promise.all([...])
navCounts = buildNavCounts(entries, tasks, clips)
```

この対応でHAR上のRSCリクエスト数は減少した。

```text
RSC/fetch entries: 22 -> 16
```

## Latest HAR

navCounts重複fetch削減後のHAR代表値:

```text
/entry      max 547ms avg 505ms
/inbox      max 542ms avg 445ms
/task       max 491ms avg 352ms
/kanban     max 401ms avg 278ms
/dashboard  max 179ms avg 179ms
```

最初の `/task` は以下まで改善した。

```text
/task max 1425ms -> 491ms
```

ただし、全体としてまだ十分に速いとは言い切れない。

## Current Bottleneck

Frontend側の明確なN+1と重複fetchは削減済み。

次に残っている構造的な改善余地は、サイドバー件数を出すためにfull listを取得している点。

現在:

```text
navCountsを作るために
  entries full list
  tasks full list
  inbox clips full list
を取得している
```

理想:

```text
GET /api/v1/nav-counts
=> {
  entry: 10,
  task: 3,
  inbox: 2
}
```

Backend側で `COUNT(*)` に寄せることで、Next Server側のpayload処理とDB/Network転送量をさらに減らせる。

## Design Note

今回の改善では、Clean Architecture / DDDの境界は崩していない。

- RSC pageはserver-side data compositionを担当
- API contractはOpenAPIから追加
- UseCase / RepositoryはTask一覧取得の責務を追加
- SQLは必要な集約やjoinをDB側に寄せる
- UI用の会社名joinはFrontend Server側で取得済みデータを突き合わせる

今後 `GET /api/v1/nav-counts` を追加する場合も、UseCaseとして「sidebar summary」を切り出し、SQL側ではfull listではなくCOUNT系クエリに寄せる。

