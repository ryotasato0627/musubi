---
paths:
  - "internal/cli/**"
  - "cmd/**"
---

# CLI のルール（design.md 5・6 章、NFR-04・NFR-06）

- コマンドは薄く保つ。ロジックは core / analyzers / rules に置き、コマンドは `IndexSnapshot` を作って Presenter に渡すだけ
- コマンド名・オプション名は design.md 5.1 の表に合わせる。勝手に増やす・改名するときは表も更新する
- 共通オプション：`--config` `--root`（1 つだけ）`--no-color` `--quiet` `--verbose`
- 終了コード：`0` 成功 ／ `1` lint で error（または `--max-warnings` 超過）／ `2` 使い方の誤り・設定エラー
- リストを受けるオプションはカンマ区切り（`--requires a,b`）。依存宣言の書式とそろえる

## 書き込み系コマンド（`new` `rename` `rm` `lint --fix`）

- 既定は読み取り専用。書き込み系は必ず `--dry-run` と差分表示を持つ
- 「`ChangePlan` を作る → 表示 → 適用」の 2 段構え。計画作成と適用を混ぜない
- 適用手順（design.md 6 章）を守る：ハッシュ・移動先・権限の事前確認 → 作業ツリーが dirty なら停止（`--allow-dirty`）→ `.musubi/backup/<日時>/` へ退避 → 一時ファイルに書いて置き換え → ディレクトリ改名は最後 → 失敗時は退避から復元 → 成功時は退避を削除
- `rename` は本文中の言及（mention）を書き換えない
- `rm` は他 SKILL の依存宣言を自動で消さない。被依存があれば拒否（`--force` で続行）
- `scripts/` 配下のコードは決して実行しない
