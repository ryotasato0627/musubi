# CLAUDE.md

musubi は、リポジトリ内の Claude Code SKILL（`.claude/skills/`）を一覧・検証・作成・改名・削除し、SKILL 間の依存関係を解析・可視化する Go 製 CLI。

このファイルには、リポジトリに実在するものだけを書く。まだ無いファイルやパッケージの決まりは `.claude/rules/` の該当ファイル（そのパスを触るときに読み込まれる）に書く。

## 現状と参照先

- 実装の進み具合：[docs/plans/v0.1.md](docs/plans/v0.1.md)（v0.1 は `init` `list` `show` `new` `lint` `deps`）
- 何を作るか：[docs/spec/requirements.md](docs/spec/requirements.md)／どう作るか：[docs/spec/design.md](docs/spec/design.md)（リポジトリ構成の予定は 11 章、決定済み事項は 12 章）
- 設計判断の経緯：[docs/adr/](docs/adr/README.md)
- 開発フローの全体像（人間向け）：[.claude/README.md](.claude/README.md)

## 検証

コミット前・push 前に `.claude/rules/git-workflow.md` の検証コマンド（build・vet・golangci-lint・test）をすべて通す。失敗したままコミットしない。実行できなければ成功と報告しない。

## アーキテクチャの要点

解析はコアで 1 回だけ行い、確定した読み取り専用の `IndexSnapshot` を CLI・エクスポート・GUI が共有する（どれも依存関係を再計算しない）。

```
Scanner → Parser → Path Resolver → Index Builder → Analyzers → Rules → IndexSnapshot
                                                                         ├→ CLI（表・JSON・SARIF・終了コード）
                                                                         ├→ export（Mermaid・DOT・JSON）
                                                                         └→ server（musubi ui、localhost・閲覧専用）
```

流れは一方向で、診断を作るのは Rules 層だけ。パッケージごとの責務は `.claude/rules/architecture.md`、CLI・出力・GUI・テストの決まりも `.claude/rules/` の各ファイルにある。

## 進め方

- 実装は承認済みの計画にある PR だけを、1 本ずつ `/dev-cycle` で進める（いきなり実装しない）。計画は `/release-plan <版>` で作る。人間が判断するのは計画の承認と PR のレビューだけ
- 実装とドキュメントが食い違ったら勝手に寄せず、ユーザーに確認する。仕様を変えたら docs も同じ変更で更新する
- 決定済み事項を覆すとき、スコープ外に触れるとき、後の段階の機能を先回りしそうなときは止めて確認する（`.claude/rules/project.md`）。ただし `dev-cycle` の実行中は止まらず、PR の「人間に判断してほしい点」に回す
- Git Flow：`develop` から `feature/*` / `fix/*` を切る。`main` / `develop` へ直接コミットしない、force push しない（`.claude/rules/git-workflow.md`）
