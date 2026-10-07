# AGENTS.md

Fork of `glanceapp/glance` (personal dashboard, Go monolith). This fork collects bug fixes and features missing upstream.

## Hard rule: repository boundaries

- **All writes go to `felix-schindler/glance` only**: pushes, branches, PRs, comments, releases.
- **Never interact with `glanceapp/glance`**: no PRs, no comments, no reviews, no pushes there — read-only `git fetch` / `gh pr view` is allowed only to port a specific PR on explicit request.
- `gh` CLI defaults: always pass `--repo felix-schindler/glance` (or `-R`) so nothing ever targets upstream by accident. `upstream` git remote is fetch-only in practice — never push to it.

## Build / test / run

```sh
go build ./...        # build
go vet ./...          # vet (run before committing)
go test ./...         # full test suite
go test ./internal/glance/ -run 'TestName' -v   # single test
gofmt -l internal/ pkg/ main.go                 # must output nothing new (address_test.go + widget-search.go have pre-existing findings, leave them)
./glance --config glance.yml                    # run; config template in docs/glance.yml
./glance config:validate --config glance.yml
```

Go 1.27. Module path stays `github.com/glanceapp/glance` (do not rename — keeps ports clean).

## Layout

- `main.go` → `internal/glance/` — almost everything lives here (one package).
- `internal/glance/templates/` + `static/` — embedded via `embed.go`; CSS bundled at runtime, no build step. JS is vanilla + `static/js/templating.js` helper.
- `docs/configuration.md` — user-facing config reference; update it with every config-affecting change.
- `docs/glance.yml` — example config; keep in sync with new options.
- Tests are colocated `*_test.go`; run the suite for every change.

## Config changes

Validation lives in `isConfigStateValid` (`internal/glance/config.go`). Secrets use `${secret:name}` (Docker secrets) / `${VAR}` (env) — see `parseConfigVariables`. Auth: local users hash via `password:hash`, session is an HMAC cookie (`auth.go`); OIDC issues the same cookie (`auth_oidc.go`).

## Porting upstream PRs (e.g. Kogoro:feature/oidc-auth)

Upstream PRs target upstream `dev` from old bases — never merge them directly:

1. `git fetch <fork-url> <branch>:pr-<name>` (no new remotes needed) and inspect with `git diff --stat $(git merge-base HEAD pr-<name>) pr-<name>`.
2. Cherry-pick/port only the feature files onto a `feature/<name>` branch from our `main`; exclude unrelated upstream-`dev` drift.
3. Resolve `go.mod`/`go.sum` by taking our versions, then `go get <new-deps> && go mod tidy`. Preserve our hardening commits (auth body limit, server timeouts).
4. Verify: `go build`, `go vet`, `go test ./...`, plus a live smoke test for auth/network-facing changes (throwaway server under `$TMPDIR`, never in repo).
5. Push branch to `origin`, open PR with `gh pr create --repo felix-schindler/glance --base main`, documenting review findings + verification.

## Conventions

- Smallest diff that works; no new deps without need; match existing code style.
- Commit messages: short imperative subject, body explaining why.
- Never commit binaries, test configs, or harness scripts — keep those in `$TMPDIR`.
