# misoca CLI — AIエージェント向けリファレンス

`misoca` は Misoca API を操作する CLI です。請求書・見積書・納品書の作成と取得、
送り先・取引先グループ・品目の管理を、すべて非対話で実行できます
（`misoca auth` の初回認証だけはブラウザを使います）。

すべて `misoca <リソース> <メソッド> [フラグ]` の形式です。出力は JSON なので
`jq` と組み合わせて処理できます。

このリファレンスはバイナリに埋め込まれています。`misoca llm` は常に実行中の
バージョンそのものを説明します。

## 認証

Misoca は OAuth 2.0 で認証します。初回だけブラウザでの許可が必要です。

```bash
misoca auth
```

- `MISOCA_CLIENT_ID` と `MISOCA_CLIENT_SECRET` が環境変数にあれば、即座に
  ブラウザ認証が始まります。
- 無ければ対話形式で、Misoca の開発者ページでアプリケーションを作成する手順が
  案内されます（コールバックURL `http://localhost:18080/callback`）。
- トークンは `~/.config/misoca-cli/token.json` に保存され、以降は自動で
  リフレッシュされます。

**エージェントは `misoca auth` を勝手に実行しないこと。** ブラウザが開き、
ユーザーの操作を待って止まります。未認証だと分かったら、ユーザーに
`misoca auth` の実行を依頼してください。

現在の認証状態は次で確認できます。

```bash
misoca user me
```

## 出力とリクエストボディ

- 出力は **stdout に JSON**。エラーは stderr に `Error:` 付きで出て、終了コードは 1。
- 作成系は `--json '<JSON>'` でリクエストボディを渡します。個別フラグと併用でき、
  **同じフィールドを両方で指定した場合はフラグが優先**されます。
- PDF・ロゴ・印影などのバイナリは `-o <path>` でファイルに保存します。

```bash
misoca invoice list --payment-status unpaid | jq '.[] | {id, subject, total_amount}'
misoca invoice pdf 12345 -o invoice.pdf
```

## 一覧の絞り込みとページネーション

- `--page` / `--per-page`（**1ページ最大100件**）。CLI は自動で全ページを辿りません。
  件数が多い場合は `--page` を進めて自分で繰り返してください。
- `--type` を渡さない場合、CLI は `type` パラメータを送らず **API 既定の
  `active`** が使われます。ゴミ箱以外をすべて見たいときは明示的に
  `--type untrashed` を指定してください。「古い請求書が出てこない」は
  たいていこれが原因です。
- 請求書の絞り込みは `--condition`（請求書番号・取引先名・件名・社内メモの
  キーワード検索）、`--payment-status`（`paid` / `unpaid`）、
  `--invoice-status`（`submitted` / `unsubmitted`）、日付範囲
  （`--from` / `--to`、`--due-date-from` / `--due-date-to`、
  `--updated-from` / `--updated-to`、いずれも `YYYY/MM/DD`）、
  並び順（`--order` `asc`/`desc`、`--order-by` `created_at`/`updated_at`/
  `issue_date`/`payment_due_on`）が使えます。

## 取り扱いに注意が必要な操作

この CLI は**実在の取引先に対する実際の帳票**を操作します。実行前にユーザーの
同意を取るべきものを挙げます。

| 操作 | 何が起きるか |
| --- | --- |
| `misoca invoice mail <id>` | **郵送を依頼します。費用が発生し、取り消せません。** |
| `misoca estimate distribute <id>` | 取引先へメールが送信されます。 |
| `invoice submit` / `unsubmit` | 送付済みフラグの変更。会計上の状態が変わります。 |
| `invoice pay` / `unpay` | 入金済みフラグの変更。 |
| `invoice trash` / `restore` | ゴミ箱への移動と復元（`restore` で戻せます）。 |
| `contact hide` / `restore` | 送り先の非表示と再表示。 |

`create` 系は下書きを作るだけですが、それでも何をどの取引先に対して作るのかを
提示してから実行してください。
