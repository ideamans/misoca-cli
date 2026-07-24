---
name: misoca-install
description: Make the misoca command available, installing it only if it is missing. Use when another skill reports that `misoca` is not on PATH, or when the user asks to install, update or upgrade the ideamans Misoca CLI. Prefers an already-installed binary, then the latest GitHub release, then a build from source with go install.
license: MIT
compatibility: Requires curl (or wget) and tar to install from a release, or a Go toolchain for the source fallback. Standalone — does not need misoca to be present already. Installs from the public repository github.com/ideamans/misoca-cli, so no GitHub authentication is needed.
allowed-tools: Bash(curl:*) Bash(wget:*) Bash(tar:*) Bash(unzip:*) Bash(go:*) Bash(uname:*) Bash(command:*) Bash(which:*) Bash(mkdir:*) Bash(mv:*) Bash(cp:*) Bash(rm:*) Bash(chmod:*) Bash(ls:*) Bash(test:*) Bash(echo:*) Read
---

# misoca-install

Make the `misoca` command usable, doing the least work that achieves it.

## Route 1 — an existing installation on PATH

```bash
command -v misoca && misoca --version
```

If that resolves, **use it and stop here.** Do not check for a newer release —
it costs an API call and the user did not ask for an upgrade.

Two checks before trusting the hit:

- **It is the right tool.** `misoca llm | head -1` must read
  `# misoca CLI — AIエージェント向けリファレンス`.
- **It is recent enough.** If `misoca llm` is not a known command, the binary
  predates the embedded reference. Say so and continue to route 2.

Continue past this section only when the command is missing, is the wrong tool,
is too old, or the user explicitly asked to update.

## Route 2 — the latest GitHub release

```bash
VERSION=$(curl -fsSL https://api.github.com/repos/ideamans/misoca-cli/releases/latest \
  | grep '"tag_name"' | head -1 | cut -d'"' -f4)   # e.g. v0.2.0

OS=$(uname -s | tr '[:upper:]' '[:lower:]')            # darwin | linux
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64  # amd64 | arm64
curl -fsSL -o /tmp/misoca.tar.gz \
  "https://github.com/ideamans/misoca-cli/releases/download/${VERSION}/misoca-cli_${VERSION#v}_${OS}_${ARCH}.tar.gz"
```

**The archive is named `misoca-cli_…` but the binary inside is `misoca`** —
without the `-cli` suffix. Windows ships a `.zip`.

If the download 404s, list the actual assets on the release page rather than
retrying variations.

### Install onto PATH

```bash
tar -xzf /tmp/misoca.tar.gz -C /tmp
mkdir -p ~/.local/bin && mv /tmp/misoca ~/.local/bin/ && chmod +x ~/.local/bin/misoca
```

Prefer the first writable directory already on PATH — `~/.local/bin`, then
`/usr/local/bin`. Two things not to do on your own initiative:

- If nothing on PATH is writable, leave the binary in `/tmp`, print the exact
  `sudo mv` command and let the user run it. Do not run `sudo` yourself.
- If `~/.local/bin` is not on PATH, give the user the line to add to their shell
  profile. Do not edit the profile for them.

## Route 3 — build from source

Needs a Go toolchain and compiles rather than downloads. Note the `/cmd`
suffix — the main package is not at the module root.

```bash
go install github.com/ideamans/misoca-cli/cmd@latest
```

`go install` names the binary after the last path element, so this produces
**`cmd`**, not `misoca`. Rename or symlink it:

```bash
mv "$(go env GOPATH)/bin/cmd" "$(go env GOPATH)/bin/misoca"
```

Say so explicitly — a binary called `cmd` on PATH is confusing and easy to
overwrite from another project.

## Verify

```bash
misoca --version
misoca llm | head -5
misoca user me
```

Report which route was taken, the version and the install path.

`misoca user me` failing means OAuth has not been completed. **Tell the user to
run `misoca auth` themselves** — it opens a browser and waits for their
approval, so it cannot be done for them. They will need a Misoca application's
client id and secret (or the interactive flow will walk them through creating
one), and the callback URL `http://localhost:18080/callback`.
