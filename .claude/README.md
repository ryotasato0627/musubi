# .claude（AI 開発フロー）

> 対象読者：このリポジトリで Claude Code を使って開発する人、開発フロー（Skill・Agent・ルール）を保守する人
> 次に読むもの：フロー全体 → [skills/dev-cycle/SKILL.md](skills/dev-cycle/SKILL.md)／計画の扱い → [docs/plans/](../docs/plans/README.md)

musubi は、計画の承認と PR のレビューだけを人間が行い、それ以外の工程を Claude Code の Skill と Agent で進める。Claude 向けの入口は [CLAUDE.md](../CLAUDE.md)、このファイルは人間向けの案内である。

## 流れ

```text
/release-plan v0.1          計画を作る（plan モード）→ 人間が承認 → docs/plans/v0.1.md
/dev-cycle [PR-n]           開発 → レビューで受け入れ条件を定義 → 修正と再レビュー → PR 作成
                            （人間が GitHub で PR をレビュー）
/dev-cycle resume <PR番号>  レビューでの判断を反映 → 再レビュー → push・返信
```

## 構成

| 場所 | 役割 |
| --- | --- |
| [skills/](skills/) | 手順。人間が `/名前` で呼ぶか、Claude が場面に応じて使う |
| [agents/](agents/) | 工程を担当するサブエージェント。dev-cycle が呼ぶ |
| [rules/](rules/) | 常に、または該当するパスを触るときに読み込まれる決まり |
| [hooks/](hooks/) | ツール実行前の機械的なチェック |
| [settings.json](settings.json) | 許可・禁止するコマンドと hook の登録 |

`work/`（工程間の受け渡し）と `settings.local.json` はコミットしない。

### Skill

| Skill | 使う場面 |
| --- | --- |
| [release-plan](skills/release-plan/SKILL.md) | リリース段階を PR 単位の計画にする。PR の粒度もここで決める |
| [dev-cycle](skills/dev-cycle/SKILL.md) | 計画の PR を 1 本、PR 作成まで止まらずに進める |
| [commit-strategy](skills/commit-strategy/SKILL.md) | コミットの粒度とメッセージの書式 |
| [pr-review](skills/pr-review/SKILL.md) | 差分を要件・設計・ルールに照らしてレビューする |
| [address-pr](skills/address-pr/SKILL.md) | PR のレビューコメントに対応する |
| [adr-writer](skills/adr-writer/SKILL.md) | 設計判断を ADR に記録する |

各 Skill の frontmatter には、musubi 自身の書式（`deps.requires` / `deps.related`）で Skill 間の依存を書いている。

### Agent

| Agent | 担当 | モデル |
| --- | --- | --- |
| [go-dev](agents/go-dev.md) | Go の実装・修正 | Opus |
| [web-dev](agents/web-dev.md) | GUI（`web/`）の実装・修正。v0.3 から | Opus |
| [test-writer](agents/test-writer.md) | 既存実装へのテスト追加 | Opus |
| [reviewer](agents/reviewer.md) | 受け入れ条件の定義と判定。毎回新しく起動し、開発の経緯を渡さない | Sonnet |
| [pr-creator](agents/pr-creator.md) | PR の作成と、マージできる状態かの確認 | Sonnet |

### ルール

| ファイル | 読み込まれるとき |
| --- | --- |
| [project.md](rules/project.md)、[git-workflow.md](rules/git-workflow.md) | 常に |
| [architecture.md](rules/architecture.md)、[cli.md](rules/cli.md)、[lint-rules.md](rules/lint-rules.md)、[output-schema.md](rules/output-schema.md)、[gui-server.md](rules/gui-server.md)、[testing.md](rules/testing.md)、[docs.md](rules/docs.md) | 各ファイル先頭の `paths` に当てはまるファイルを触るとき |

## 保守するときの注意

- 同じ内容を複数の場所に書かない。検証コマンドは `rules/git-workflow.md`、ADR の状態の決め方は `skills/adr-writer/SKILL.md` が正で、他は参照する
- この説明は `skills/` の中ではなくここに置く（`skills/` は musubi の探索対象なので、SKILL 以外のファイルを置かない）
