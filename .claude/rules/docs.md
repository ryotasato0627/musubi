---
paths:
  - "docs/**"
---

# ドキュメント執筆のルール

- 日本語で書く。文体は「〜する」「〜である」の常体
- 役割を分ける：`docs/spec/` は仕様、`docs/adr/` は設計判断の記録、`docs/plans/` は実装計画。仕様のうち requirements.md は「何を・なぜ」、design.md は「どう実現するか」。design.md は要件 ID（`F-xx` / `NFR-xx` / `E/W/Ixxx`）を引いて書く
- ID は振り直さない。廃止するときは欠番にする
- 設計判断は `adr-writer` Skill で ADR（`docs/adr/NNNN-<slug>.md`）に記録し、design.md 12 章の表には 1 行の要約と ADR へのリンクを書く（見出し下の日付も更新する）
- `承認` の ADR は書き換えない。判断を変えるときは新しい ADR で置き換える
- 未確定の事項は「（要確認）」と明記する
- 不採用にした案は、理由とあわせて短く残す
