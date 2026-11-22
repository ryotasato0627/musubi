# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 現状

musubi は、リポジトリ内の Claude Code SKILL（`.claude/skills/`）を一覧・検証・作成・改名・削除し、SKILL 間の依存関係を解析・可視化する Go 製 CLI。

現時点では実装コードは無く、`docs/requirements.md`（何を作るか）と `docs/design.md`（どう作るか）だけがある。git リポジトリもまだ初期化されていない。最初の実装は v0.1（F-01〜F-07・F-18：`init` `list` `show` `new` `lint` `deps`）で、`go.mod`・`.golangci.yml` もその最初の PR で追加する。リポジトリ構成の予定は design.md 11 章にある。

## コマンド

コミット前・push 前に必ずすべて通す（失敗したままコミットしない。実行できなければ成功と報告しない）。

```bash
go build ./...
go vet ./...
golangci-lint run
go test ./...
```

- 単一テスト：`go test ./internal/rules -run 'TestE007_CaseMismatch' -v`（テスト名には要件 ID・ルール ID を入れる規約）
- ゴールデンファイル（`testdata/**/*.golden`、JSON / SARIF / Mermaid / DOT 出力）の更新：`go test ./... -update`
- 手元での動作確認：`go run ./cmd/musubi lint --root testdata/<ケース名>/.claude/skills`（`--root` はテスト用に探索ルートを 1 つだけ上書きする）

## アーキテクチャの要点

解析はコアで 1 回だけ行い、確定した読み取り専用の `IndexSnapshot` を CLI・エクスポート・GUI が共有する（どれも依存関係を再計算しない）。

```
Scanner → Parser → Path Resolver → Index Builder → Analyzers → Rules → IndexSnapshot
                                                                         ├→ CLI（表・JSON・SARIF・終了コード）
                                                                         ├→ export（Mermaid・DOT・JSON）
                                                                         └→ server（musubi ui、localhost・閲覧専用）
```

- 流れは一方向。逆向きの依存を作らない
- **`Diagnostic` を作るのは `internal/rules` だけ。** Scanner / Parser のエラーは事実として記録し、Rules が E001・E008 などに変換する。重大度の上書き・無効化もここで一括適用する（NFR-07、1 ルール 1 ファイル）
- Analyzers は意味的解析だけで診断を作らない。`graph` はエッジの意味を解釈しない汎用処理
- `internal/cli` は cobra のコマンド定義だけでロジックを持たない
- Claude Code 固有の情報は `internal/adapters/claude/` に閉じ込める
- 書き込み系（`new` `rename` `rm` `--fix`）は「変更計画 → 表示 → 適用」の 2 段構えで `--dry-run` を実現する（design.md 6 章）
- 終了コード：`0` 成功 ／ `1` lint error（または `--max-warnings` 超過）／ `2` 使い方・設定エラー
- JSON 出力は `schemaVersion` を持ち、並び順を決定的にする。契約は `schemas/index.schema.json`

## ルールとワークフロー

詳細なルールは `.claude/rules/` にあり、`project.md`・`git-workflow.md` は常時、それ以外は該当パスを触るときに読み込まれる（`architecture.md`・`cli.md`・`lint-rules.md`・`output-schema.md`・`gui-server.md`・`testing.md`・`docs.md`）。ここに書いていない決まりはそちらを参照する。特に：

- 実装とドキュメントが食い違ったら勝手に寄せず、ユーザーに確認する。仕様を変えたら docs も同じ変更で更新する
- 決定済み事項（design.md 12 章）を覆すとき、スコープ外（ネットワークアクセス・`scripts/` 実行・GUI 書き込みなど）に触れるとき、後の段階の機能を先回りしそうなときは止めて確認する
- Git Flow：`develop` から `feature/*` / `fix/*` を切り、PR の既定マージ先は `develop`。`main` / `develop` へ直接コミットしない、force push しない

## 開発サイクル

実装は `dev-cycle` Skill で 1 PR ずつ進める。人間が判断するのは PR のレビューだけで、それ以外は止まらずに進める。

```
/dev-cycle <タスク>        0 準備 → 1 開発 → 2 初回レビュー（受け入れ条件を定義）→ [3 修正 → 4 再レビュー]×最大3回 → 5 PR 作成 → 終了
（人間が PR をレビュー）
/dev-cycle resume <PR番号>  7 判断を反映 → 8 再レビュー → push・返信 → 終了
```

- 進行役はメインのセッション（サブエージェントは別のサブエージェントを呼べない）。各工程は `.claude/agents/` の `go-dev` / `web-dev`（開発・修正）、`test-writer`、`reviewer`（受け入れ条件の定義と判定）、`pr-creator` が担当する
- 工程間の受け渡しは `.claude/work/<ブランチ>/`（`task.md`・`ac.md`・`review-<N>.md` など。コミットしない）
- **reviewer は毎回新しいインスタンスで起動し、開発のコンテキストを渡さない**（SendMessage での継続・fork は使わない。プロンプトはモード・作業ディレクトリ・ブランチ・round だけ。`dev-report` / `fix-report` / 会話の要約は渡さない）。reviewer はファイルを書かず、返した `ac.md` / レビューをメインが保存する
- 受け入れ条件が 3 回で満たせない、仕様の判断が要る、コンフリクトが解消できない場合も止まらず、Draft PR と「人間に判断してほしい点」に回す
- `dev-cycle` の実行中は、上の「ユーザーに確認する」「止めて確認する」場面でも止まらない。その部分は実装せず、PR の「人間に判断してほしい点」に選択肢と推奨を書いて回す
- 関連 Skill：`pr-strategy`（PR の粒度と計画）、`commit-strategy`、`pr-review`、`address-pr`
