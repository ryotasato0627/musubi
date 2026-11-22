# ADR-0006: レイヤー間の依存方向を depguard で強制する

- 状態: 承認
- 日付: 2026-11-22
- 決定者: プロジェクトオーナー（v0.1 計画の承認時）
- 関連: NFR-07、[docs/plans/v0.1.md](../plans/v0.1.md) D6・PR-1、ADR-0005

## 背景

design.md 1 章は、解析を `Scanner → Parser → Path Resolver → Index Builder → Analyzers → Rules → IndexSnapshot` の一方向に流し、診断は Rules 層だけで作ると定めている。また Claude Code 固有の情報は `internal/adapters/claude/` に閉じ込める（design.md 10 章）。開発は Agent が進めるため、文章のルールだけでは逆向きの import が紛れ込みやすい。

## 検討した案

### 案 A: golangci-lint の depguard でパッケージごとの import を制限する

- 良い点: 既に必須にしている golangci-lint（ADR-0005）に設定を足すだけで済む。CI と手元の両方で止められる
- 悪い点: パッケージを増やすたびに設定の更新が要る

### 案 B: go-arch-lint などの専用ツールを使う

- 良い点: レイヤー構成を宣言的に書け、表現力が高い
- 悪い点: ツールが 1 つ増え、CI と手元の導入作業が増える

### 案 C: レビューで確認する

- 良い点: 追加の設定が要らない
- 悪い点: 見落としが起き得る。機械的に止められない

## 決定

案 A を採用する。`.golangci.yml` の depguard で、design.md 1 章の流れに逆らう import と、adapters 以外からの Claude Code 固有情報への依存を禁止する。

## 影響

- `.golangci.yml` を PR-1 で追加し、逆向きの import が検出されることを確かめる
- パッケージを追加・移動する PR は depguard の設定も更新する
- 見直す条件: 設定の保守が重くなったら、案 B を検討する
