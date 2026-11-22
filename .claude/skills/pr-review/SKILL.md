---
name: pr-review
description: PR やブランチの差分を、musubi の要件・設計・.claude/rules に照らしてレビューする。「レビューして」「PR を見て」と言われたとき、または PR 作成前のセルフレビューで使う。コードは変更しない。
---

# PR Review Skill

## Purpose

Pull Request の変更内容をレビューし、プロジェクトのルール・設計・品質基準に対する問題を検出する。

## Required Context

レビュー開始前に、以下を確認する。

```text
CLAUDE.md                     # 存在する場合
.claude/rules/                # 変更対象のパスに対応するルール（project.md は常に）
docs/spec/requirements.md          # 機能要件 F-xx、検証ルール E/W/Ixxx、非機能要件 NFR-xx
docs/spec/design.md                # レイヤー構成、CLI、データモデル、設計判断の索引（12 章）
docs/adr/                     # 設計判断の詳細（承認済みの ADR）
schemas/index.schema.json     # JSON 出力の契約（存在する場合）
```

ファイルが存在しない場合は、そのルールについて推測しない。

## Review Process

### 1. Pull Request Context

以下を確認する。PR が無い場合は `git diff origin/develop...HEAD`（マージ先が別ならそのブランチ） で対象差分を特定する。

* PR title / description（対応する要件 ID が書かれているか）
* source / target branch
* changed files
* related Issue
* CI result

### 2. Branch Review

`.claude/rules/git-workflow.md` に従って、以下を確認する。

* ブランチ名（`feature/*` `fix/*` `release/*` `hotfix/*`）
* source / target の組み合わせが Git Flow 上の許可された遷移か（既定は `develop` 向け。別のマージ先なら理由があるか）
* `main` / `develop` への直接 push が無いか
* `release/*` に新機能が混入していないか
* マージ先とコンフリクトせず、マージが完了できる状態か（`mergeable` / CI）

### 3. Change Review

以下の優先順位で確認する。

1. Correctness（検証ルールの判定が requirements.md 7 章どおりか、誤検出・見逃し）
2. Safety（NFR-04：書き込み系の `--dry-run` と ChangePlan、`scripts/` を実行しない、GUI サーバーの localhost 限定・Host 検証・書き込み API 無し）
3. Determinism（NFR-05：出力の並び順、`schemaVersion`、OS 間で結果がそろうか）
4. Architecture（`.claude/rules/architecture.md`：一方向のレイヤー、診断は Rules 層だけ、Claude Code 固有情報は `internal/adapters/claude/`）
5. Error handling（読み取りエラーを握りつぶしていないか、終了コード 0 / 1 / 2）
6. Test coverage（正常系と異常系の `testdata/`、ルールごとのテスト）
7. Maintainability / Readability（`golangci-lint run` が通っているか）
8. Performance（NFR-01：100 SKILL で 1 秒程度。ただしキャッシュは入れない方針）

### 4. Project Rules

`.claude/rules/`・`docs/` に定義されたプロジェクト固有ルールを優先する。一般的なベストプラクティスとプロジェクトルールが異なる場合は、プロジェクトルールを根拠として指摘する。

次の変更は、ドキュメントの更新が同じ PR に含まれているかを確認する。

* 検証ルールの追加・変更 → requirements.md 7 章
* コマンド・オプションの追加・変更 → design.md 5.1
* データモデルの変更 → `schemas/index.schema.json` と design.md 8 章
* 決定済み事項（design.md 12 章・承認済みの ADR）に反する変更、または新しい設計判断 → ADR（`docs/adr/`）が含まれているか。`承認` の ADR が書き換えられていないか

### 5. Avoid False Positives

以下の場合、問題として断定せず `[Question]` として確認する。

* コードだけでは意図を判断できない
* 要件・設計に書かれていない
* 実行結果を確認できない
* 一般論だけを根拠とする問題

## Severity

* **Blocking**: マージ前に対応が必要（バグ、検証ルールの誤判定、NFR-04 の安全性違反、出力の非決定性、レイヤー構成の違反、要件・決定済み事項との矛盾、必須テストの欠落）
* **Non-blocking**: 改善を推奨するが、マージを妨げない（可読性、保守性、命名、軽微な設計改善、リファクタリング候補）
* **Question**: 実装意図や仕様の確認

人間がレビューコメントを書くときは、先頭に `[Blocking]` / `[Non-blocking]` / `[Question]` を付ける。

## Output

```markdown
## Summary
<レビュー概要>

## Blocking
- [file:line]
  <問題>
  <理由（根拠となる要件 ID・設計の章・ルール）>
  <対応案>

## Non-blocking
- [file:line]
  <改善提案>
  <理由>

## Questions
- [file:line]
  <確認事項>

## Positive
- <良かった点>
```

問題がない場合は、無理に指摘を作らない。

## Important Rules

* コードを変更しない
* ユーザーの明示的な依頼なしに PR を Approve / Merge しない
* 一般論だけで Blocking 判定しない
* プロジェクトドキュメントを確認せずにルール違反を断定しない
