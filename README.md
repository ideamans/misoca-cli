# misoca

[Misoca API v3](https://doc.misoca.jp/) を操作するコマンドラインツールです。

## インストール

```bash
go install github.com/miyanaga/misoca-cli/cmd/main.go@latest
# バイナリ名を misoca にリネーム
mv $(go env GOPATH)/bin/main $(go env GOPATH)/bin/misoca
```

または、リポジトリをクローンしてビルド:

```bash
git clone https://github.com/miyanaga/misoca-cli.git
cd misoca-cli
go build -o misoca ./cmd/main.go
cp misoca /usr/local/bin/
```

## 初期設定（認証）

### 環境変数にクレデンシャルがない場合

```bash
misoca auth
```

対話形式でセットアップが始まります:

1. ブラウザでMisocaの開発者ページが開きます
2. アプリケーションを作成（コールバックURL `http://localhost:18080/callback` は自動でクリップボードにコピーされます）
3. アプリケーションIDとシークレットを貼り付け
4. ブラウザで認証を許可

### 環境変数にクレデンシャルがある場合

```bash
export MISOCA_CLIENT_ID='your-client-id'
export MISOCA_CLIENT_SECRET='your-client-secret'
misoca auth
```

即座にブラウザ認証が開始します。

トークンは `~/.config/misoca-cli/token.json` に保存され、以降は自動的にリフレッシュされます。

## コマンド一覧

### 請求書 (`invoice`)

```bash
misoca invoice list                          # 一覧（デフォルト: 下書き）
misoca invoice list --type untrashed         # ゴミ箱以外の全請求書
misoca invoice list --condition "キーワード"   # キーワード検索
misoca invoice list --payment-status paid    # 入金済みのみ
misoca invoice list --from 2025/01/01 --to 2025/12/31  # 請求日で絞り込み
misoca invoice get <id>                      # 詳細取得
misoca invoice create --json '{...}'         # 作成
misoca invoice pdf <id> -o invoice.pdf       # PDF取得
misoca invoice trash <id>                    # ゴミ箱に移動
misoca invoice restore <id>                  # ゴミ箱から復元
misoca invoice submit <id>                   # 送付済みにする
misoca invoice unsubmit <id>                 # 送付済みを取り消し
misoca invoice pay <id>                      # 入金済みにする
misoca invoice unpay <id>                    # 入金済みを取り消し
misoca invoice mail <id>                     # 郵送依頼
```

#### list フィルタオプション

| フラグ | 説明 |
|--------|------|
| `--type` | `active`(デフォルト) / `archived` / `trashed` / `untrashed` |
| `--condition` | 請求書番号・取引先名・件名・社内メモのキーワード検索 |
| `--contact-group-id` | 取引先グループID |
| `--payment-status` | `paid` / `unpaid` |
| `--invoice-status` | `submitted` / `unsubmitted` |
| `--from` / `--to` | 請求日の範囲 (YYYY/MM/DD) |
| `--due-date-from` / `--due-date-to` | 支払期限の範囲 |
| `--updated-from` / `--updated-to` | 更新日の範囲 |
| `--order` | `asc` / `desc` (デフォルト: desc) |
| `--order-by` | `created_at` / `updated_at` / `issue_date` / `payment_due_on` |
| `--page` / `--per-page` | ページネーション (最大100件/ページ) |

### 見積書 (`estimate`)

```bash
misoca estimate list                         # 一覧
misoca estimate list --type untrashed        # ゴミ箱以外
misoca estimate get <id>                     # 詳細取得
misoca estimate create --json '{...}'        # 作成
misoca estimate pdf <id> -o estimate.pdf     # PDF取得
misoca estimate logo <id> -o logo.png        # ロゴ画像取得
misoca estimate stamp <id> -o stamp.png      # 印影画像取得
misoca estimate distribute <id> --json '{...}'  # メール送信
```

### 納品書 (`delivery-slip`, エイリアス: `ds`)

```bash
misoca ds list                               # 一覧
misoca ds list --type untrashed              # ゴミ箱以外
misoca ds get <id>                           # 詳細取得
misoca ds create --json '{...}'              # 作成
misoca ds pdf <id> -o slip.pdf               # PDF取得
```

### 送り先 (`contact`)

```bash
misoca contact list                          # 一覧
misoca contact list --contact-group-id 123   # 取引先グループで絞り込み
misoca contact list --trashed                # 非表示の送り先を取得
misoca contact get <id>                      # 詳細取得
misoca contact create --json '{...}'         # 作成
misoca contact hide <id>                     # 非表示にする
misoca contact restore <id>                  # 再表示する
```

### 取引先グループ (`contact-group`, エイリアス: `cg`)

```bash
misoca cg list                               # 一覧
misoca cg list --trashed                     # 非表示を含む
misoca cg get <id>                           # 詳細取得
misoca cg create --json '{...}'              # 作成
misoca cg hide <id>                          # 非表示にする
misoca cg restore <id>                       # 再表示する
```

### 品目 (`item`)

```bash
misoca item list                             # 一覧
misoca item get <id>                         # 詳細取得
misoca item create --json '{...}'            # 作成
```

### ユーザー (`user`)

```bash
misoca user me                               # 認証ユーザーの情報を取得
```

## 使用例

### 請求書を作成してPDFをダウンロード

```bash
# 請求書を作成
misoca invoice create --json '{
  "contact_id": 1234567,
  "subject": "開発費用",
  "issue_date": "2025/03/18",
  "payment_due_on": "2025/04/30",
  "items": [
    {
      "name": "Webアプリケーション開発",
      "unit_price": 500000,
      "quantity": 1,
      "unit_name": "式",
      "tax_type": "STANDARD_TAX_10"
    }
  ]
}'

# PDFをダウンロードして開く
misoca invoice pdf <id> -o invoice.pdf && open invoice.pdf
```

### キーワードで請求書を検索

```bash
misoca invoice list --type untrashed --condition "プロジェクト名" --per-page 100
```

## 技術スタック

- Go
- [Cobra](https://github.com/spf13/cobra) (CLI フレームワーク)
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) (OAuth2 認証)

## ライセンス

MIT
