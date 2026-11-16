# AGENTS.md

Go CLI that updates an existing Hetzner DNS record (Hetzner Cloud API) with the current external IP.
Sister project: `github.com/Mario-F/cloudflare-dyndns` (same structure, Cloudflare instead of Hetzner). No over-engineering.

## Layout

- `main.go` → `cmd.Execute()`
- `cmd/` – Cobra commands: `root.go` (global flags `--debug`, `--quiet`, `--config`, `--version ipv4|ipv6`), `myip.go`, `records.go`, `update.go`, `version.go`, `setup.go` (interactive, unfinished)
- `internal/hetzner/` – Hetzner Cloud DNS client (`https://api.hetzner.cloud/v1/`, `net/http`)
- `internal/externalip`, `internal/providers`, `internal/network`, `internal/logger` – external IP detection, shared with cloudflare-dyndns
- `internal/setup/` – cron job helper for `setup`
- Tooling: `Makefile`, `Dockerfile`, `.goreleaser.yml`, `.github/workflows/{test,release,docker}.yml`, `install.sh`, `debug` (delve) + `.vscode/launch.json`

## Commands

```shell
make test     # go test -race with coverage
make lint     # golangci-lint + misspell
make godry    # goreleaser snapshot
docker run --rm -v "$PWD":/w -w /w goreleaser/goreleaser:latest check   # validate .goreleaser.yml
```

## Toolchain pitfalls

- If Go auto-downloads a newer toolchain (`GOTOOLCHAIN=auto`), `go test -cover…` can fail with `go: no such tool "covdata"`. Use a full SDK (`go install golang.org/dl/goX.Y.Z@latest && goX.Y.Z download`, `~/sdk/goX.Y.Z/bin` first in `PATH`, `GOTOOLCHAIN=local`).
- A golangci-lint built with an older Go than the code fails with type-check errors; run the current one via `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...`.

## Behaviour

- External IP: shuffled provider list per IP version; an IP is accepted once two providers confirm it; values are validated by regex. Keep provider logic and regexes in sync with cloudflare-dyndns.
- Record reference: `<zone_id>/<rr_name>/<rr_type>` as printed by `records`.
- Update: `POST zones/{zone}/rrsets/{name}/{type}/actions/set_records`; same IP → "No IP diff detected", no write.
- Token check: `GET zones?page=1&per_page=1`. Requests: 20 s timeout, up to 3 attempts on network errors, 429 and 5xx.
- The token must never be logged, also not with `--debug`.

## Tests

- Unit tests only (`internal/hetzner`, `internal/providers`). Never call the real Hetzner API or real infrastructure from tests; use `httptest.Server` for new client tests.

## Privacy and secrets

- No real infrastructure in the repo: no own domains, hostnames, IPs, zone/record IDs or tokens – not in README, tests, examples or commit messages.
- Use reserved placeholders only: `example.com` / `example.net` / `example.org`, IPv4 `192.0.2.0/24`, `198.51.100.0/24`, `203.0.113.0/24`, IPv6 `2001:db8::/32`, token `your-api-token`.
- Before committing, grep the diff for anything that looks like real data. Once pushed, removing it needs a history rewrite plus GitHub-side cleanup (dangling commits stay reachable by SHA).
- Do not commit local debug binaries (`__debug_bin*`).

## Workflow

- Branch `master`. Commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`, …).
- Commit, push, tags and history rewrites only with explicit approval of the maintainer.
- Releases: tag → goreleaser (uses secret `GH_TOKEN`, so the published release triggers the multi-arch image `ghcr.io/mario-f/hetzner-dyndns`).
