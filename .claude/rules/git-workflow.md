# ブランチ運用と検証コマンド

## ブランチ（Git Flow）

| ブランチ | 用途 | 分岐元 | マージ先 |
| --- | --- | --- | --- |
| `main` | リリース済みのコード。タグ（`v0.1.0` など）を打つ | - | - |
| `develop` | 次のリリースに向けた統合先 | `main` | - |
| `feature/<内容>` | 機能追加・改善（例 `feature/f04-lint-e007`） | `develop` | `develop` |
| `fix/<内容>` | `develop` 上のバグ修正 | `develop` | `develop` |
| `release/<版>` | リリース準備（例 `release/v0.1.0`）。バグ修正とドキュメントだけ | `develop` | `main` と `develop` |
| `hotfix/<内容>` | リリース済みの緊急修正 | `main` | `main` と `develop` |

- `main` / `develop` に直接コミット・push しない。force push しない
- `release/*` に新機能を入れない
- ブランチ名の `<内容>` は英小文字とハイフン。対応する要件 ID があれば入れてよい

## PR のマージ先

- 既定のマージ先は `develop`（`release/*` と `hotfix/*` は上の表に従う）
- `develop` は必須ではない。リモートに `develop` が無い、またはユーザーが別のブランチを指示した場合はそれに従う。`develop` が無いときは勝手に作らず、ユーザーに確認する
- **どのブランチを選んでも、マージが完了できることを確認する。** PR 作成前と push 後に次を確かめ、満たさなければ止めて報告する
  - `git fetch` 後、`git merge-tree --write-tree origin/<base> HEAD` がコンフリクトなしで終わる（コンフリクトがあれば、マージ先を取り込んで解消する。rebase して force push はしない）
  - PR 作成後、`gh pr view <N> --json mergeable,mergeStateStatus` が `MERGEABLE` で、`mergeStateStatus` が `DIRTY` / `BLOCKED` でない
  - CI がある場合は `gh pr checks <N>` がすべて成功している

## 検証コマンド（コミット前・push 前に必ず実行）

```bash
go build ./...
go vet ./...
golangci-lint run
go test ./...
```

- `golangci-lint` は必須。設定は `.golangci.yml` に置く（無ければ最初の実装 PR で追加する）
- 失敗したままコミット・push しない。実行できなかった場合は、成功と報告せず理由を書く
