# ADR-0005: CI は GitHub Actions で 3 OS を検証する

- 状態: 承認
- 日付: 2026-11-22
- 決定者: プロジェクトオーナー（v0.1 計画の承認時）
- 関連: NFR-02、NFR-03、NFR-06、[docs/plans/v0.1.md](../plans/v0.1.md) D5・PR-1、ADR-0004、ADR-0006

## 背景

musubi は macOS・Linux・Windows で動くこと（NFR-02）を求められ、大文字小文字やパス区切り、シンボリックリンクの扱いで OS 差が出やすい（design.md 3 章）。開発は Agent が止まらずに進める（dev-cycle）ため、PR ごとに機械的に品質を確かめる仕組みが要る。

## 検討した案

### 案 A: PR ごとに ubuntu / macos / windows の 3 OS で検証する

- 良い点: OS 差による不具合を PR の段階で止められる
- 悪い点: CI の実行時間と、macOS・Windows ランナーの消費が増える

### 案 B: PR は ubuntu だけ、develop への push で 3 OS

- 良い点: PR の CI が速く、消費が少ない
- 悪い点: OS 差の不具合がマージ後に見つかり、手戻りになる

## 決定

案 A を採用する。GitHub Actions で、PR ごとに 3 OS で次を実行する。

- `go build ./...`、`go vet ./...`、`golangci-lint run`、`go test ./...`
- `govulncheck ./...`
- `go mod tidy` で差分が出ないことの確認

## 影響

- `.github/workflows/ci.yml` を PR-1 で追加する
- 手元でも同じ検証コマンドを通す（`.claude/rules/git-workflow.md`）
- 見直す条件: CI の時間や消費が問題になったら、案 B への切り替えを検討する
