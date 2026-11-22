---
name: address-pr
description: PR のレビューコメント(未解決スレッド)を読み、分類 → 方針提示(承認待ち) → 修正・コミット → 検証 → push → 返信・resolve まで行う。引数は PR 番号。「レビュー対応して」「PR のコメント見て」と言われたら使う。
metadata:
  deps.requires: "commit-strategy,pr-review,adr-writer"
  deps.related: "dev-cycle"
---

# PR レビュー対応

作業系の指摘はさっさと直し、設計系の指摘は止まって合意を取る。そのために **フェーズ 3 で必ず停止する**。

引数: PR 番号(省略時は `gh pr view --json number -q .number` で現在ブランチの PR を使う)

## dev-cycle モード

`dev-cycle` Skill の resume(工程 7・8)から呼ばれた場合は、人間の PR レビューそのものを「判断」として扱い、**フェーズ 3 で停止しない**。

- フェーズ 0〜2 は同じ
- フェーズ 3 の方針提示と承認待ちを省き、`dev-cycle` の工程 7 の表(人間が案を選んだ design は反映、判断が読み取れない design は選択肢を返信して未解決のまま残す)に従う
- 人間の設計判断は、実装の前に「設計判断の記録(ADR)」に従って ADR に記録する(`dev-cycle` の工程 7-3)
- question への回答は承認を待たずに返信する
- 修正は開発 Agent(`go-dev` / `web-dev`)の修正モードで行い、確認は新しく起動した `reviewer` の再レビューモードで行う(`dev-cycle` の「reviewer の起動ルール」)。フェーズ 4・5 はそれで代える
- フェーズ 6 は同じ(未達の受け入れ条件に対応するスレッドは resolve しない)

## Required Context

```text
CLAUDE.md                                # 存在する場合
.claude/rules/git-workflow.md            # ブランチ運用・検証コマンド
docs/spec/requirements.md                     # 要件 ID・検証ルール
docs/spec/design.md                           # 設計、決定済み事項の索引(12 章)
docs/adr/                                # 設計判断の詳細(承認済み・提案中の ADR)
.claude/skills/adr-writer/SKILL.md       # 設計判断の記録方法
.claude/skills/commit-strategy/SKILL.md
.claude/skills/pr-review/SKILL.md        # レビューコメントの優先度(Blocking / Non-blocking / Question)
.claude/rules/                           # 修正対象のファイルに対応するルール
```

## フェーズ 0: 前提チェック(失敗したら止まって報告)

- `gh auth status` が成功すること
- `gh pr view <N> --json headRefName,baseRefName,state` で、PR が OPEN であること
- `headRefName` と `git branch --show-current` が一致すること
- 現在のブランチが `main` / `develop` でないこと
- `git status --porcelain` が空であること(未コミットの変更があれば stash せず報告)
- `git fetch` 後、ローカルがリモートより遅れていないこと(遅れていれば報告)

## フェーズ 1: 収集

- `bash .claude/skills/address-pr/scripts/fetch-threads.sh <N>` で未解決スレッドを JSON で取得する。
  各スレッドは `id` / `path` / `line` / `outdated` / `comments[]` を持つ。
- PR 本文への会話コメント(スレッドに紐づかないもの)とレビュー本文は、次のコマンドで読む。

  ```bash
  gh pr view <N> --json comments,reviews \
    -q '{comments: [.comments[] | {author: .author.login, body}], reviews: [.reviews[] | {author: .author.login, state, body}]}'
  ```

  `gh pr view <N> --comments` は使わない(gh 2.45 では廃止された Projects (classic) を参照してエラーになる)。
- PR 本文の「人間に判断してほしい点」(`D1`…)と、PR に含まれる `提案` の ADR(`docs/adr/` の差分)を一覧にする。レビュー本文・会話コメント・スレッドのどこかでそれに答えていないかを探す(「D1: 案 A」「ADR-0005 は却下」など)。
- 指摘されたコードを実際に読み、指摘が現在のコードに当てはまるかを確認する。

## フェーズ 2: 分類

各スレッドと、スレッドに紐づかない判断(レビュー本文・会話コメントで「人間に判断してほしい点」や `提案` の ADR に答えたもの)を、次の 4 つに分ける。**迷ったら design に倒す。**

| 分類 | 基準 | 扱い |
|---|---|---|
| fix | 命名、typo、テスト・`testdata/` の追加、エラー処理の漏れ、lint、既存パターンへの統一など、設計を変えない修正 | 承認後に実装 |
| design | 責務分割、型・インターフェースの形、レイヤー間の依存の向き、エラー方針、「なぜこうした?」を含むもの | 実装せず選択肢を提示 |
| question | 説明を求めているだけで変更要求がない | 返信文案のみ作成 |
| skip | 既に別コミットで対応済み、または指摘箇所が消えている | 理由を添えて報告 |

このプロジェクト固有の基準:

- 次のいずれかを変更する指摘は、内容が小さくても **design** にする。
  - `docs/spec/requirements.md` / `docs/spec/design.md` の内容(要件、検証ルールの ID・判定条件・既定の重大度)
  - design.md 12 章・承認済みの ADR の決定済み事項(探索ルート、`deps.*` の書式、設定ファイルの場所など)
  - `schemas/index.schema.json` とデータモデル(JSON 出力の契約)
  - CLI のコマンド・オプション・終了コード
  - レイヤー間の依存の向き、診断を作る場所(Rules 層だけ)
  - 安全性(NFR-04:書き込み系の手順、GUI サーバーのセキュリティ)
- コメントの先頭の `[Blocking]` / `[Non-blocking]` / `[Question]`(`pr-review` Skill)は優先度であって分類ではない。`[Question]` は原則 question、それ以外は内容で fix / design を判断する。
- 指摘が要件や設計と矛盾する場合は、指摘に従わず design として矛盾を提示する。

## 設計判断の記録(ADR)

人間がレビューで設計判断を下したら、実装の前に `adr-writer` Skill で ADR に記録する。記録はメインのセッションが Skill ツールで行う(サブエージェントには任せない)。

**ADR を書く対象**(1 件の判断につき ADR 1 件)

- design に分類され、人間が案を選んだ指摘
- 「人間に判断してほしい点」(`D1`…)への回答
- `提案` の ADR への回答(採用・却下・別案)
- 決定済み事項(design.md 12 章・承認済みの ADR)を変える指摘(fix に見えても対象)

状態の決め方と記録の手順は `adr-writer` Skill の「人間の判断と ADR の状態」「手順」に従う。ADR は `docs: ADR-NNNN <タイトル>` として実装とは別にコミットし、判断に対応するスレッド・コメントへの返信に ADR へのリンクを入れる(フェーズ 6)。

## フェーズ 3: 方針提示 ★必ず停止

`templates/plan.md` の形式で出力し、ユーザーの承認を待つ。

- fix: 変更内容を 1 行で
- design: 選択肢を 2〜3 個、トレードオフ、推奨とその理由。**コードは書かない**。ユーザーが選んだら「設計判断の記録(ADR)」に従って ADR を書いてから実装する
- question: 返信文案
- skip: 理由

ユーザーが「進めて」と言うまでフェーズ 4 に進まない。
design の項目は、ユーザーが選択肢を選ぶまで実装しない。選ばれなかった design は未対応として残す。

## フェーズ 4: 実装

`commit-strategy` Skill のルールに従う。

- fix は原則 1 スレッド = 1 コミット。関連する指摘はまとめてよい。
- design はユーザーの決定どおりに実装し、コミット本文に決定内容と理由を書く。決定が設計判断にあたる場合は `adr-writer` Skill で ADR に記録する。
- コミットメッセージは `<type>: <要約>`(type は feat / fix / refactor / test / docs / chore)。本文に PR 番号と対応したスレッド ID をすべて書き、関係する要件 ID があれば `Refs:` に書く。

  ```text
  fix: 大文字小文字だけが違うリンクを E007 にする

  PR #42 のレビュー指摘への対応。
  Thread: PRRT_xxxxxxxx
  Refs: E007
  ```

- 生成物は手で編集しない。
- 承認された範囲以外の変更(ついでのリファクタリングなど)をしない。

## フェーズ 5: 検証

`.claude/rules/git-workflow.md` の「検証コマンド」（build・vet・golangci-lint・test）をすべて実行する。

失敗したら push せずに報告する。実行できなかった場合は、成功と報告せず理由を書く。

## フェーズ 6: push・返信・resolve

1. `git fetch` し、`git merge-tree --write-tree origin/<base> HEAD` でマージ先(`baseRefName`)とコンフリクトしないことを確認してから `git push`(PR のブランチへ。force push しない)。コンフリクトがあればマージ先を取り込んで解消し、フェーズ 5 をやり直す
2. `gh pr view <N> --json mergeable,mergeStateStatus` と(CI があれば)`gh pr checks <N>` で、マージ可能な状態に戻ったことを確認する
3. 各スレッドに返信する。

   ```bash
   bash .claude/skills/address-pr/scripts/reply-resolve.sh <thread_id> "<本文>" [true|false]
   ```

   | 分類 | 本文 | resolve |
   |---|---|---|
   | fix | 「対応しました: <commit sha>」と変更の要約 | true |
   | design(決定済み) | 「<選んだ案> で対応しました: <commit sha>。記録: ADR-NNNN(<リンク>)。理由: ...」 | true |
   | design(未決定) | 返信しない | - |
   | question | 承認された回答文 | **false**(質問者が resolve する) |
   | skip | 対応済みのコミット、または指摘箇所が消えた理由 | true |

4. 最後に対応サマリを出す。
   - スレッド数(分類ごと)
   - コミット一覧(sha とメッセージ)
   - build / vet / golangci-lint / test の結果、マージ可能かどうか
   - 未対応のスレッドとその理由

design で決めた内容が要件・設計・JSON スキーマに関わる場合は、`docs/` / `schemas/` の更新も同じ PR に含める(フェーズ 3 の方針に含めて承認を取る)。

## やってはいけないこと

- フェーズ 3 の承認なしにコードを変更する(dev-cycle モードを除く)
- design に分類した指摘を独断で実装する
- force push、PR のブランチ以外への push、`main` / `develop` への直接コミット
- 返信せずにスレッドを resolve する
- question のスレッドを resolve する
- build / vet / golangci-lint / test が失敗した状態で push する
