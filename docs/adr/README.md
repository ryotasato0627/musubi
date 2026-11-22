# ADR（Architecture Decision Record）

> 対象読者：設計判断の経緯を追いたい人、判断を覆す前に背景を確かめたい人
> 次に読むもの：決定済み事項の一覧 → [design.md 12 章](../spec/design.md#12-設計判断の記録)

musubi の設計判断を 1 件ずつ記録する。何を決めたかに加えて、なぜそうしたか、何と比べたかを残す。

- 書き方は `.claude/skills/adr-writer/SKILL.md`、形式は `.claude/skills/adr-writer/templates/adr.md`
- 決定済み事項の索引は [design.md 12 章](../spec/design.md#12-設計判断の記録)。詳細はこちらが正
- `承認` の ADR は書き換えない。判断を変えるときは新しい ADR で置き換える

| 番号 | タイトル | 状態 | 日付 |
| --- | --- | --- | --- |
| [0001](0001-record-decisions-as-adr.md) | 設計判断を ADR で記録する | 承認 | 2026-11-22 |
| [0002](0002-yaml-library.md) | frontmatter の YAML 解析に goccy/go-yaml を使う | 承認 | 2026-11-22 |
| [0003](0003-toml-library.md) | 設定ファイルの TOML 解析に BurntSushi/toml を使う | 承認 | 2026-11-22 |
| [0004](0004-go-version.md) | Go は最新の安定版を使う | 承認 | 2026-11-22 |
| [0005](0005-ci-github-actions.md) | CI は GitHub Actions で 3 OS を検証する | 承認 | 2026-11-22 |
| [0006](0006-depguard-layer-rules.md) | レイヤー間の依存方向を depguard で強制する | 承認 | 2026-11-22 |
