# 帳票と状態

## 3種類の帳票

| リソース | コマンド | エイリアス |
| --- | --- | --- |
| 請求書 | `misoca invoice` | — |
| 見積書 | `misoca estimate` | — |
| 納品書 | `misoca delivery-slip` | `misoca ds` |

3つとも `list` / `get <id>` / `create --json` / `pdf <id> -o <file>` を持ちます。
できることの差は次のとおりです。

- **請求書だけ**が状態遷移（`submit` / `unsubmit`、`pay` / `unpay`）と
  郵送依頼（`mail`）を持ちます。
- **見積書だけ**がメール送信（`distribute --json`）とロゴ・印影の取得
  （`logo` / `stamp`）を持ちます。
- 納品書は取得と作成と PDF のみです。

## 請求書の状態

請求書には独立した2つの状態があります。混同しないこと。

- **送付済みか**（`submit` / `unsubmit`、絞り込みは `--invoice-status
  submitted|unsubmitted`）
- **入金済みか**（`pay` / `unpay`、絞り込みは `--payment-status paid|unpaid`）

「未入金の請求書」は `--payment-status unpaid` であって、`--invoice-status` では
ありません。

ゴミ箱は上の2つとは別の軸です（`trash` / `restore`、絞り込みは `--type`）。

## 送り先と取引先グループ

- 送り先（`contact`）は取引先グループ（`contact-group`、エイリアス `cg`）に
  属します。
- どちらも削除ではなく**非表示**（`hide`）で、`restore` で戻せます。
- 一覧で非表示のものを見るには `--trashed` を付けます（帳票側の `--type` とは
  フラグ名が違うので注意）。

## 品目とユーザー

- `misoca item` は品目マスタ（`list` / `get` / `create`）。
- `misoca user me` は認証中のユーザー情報。認証が通っているかの確認に使えます。
