# musubi プロジェクト共通ルール

musubi は、リポジトリ内の Claude Code SKILL（`.claude/skills/`）を一覧・検証・作成・改名・削除し、SKILL 間の依存関係を解析・可視化する Go 製 CLI。

## 正とするドキュメント

- 要件：[docs/spec/requirements.md](../../docs/spec/requirements.md)（何を作るか）
- 設計：[docs/spec/design.md](../../docs/spec/design.md)（どう作るか）
- 実装がドキュメントと食い違う場合は、勝手にどちらかへ寄せず、ユーザーに確認する。仕様を変えたら該当ドキュメントも同じ変更で更新する
- 設計判断の詳細（背景・比較した案）は [docs/adr/](../../docs/adr/README.md)。新しい判断は `adr-writer` Skill で ADR に記録する
- 要件・ルールは ID で参照する：機能 `F-xx`、非機能 `NFR-xx`、検証ルール `E0xx` / `W0xx` / `I0xx`。コミットメッセージ・コメント・テスト名でも ID を使う

## 決定済み事項（design.md 12 章と承認済みの ADR。覆すときは必ず相談し、新しい ADR を書く）

- 実装言語は Go。バイナリ名・モジュール名・設定フォルダ名は `musubi` / `.musubi/` でそろえる
- 対象は Claude Code の SKILL のみ。探索ルートは `.claude/skills/` に固定（`--root` はテスト用に 1 つだけ上書き可）。`~/.claude/skills` は扱わない
- 依存宣言は frontmatter の `metadata` 直下のフラットキー：`deps.requires` / `deps.related` / `deps.ignore`。値はカンマ区切り文字列
- `requires` は「人間と CI のための宣言」。Claude Code が自動ロードするわけではない
- 設定ファイルは `.musubi/config.toml`。frontmatter プロファイルの既定は `claude-code`、`--max-warnings` の既定は無制限
- GUI は localhost の閲覧専用（v1.0 まで書き込み API なし）
- ブランチは Git Flow、静的解析は `golangci-lint` 必須（詳細は `.claude/rules/git-workflow.md`）
- 外部取り込み（F-15）は今は設計しない

## スコープ外（v1）

外部レジストリからの取り込み・GUI 上での編集・SKILL の実行や evals・`scripts/` の実行・ネットワークアクセス（NFR-03）。これらを実装しようとしたら止めて確認する。

## 進め方

`dev-cycle` Skill の実行中は、このファイルや他のルールにある「ユーザーに確認する」「止めて確認する」場面でも止まらない。その部分は実装せず、PR の「人間に判断してほしい点」に選択肢と推奨を書いて回す。

リリース計画（requirements.md 9 章）の段階に沿って実装する。v0.1 は F-01〜F-07・F-18（`init` `list` `show` `new` `lint` `deps`）。後の段階の機能を先回りして作らない。
