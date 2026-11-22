---
paths:
  - "internal/export/**"
  - "schemas/**"
  - "internal/core/index/**"
---

# 出力とデータモデルのルール（design.md 5.3・8 章、NFR-05・NFR-06）

- `schemas/index.schema.json` が JSON 出力の契約。Go の構造体（`Skill` `Edge` `Diagnostic` `IndexSnapshot` など）と GUI の TS 型はこれに従う。構造体を変えたらスキーマも同じ変更で更新する
- JSON 出力には必ず `schemaVersion` を持たせる。互換性を壊す変更ではバージョンを上げる
- 並び順は決定的にする（git の差分が読めるように）。map をそのまま列挙して出力しない
- `metadata` はいったん `any` で読み、文字列→文字列を検証してから `map[string]string` に入れる（違反は E009）
- `musubi graph --format json` と GUI の `GET /api/index` は同一の内容を返す

## SARIF（2.1.0）

- level 対応：error → `error`、warn → `warning`、info → `note`
- `artifactLocation.uri` はリポジトリルートからの相対パス。`startLine` / `startColumn` は 1 始まり
- 循環グループの他メンバーは `relatedLocations`
- ファイル位置を持たない診断は探索ルートのパスに紐づける
