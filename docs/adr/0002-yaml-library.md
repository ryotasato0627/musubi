# ADR-0002: frontmatter の YAML 解析に goccy/go-yaml を使う

- 状態: 承認
- 日付: 2026-11-22
- 決定者: プロジェクトオーナー（v0.1 計画の承認時）
- 関連: F-04、F-06、E001、E009、NFR-05、[docs/plans/v0.1.md](../plans/v0.1.md) D1・PR-4

## 背景

SKILL.md の frontmatter を解析し、E001（YAML 不正）・E002・E006・E009（`metadata` の値が文字列でない）などの診断を出す。診断は SARIF や表にファイル・行・列を出す（design.md 5.3・8 章の `Location`）ため、YAML のノードごとの位置が取れる必要がある。`metadata` は値の型を検証するため、いったん型を決めずに読める必要もある（design.md 8 章）。

## 検討した案

### 案 A: github.com/goccy/go-yaml

- 概要: AST を公開している YAML ライブラリ
- 良い点: ノードごとに行・列が取れる。構文エラーにも位置が付く。保守が続いている
- 悪い点: yaml v3 系に比べて利用実績が少ない

### 案 B: yaml v3 系（go.yaml.in/yaml/v3）

- 概要: `yaml.Node` で位置を取れる定番ライブラリの後継
- 良い点: 利用実績が多い。`yaml.Node` の API が広く知られている
- 悪い点: 元の `gopkg.in/yaml.v3` は保守が終了しており、後継への移行状況を見ながら使うことになる

## 決定

案 A（goccy/go-yaml）を採用する。位置つきの AST とエラーが標準で取れ、保守が続いている点を重視した。

## 影響

- `internal/core/parser/frontmatter.go` は goccy/go-yaml の AST から値と位置を取り出す（PR-4）
- ライブラリの型を parser の外に出さない（parser の戻り値は musubi のデータモデルにする）
- 見直す条件: 位置情報の不具合や保守の停止が起きたら、案 B への置き換えを検討する
