## PR #<N> レビュー対応方針

未解決スレッド: <X> 件(fix <a> / design <b> / question <c> / skip <d>)

### fix(承認で一括実装)

| # | 場所 | 指摘要約 | 対応 | コミット |
|---|---|---|---|---|
| F1 | internal/core/resolver/resolve.go:42 | `%20` を含むリンクが E007 になる | 結合前に URL デコードする + testdata 追加 | fix: ... |

### design(選択が必要)

#### D1: internal/analyzers/dependency.go:10 「循環の診断をここで作っている」

- 案 A: Analyzer は強連結成分を返すだけにし、E004 ルールで診断にする(推奨)— design.md 1 章「診断を作るのは Rules だけ」と一致。変更 2 ファイル
- 案 B: Analyzer のまま、重大度の上書きだけ後で適用する — 変更は小さいが、`[lint.rules]` の無効化が効く場所が 2 か所に分かれる
- 推奨: A。理由: ...

### question(返信のみ・resolve しない)

| # | 場所 | 返信文案 |
|---|---|---|
| Q1 | ... | ... |

### skip

| # | 場所 | 理由 |
|---|---|---|
| S1 | ... | <commit sha> で対応済み |

→ fix を進めてよいか、design は案 A / B のどちらにするか、question の返信文案でよいかを教えてください。
