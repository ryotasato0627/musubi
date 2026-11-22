# musubi

> 対象読者：musubi を使いたい人、このリポジトリを初めて開いた人
> 次に読むもの：仕様を知りたい → [docs/](docs/README.md)／開発に参加する → [.claude/](.claude/README.md)

リポジトリ内の Claude Code SKILL（`.claude/skills/`）を一覧・検証・作成・改名・削除し、SKILL 間の依存関係を解析・可視化する Go 製 CLI。

SKILL が増えると「どの SKILL がどれを前提にしているか」「壊れた参照や仕様違反がないか」「消してよい SKILL はどれか」を目視で追えなくなる。musubi は依存関係を SKILL 自身の frontmatter に宣言させ、壊れた依存や循環を CI で止められるようにする。

## 現状

**開発中（v0.1 の実装前）。** まだ動くバイナリは無い。進み具合は [v0.1 実装計画](docs/plans/v0.1.md) を見る。

## できるようになること

| 段階 | コマンド・機能 |
| --- | --- |
| v0.1 | `init`（設定の生成）、`list`（一覧）、`show`（詳細）、`new`（雛形作成）、`lint`（検証。JSON・SARIF 出力）、`deps`（依存ツリー・逆引き・循環検出） |
| v0.2 | Mermaid・DOT・JSON へのエクスポート、本文中の言及の推論、改名、削除 |
| v0.3 | ブラウザでの依存グラフの閲覧（localhost・閲覧専用）、watch |
| v1.0 | `lint --fix`、JSON スキーマの固定 |

詳しくは [requirements.md 9 章](docs/spec/requirements.md#9-リリース計画)。

## 依存関係の宣言

SKILL の frontmatter の `metadata` に、カンマ区切りで書く。Agent Skills 仕様の範囲内なので、musubi を使わない環境でも SKILL は壊れない。

```yaml
---
name: api-review
description: API変更のレビュー手順。PRにAPI変更が含まれるときに使う。
metadata:
  deps.requires: "repo-conventions,test-writing"   # 前提とする SKILL（欠落・循環は error）
  deps.related: "adr-writer"                       # 関連する SKILL
---
```

`requires` は人間と CI のための宣言で、Claude Code が自動でロードするわけではない。

## ドキュメント

| 場所 | 内容 |
| --- | --- |
| [docs/](docs/README.md) | 要件・設計・設計判断（ADR）・実装計画 |
| [.claude/](.claude/README.md) | このリポジトリの AI 開発フロー（Skill・Agent・ルール） |
