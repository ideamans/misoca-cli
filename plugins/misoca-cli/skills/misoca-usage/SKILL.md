---
name: misoca-usage
description: Work with Misoca through the misoca CLI — list and create invoices (請求書), estimates (見積書) and delivery slips (納品書), download their PDFs, mark invoices as sent or paid, and manage contacts, contact groups and items. Use when the user asks about a Misoca invoice or quote, wants to issue or look up billing documents, or needs a PDF of one.
license: MIT
compatibility: Requires the `misoca` binary on PATH — run the misoca-install skill if it is missing. Needs a completed OAuth login (`misoca auth`, which opens a browser and must be run by the user). Operates on real billing documents for real customers.
allowed-tools: Bash(misoca:*) Bash(jq:*) Bash(command:*) Read Write
---

# misoca-usage

Drive Misoca through the `misoca` CLI. Commands have the shape
`misoca <resource> <method> [flags]` and print JSON.

## 1. Confirm the tool and the login

```bash
command -v misoca && misoca --version
misoca user me
```

Missing binary? Run the `misoca-install` skill.

`misoca user me` returning the account details means OAuth is set up. If it
fails, **do not run `misoca auth` yourself** — it opens a browser and waits for
the user. Ask them to run it; the token is then cached in
`~/.config/misoca-cli/token.json` and refreshed automatically.

## 2. Read the reference

```bash
misoca llm | head -80
misoca llm | grep -A 20 '請求書の状態'
```

Embedded in the binary (~570 lines), so it matches the installed version.

## 3. Listing: the two things that hide data

- **`--type` is not sent unless you pass it**, so you get the API default
  (`active`). If the user says an invoice is missing, try
  `--type untrashed` before concluding it does not exist.
- **There is no automatic pagination.** `--per-page` caps at 100, and the CLI
  does not follow pages for you. Advance `--page` yourself and say how many
  records your answer covers.

Two independent invoice states, easy to confuse:

| Question | Flag |
| --- | --- |
| Has it been sent? | `--invoice-status submitted` / `unsubmitted` |
| Has it been paid? | `--payment-status paid` / `unpaid` |

"Unpaid invoices" is `--payment-status unpaid`.

```bash
misoca invoice list --payment-status unpaid --type untrashed | jq '.[] | {id, subject, total_amount}'
```

## 4. Creating and fetching documents

```bash
misoca invoice create --json '{...}'
misoca invoice pdf <id> -o invoice.pdf
```

`--json` carries the body; individual flags override the same fields. Show the
user the body you are about to send — a created invoice is a real document
against a real contact, even as a draft.

## 5. Actions that cost money or reach the customer

**Get explicit agreement before these, and never run one speculatively:**

| Command | Effect |
| --- | --- |
| `misoca invoice mail <id>` | **Requests physical postal delivery. It costs money and cannot be undone.** |
| `misoca estimate distribute <id> --json` | Emails the estimate to the customer. |
| `invoice submit` / `unsubmit` | Changes the sent state. |
| `invoice pay` / `unpay` | Changes the paid state — an accounting fact. |
| `invoice trash` / `restore` | Moves to and from the trash (reversible). |
| `contact hide` / `restore` | Hides and restores a contact (reversible). |

## 6. Report

Give back the document id, the subject and amount, and the state you left it in.
When you listed, say which filters were applied — an answer computed from the
default `active` page is not "all invoices".

## Failure modes

| Symptom | Cause | Fix |
| --- | --- | --- |
| `command not found: misoca` | not installed | run the `misoca-install` skill |
| auth errors on every command | OAuth not completed or token revoked | ask the user to run `misoca auth` |
| an expected invoice is missing | default `type` filter | retry with `--type untrashed` |
| totals look short | no automatic pagination | advance `--page`, `--per-page` max 100 |
| a hidden contact does not appear | contacts use `--trashed`, not `--type` | pass `--trashed` |
