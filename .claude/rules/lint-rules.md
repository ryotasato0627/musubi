---
paths:
  - "internal/rules/**"
  - "internal/analyzers/**"
---

# 検証ルール実装のルール（requirements.md 7 章、design.md 4 章）

- 1 ルール 1 ファイル（`internal/rules/e001.go`, `w004.go` …）。ルール追加がそのファイルと登録だけで済むようにする（NFR-07）
- ルール ID・既定の重大度・意味は requirements.md 7 章の表が正。新しいルールを足すときは先に表へ追記する
- 重大度の上書き・無効化（`[lint.rules]`）は Rules 層で一括適用する。個々のルールに分岐を書かない
- 終了コード：error が 1 件でもあれば 1。warning だけでは失敗にしない（`--max-warnings` 指定時を除く）
- 診断の並び順・主な位置は決定的にする（NFR-05）

## 依存解析の決まり

- 循環検出（E004）は `requires` だけが対象。`related`・推論エッジ（`file-ref` / `mention`）・`ambiguous` なエッジは含めない
- 自己依存も循環。強連結成分ごとに診断 1 件。主な位置はグループ内で ID が最小の SKILL の宣言行、他メンバーは `Diagnostic.Related`
- `related` は対称。両側に同じ宣言があれば 1 本に統合する
- `ambiguous` は E005 を出し、循環検出・到達性から外す

## 推論（F-12、v0.2〜）

- `mention` の信頼度：スラッシュコマンド形式・バッククォート内＝高、「〜スキル」「〜 SKILL」＝中、通常の文中＝低（対象外）
- 抑制コメント `<!-- musubi-ignore mention -->`（直後の 1 ブロック）、`<!-- musubi-ignore-file mention -->`（ファイル全体）、`metadata.deps.ignore`（SKILL 全体）。「その行だけ」の抑制は作らない
- I001 は v0.2 で F-12 と同時に有効化する。W004 は v0.1 から有効
