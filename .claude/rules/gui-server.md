---
paths:
  - "internal/server/**"
  - "web/**"
---

# GUI・ローカルサーバーのルール（design.md 9 章、NFR-04）

- GUI は表示とフィルタだけ。解析ロジックを持たない（`IndexSnapshot` を受け取って描くだけ）
- API は `GET /api/index` と `GET /api/events`（SSE）の 2 本だけ。**書き込み API は作らない**（v1.0 まで）

## セキュリティ（必須）

- `127.0.0.1` に明示的にバインドする（`0.0.0.0` や名前解決に任せない）
- `Host` ヘッダーが `127.0.0.1:<port>` / `localhost:<port>` 以外なら拒否（DNS リバインディング対策）
- CORS ヘッダーは出さない
- SSE の同時接続数に上限を設け、ファイル監視はサーバー全体で 1 つにする

## watch

- 監視対象は探索ルート配下と `.musubi/config.toml`（`.gitignore` 除外対象は監視しない）
- 約 300 ms のデバウンス後に全体を再解析し、SSE で新しい Snapshot の番号を通知する

## フロントエンド

- Vite + Cytoscape.js + dagre。ビルド成果物は `go:embed` でバイナリに同梱
- エッジの線種：実線 `requires`、破線 `related`、点線 推論。循環は赤で強調。ノード色は ok / warn / error
