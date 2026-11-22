# docs

> 対象読者：musubi の仕様と設計を理解したい人、仕様の変更をレビューする人
> 次に読むもの：初めてなら [spec/requirements.md](spec/requirements.md) → [spec/design.md](spec/design.md) の順

## フォルダーの役割

| 場所 | 内容 | 読む場面 |
| --- | --- | --- |
| [spec/requirements.md](spec/requirements.md) | 何を・なぜ作るか。機能 `F-xx`、非機能 `NFR-xx`、検証ルール `E/W/Ixxx` を定義する | 機能や検証ルールの定義を確かめるとき |
| [spec/design.md](spec/design.md) | どう作るか。要件 ID を引いて、レイヤー構成・CLI・データモデル・GUI を書く。12 章が決定済み事項の索引 | 実装の方針を確かめるとき |
| [adr/](adr/README.md) | 設計判断 1 件ごとの記録。背景・比較した案・決定・影響 | なぜその設計なのかを追うとき、判断を覆す前 |
| [plans/](plans/README.md) | リリース段階ごとの PR 単位の実装計画と、その実績 | 実装の進み具合を見るとき、計画を承認するとき |

## 正とするもの

- 要件と設計は `spec/` が正。実装と食い違ったら、どちらかに勝手に寄せず確認する
- 設計判断の詳細は `adr/` が正。design.md 12 章は 1 行の要約とリンクだけを持つ
- ID は振り直さない。廃止するときは欠番にする
