---
paths:
  - "internal/cmd/*.go"
  - "internal/llmdocs/0*.md"
  - "internal/llmdocs/1*.md"
---

# 埋め込みリファレンスの原本に触れました

コマンド・フラグ・ヘルプ文字列を変えたなら、終わる前に `/regen-ai` を実行して
`internal/llmdocs/90-commands.md` を一致させること。CI が再生成して差分が出たら
落ちます。

一覧の既定の絞り込みや帳票ごとの操作の差を変えたなら、`00-guide.md` /
`10-documents.md` を手で直すこと。ここがズレると、エージェントは
**エラーにならずに一部のデータだけを見て結論を出します**。

`90-commands.md` は直接編集しないこと。
