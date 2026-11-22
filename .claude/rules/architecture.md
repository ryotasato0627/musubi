---
paths:
  - "cmd/**"
  - "internal/**"
---

# アーキテクチャのルール（design.md 1・2・3・4・10 章）

## レイヤーと一方向の流れ

`Scanner → Parser → Path Resolver → Index Builder → Analyzers → Rules → IndexSnapshot` の一方向。逆向きの依存を作らない。

| パッケージ | やること | やらないこと |
| --- | --- | --- |
| `internal/core/scanner` | `<root>/<name>/SKILL.md` の列挙、`.gitignore` 尊重（go-git の gitignore パッケージ。git コマンドに依存しない） | SKILL ディレクトリ内でさらに SKILL を探す |
| `internal/core/parser` | frontmatter（YAML）と Markdown AST（goldmark）からリンク・言及を抽出 | パス解決、存在確認 |
| `internal/core/resolver` | リンクの分類・正規化・SKILL 境界の判定 | 診断の生成 |
| `internal/core/index` / `graph` | 索引の構築、汎用グラフ処理（`FindCycles` `Reachable` `ReverseEdges`） | エッジの意味の解釈 |
| `internal/analyzers` | 依存・参照・推論の意味的解析 | **診断の生成** |
| `internal/rules` | Analyzer の結果を `Diagnostic` に変換、重大度の上書き・無効化 | 解析そのもの |
| `internal/cli` | 引数解析とコマンド定義（cobra） | ロジック |

## 守ること

- **`Diagnostic` を作るのは `internal/rules` だけ。** Scanner や Parser のエラー（読み取り失敗、YAML 不正）は事実として記録し、Rules が E001・E008 などに変換する
- `IndexSnapshot` は生成後は読み取り専用。CLI・エクスポート・GUI は依存関係を再計算しない
- 解析は純粋関数に寄せる（FS 読み取り → スナップショット → 診断）
- キャッシュや差分解析は入れない。毎回全体を解析する（NFR-01）
- SKILL は内部で実体パス由来の `ID` で識別し、`name` は表示と依存宣言だけに使う。名前解決は `resolved` / `missing` / `ambiguous`
- Claude Code 固有の情報（`.claude/skills` のパス、拡張 frontmatter キー）は `internal/adapters/claude/` に閉じ込め、コアから直接参照しない
- 読み取りエラーを握りつぶさない（E008 にして探索は続行）
- ネットワークに出ない。外部 URL のリンクは検証しない（NFR-03）

## パス解決の決まり（3 章）

- 抽出対象：インラインリンク・参照形式リンク・画像。コードブロック・インラインコード・HTML コメント・生 HTML の `href` は対象外
- 正規化：`#`・`?` 以降を除去 → URL デコード → `SKILL.md` のディレクトリ基準で結合し `..` を正規化 → 自分／他 SKILL／どれでもない を判定
- 大文字小文字が実ファイル名と完全一致しなければ E007（OS 間で結果をそろえる）
- シンボリックリンクはルート直下のディレクトリへのものだけ解決する。SKILL 内部のリンクはたどらない
