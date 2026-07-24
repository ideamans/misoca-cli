# コマンドカタログ

cobra のコマンド定義から `go generate ./...` で生成しています。
手編集しないこと。直すのはコマンド定義側です。

## `misoca auth`

Misoca APIのOAuth2認証を行います

Misoca APIのOAuth2認証を行います。

環境変数 MISOCA_CLIENT_ID, MISOCA_CLIENT_SECRET が設定されていれば
即座にブラウザ認証を開始します。

未設定の場合はアプリケーション作成から案内します。
トークンは ~/.config/misoca-cli/token.json に保存され、自動的にリフレッシュされます。

## `misoca contact`

送り先の操作

Aliases: contacts

### `misoca contact create`

送り先を作成

| flag | type | default | description |
| --- | --- | --- | --- |
| `--address1` | string | — | 住所1 |
| `--address2` | string | — | 住所2 |
| `--code` | string | — | 取引先コード |
| `--contact-group-id` | int | `0` | 取引先グループID |
| `--fax` | string | — | FAX番号 |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--mail` | string | — | メールアドレス |
| `--mail-cc` | string | — | CCメールアドレス |
| `--memo` | string | — | 社内メモ |
| `--name` | string | — | 送り先名 |
| `--name1` | string | — | 宛名1 |
| `--name2` | string | — | 宛名2 |
| `--name3` | string | — | 宛名3 |
| `--name4` | string | — | 宛名4 |
| `--notes` | string | — | 備考 |
| `--ruby` | string | — | フリガナ |
| `--tax-option` | string | — | 税表示 (USE_SENDER/INCLUDE/EXCLUDE/EXEMPT/NONENTRY) |
| `--tel` | string | — | 電話番号 |
| `--title` | string | — | 敬称 (御中/様) |
| `--title2` | string | — | 敬称2 |
| `--zip-code` | string | — | 郵便番号 |

### `misoca contact get`

送り先の詳細を取得

```
misoca contact get <id>
```

### `misoca contact hide`

送り先を非表示にする

```
misoca contact hide <id>
```

### `misoca contact list`

送り先の一覧を取得

| flag | type | default | description |
| --- | --- | --- | --- |
| `--contact-group-id` | string | — | 取引先グループIDでフィルタ |
| `--trashed` | bool | `false` | 非表示の送り先を取得 |

### `misoca contact restore`

送り先を再表示する

```
misoca contact restore <id>
```

## `misoca contact-group`

取引先グループの操作

Aliases: contact-groups, cg

### `misoca contact-group create`

取引先グループを作成

| flag | type | default | description |
| --- | --- | --- | --- |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--name` | string | — | 取引先名 (必須) |
| `--peppol-id` | string | — | Peppol ID |
| `--ruby` | string | — | フリガナ |
| `--tax-option` | string | — | 税表示 (USE_SENDER/INCLUDE/EXCLUDE/EXEMPT/NONENTRY) |
| `--title` | string | — | 敬称 (御中/様) |

### `misoca contact-group get`

取引先グループの詳細を取得

```
misoca contact-group get <id>
```

### `misoca contact-group hide`

取引先グループを非表示にする

```
misoca contact-group hide <id>
```

### `misoca contact-group list`

取引先グループの一覧を取得

| flag | type | default | description |
| --- | --- | --- | --- |
| `--trashed` | bool | `false` | 非表示の取引先グループを取得 |

### `misoca contact-group restore`

取引先グループを再表示する

```
misoca contact-group restore <id>
```

## `misoca delivery-slip`

納品書の操作

Aliases: delivery-slips, ds

### `misoca delivery-slip create`

納品書を作成

| flag | type | default | description |
| --- | --- | --- | --- |
| `--contact-id` | int | `0` | 送り先ID (必須) |
| `--delivery-date` | string | — | 納品日 (YYYY/MM/DD) |
| `--delivery-slip-number` | string | — | 納品書番号 |
| `--issue-date` | string | — | 発行日 (必須, YYYY/MM/DD) |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--subject` | string | — | 件名 |

### `misoca delivery-slip get`

納品書の詳細を取得

```
misoca delivery-slip get <id>
```

### `misoca delivery-slip list`

納品書の一覧を取得

| flag | type | default | description |
| --- | --- | --- | --- |
| `--contact-group-id` | string | — | 取引先グループIDでフィルタ |
| `--page` | int | `0` | ページ番号 |
| `--per-page` | int | `0` | 1ページあたりの件数 (最大100) |
| `--type` | string | — | 種類 (active/archived/trashed/untrashed) |

### `misoca delivery-slip pdf`

納品書のPDFを取得

```
misoca delivery-slip pdf <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `-o`, `--output` | string | — | 出力ファイルパス |

## `misoca estimate`

見積書の操作

Aliases: estimates

### `misoca estimate create`

見積書を作成

| flag | type | default | description |
| --- | --- | --- | --- |
| `--contact-id` | int | `0` | 送り先ID (必須) |
| `--estimate-number` | string | — | 見積書番号 |
| `--expire-date` | string | — | 有効期限 (YYYY/MM/DD) |
| `--issue-date` | string | — | 見積日 (必須, YYYY/MM/DD) |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--subject` | string | — | 件名 |

### `misoca estimate distribute`

見積書をメールで送信

```
misoca estimate distribute <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `--cc-self` | bool | `false` | 自分をCCに含める |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--mail-body` | string | — | メール本文 |
| `--mail-subject` | string | — | メール件名 (必須) |

### `misoca estimate get`

見積書の詳細を取得

```
misoca estimate get <id>
```

### `misoca estimate list`

見積書の一覧を取得

| flag | type | default | description |
| --- | --- | --- | --- |
| `--contact-group-id` | string | — | 取引先グループIDでフィルタ |
| `--page` | int | `0` | ページ番号 |
| `--per-page` | int | `0` | 1ページあたりの件数 (最大100) |
| `--type` | string | — | 種類 (active/archived/trashed/untrashed) |

### `misoca estimate logo`

見積書のロゴ画像を取得

```
misoca estimate logo <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `-o`, `--output` | string | — | 出力ファイルパス |

### `misoca estimate pdf`

見積書のPDFを取得

```
misoca estimate pdf <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `-o`, `--output` | string | — | 出力ファイルパス |

### `misoca estimate stamp`

見積書の印影画像を取得

```
misoca estimate stamp <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `-o`, `--output` | string | — | 出力ファイルパス |

## `misoca invoice`

請求書の操作

Aliases: invoices

### `misoca invoice create`

請求書を作成

| flag | type | default | description |
| --- | --- | --- | --- |
| `--contact-id` | int | `0` | 送り先ID (必須) |
| `--invoice-number` | string | — | 請求書番号 |
| `--issue-date` | string | — | 請求日 (必須, YYYY/MM/DD) |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--payment-due-on` | string | — | 支払期限 (YYYY/MM/DD) |
| `--subject` | string | — | 件名 |

### `misoca invoice get`

請求書の詳細を取得

```
misoca invoice get <id>
```

### `misoca invoice list`

請求書の一覧を取得

| flag | type | default | description |
| --- | --- | --- | --- |
| `--condition` | string | — | キーワード検索 (番号・取引先名・件名・メモ) |
| `--contact-group-id` | string | — | 取引先グループIDでフィルタ |
| `--due-date-from` | string | — | 支払期限の開始日 (YYYY/MM/DD) |
| `--due-date-to` | string | — | 支払期限の終了日 (YYYY/MM/DD) |
| `--from` | string | — | 請求日の開始日 (YYYY/MM/DD) |
| `--invoice-status` | string | — | 請求状況 (submitted/unsubmitted) |
| `--order` | string | — | 並び順 (asc/desc) |
| `--order-by` | string | — | 並び順基準 (created_at/updated_at/issue_date/payment_due_on) |
| `--page` | int | `0` | ページ番号 |
| `--payment-status` | string | — | 入金状況 (paid/unpaid) |
| `--per-page` | int | `0` | 1ページあたりの件数 (最大100) |
| `--to` | string | — | 請求日の終了日 (YYYY/MM/DD) |
| `--type` | string | — | 種類 (active/archived/trashed/untrashed) |
| `--updated-from` | string | — | 更新日の開始日 (YYYY/MM/DD) |
| `--updated-to` | string | — | 更新日の終了日 (YYYY/MM/DD) |

### `misoca invoice mail`

請求書を郵送依頼する

```
misoca invoice mail <id>
```

### `misoca invoice pay`

請求書を入金済みにする

```
misoca invoice pay <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--paid-on` | string | — | 入金日 (YYYY/MM/DD) |

### `misoca invoice pdf`

請求書のPDFを取得

```
misoca invoice pdf <id>
```

| flag | type | default | description |
| --- | --- | --- | --- |
| `-o`, `--output` | string | — | 出力ファイルパス |

### `misoca invoice restore`

請求書をゴミ箱から復元

```
misoca invoice restore <id>
```

### `misoca invoice submit`

請求書を送付済みにする

```
misoca invoice submit <id>
```

### `misoca invoice trash`

請求書をゴミ箱に移動

```
misoca invoice trash <id>
```

### `misoca invoice unpay`

請求書の入金済みを取り消す

```
misoca invoice unpay <id>
```

### `misoca invoice unsubmit`

請求書の送付済みを取り消す

```
misoca invoice unsubmit <id>
```

## `misoca item`

品目の操作

Aliases: items

### `misoca item create`

品目を作成

| flag | type | default | description |
| --- | --- | --- | --- |
| `--exclude-withholding-tax` | bool | `false` | 源泉徴収対象外 |
| `--json` | string | — | リクエストボディ (JSON文字列、フラグと併用可) |
| `--name` | string | — | 品目名 (必須) |
| `--tax-type` | string | — | 税区分 (USE_SENDER/STANDARD_TAX_10/REDUCED_TAX_8/EXEMPTED_TAX等) |
| `--unit-name` | string | — | 単位名 |
| `--unit-price` | float64 | `0` | 単価 |

### `misoca item get`

品目の詳細を取得

```
misoca item get <id>
```

### `misoca item list`

品目の一覧を取得

| flag | type | default | description |
| --- | --- | --- | --- |
| `--page` | int | `0` | ページ番号 |
| `--per-page` | int | `0` | 1ページあたりの件数 (最大100) |

## `misoca user`

ユーザー情報の操作

### `misoca user me`

認証ユーザーの情報を取得
