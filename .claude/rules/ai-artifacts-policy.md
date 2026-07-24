# 生成物 — 手編集しないこと

| 生成物 | 原本 |
| --- | --- |
| `internal/llmdocs/90-commands.md` | `internal/cmd/*.go` の cobra コマンド定義（`internal/gen-llmdocs` が描画） |

手書き（編集してよい）:

- `internal/llmdocs/00-guide.md` — 認証・出力・一覧の絞り込み・注意が必要な操作
- `internal/llmdocs/10-documents.md` — 帳票の種類と状態
- `plugins/misoca-cli/skills/*/SKILL.md`
- `context7.json`

生成物を直接編集しても次の `go generate ./...` で消え、それまでは CI が古い差分で
落ちる。カタログを良くしたいときはコマンドの `Short` / `Long` / フラグ説明を直すこと。

再生成は `/regen-ai`、または `go generate ./... && go test ./...`。
