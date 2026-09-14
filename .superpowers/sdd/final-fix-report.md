# Final broad-review fix report

Implemented as one final-remediation wave without live calls or `.env` use.

## Changes

- Removed the floating-runner `sudo apt-get` SVG setup. Codegen now uses the
  repository-owned, SHA-256-guarded deterministic logo materializer and keeps
  the checked-in PNG byte-equivalent.
- Strengthened chooser validation for exactly six ordered languages, no HCL,
  and language-local installation/example content.
- Structurally checked every runbook external action as its own unchecked item,
  rejecting both lowercase and uppercase checked states.
- Hardened public security-reporting and credential-disclosure governance
  assertions, and guarded Makefile section extraction before slicing.

## Evidence

- Domain, Mount, and MountDispatch successful creates now register bounded
  fallback ownership immediately. Explicit deletion releases ownership before
  absence verification, preventing duplicate cleanup after verification errors.
- SSH key create performs one post-error list discovery. A single novel
  same-name key returns partial state and `initFailed`; no candidate returns the
  sanitized original create error; ambiguous discovery returns a sanitized
  ambiguity initialization error without an ID. Create is never retried.
- Workload call-path diagnostic tests use a scripted fake client and verify
  raw IDs, paths, SQL, and content do not appear in structural diagnostics.
- Tier 2 documentation records the final sanitized evidence and preserves the
  cleanup limitation and historical findings.

## Verification

- `go test ./provider -run '^TestRegistry' -count=1` — PASS.
- `go test ./provider -run 'Test(OwnedWorkflow|ReleaseSmokeWorkflow|Workflow|Registry)' -count=1` — PASS (workflow/action YAML contracts and registry checks).
- `go test -short -count=1 ./provider/... ./internal/...` — PASS.
- `go test ./provider/... ./internal/... -count=1` — PASS.
- `env -i HOME="$HOME" PATH="$PATH" GOPATH="$GOPATH" GOMODCACHE="$GOMODCACHE" go test -race ./provider/... ./internal/... -count=1` — PASS.
- `tmp=$(mktemp) && python3 scripts/generate-logo-png.py website/public/logo.svg "$tmp" && cmp -s "$tmp" sdk/dotnet/logo.png && rm -f "$tmp"` — PASS (byte-equivalent).
- `make check_codegen` — PASS (generated schema/SDK diff clean).
- `mise exec java@11 gradle@8.14.3 -- make build_sdks` — PASS (all SDK builds; Java toolchain supplied explicitly).
- `make docs_check` — PASS (Astro: 0 errors/warnings/hints; 44 docs tests and 2 built-site tests).
- `git diff --check` — PASS.

The first concurrent full race invocation overlapped the ordinary Go suite and
flaked in the existing cancellation test with one scripted request remaining;
the required sanitized race command was rerun alone and passed.

## Earlier concurrent main review evidence

- Backup snapshots retain every non-empty `backupId`, mark malformed observations
  invalid, and never adopt them on a later complete observation.
- Target-backup discovery failures use the fixed operation-level message
  `backup.create could not read target backups`; cancellation and deadline
  wrapping remains intact for `errors.Is`.
- The focused backup regression suite, short provider/internal suite, provider/
  internal race suite, and whitespace check passed. Live Dokploy acceptance was
  not run for those backup fixes.

## Concerns

- `make docs_check` reports three existing high-severity npm audit findings and
  a non-fatal astro-icon missing-directory warning during the static build.
- SDK builds retain existing .NET nullable and Python packaging warnings.
- The first bare `make build_sdks` attempt could not find `gradle`; the required
  command passed with pinned mise-provided Java 11 and Gradle 8.14.3.

## Final re-review remediation

- `.mise.toml` now declares `java = "11"` and `gradle = "8.14.3"`; the setup
  action enables mise installation, so normal repository setup supplies both.
- `scripts/generate-logo-png.py` now parses the canonical SVG XML and cubic/line
  geometry, rasterizes it with standard-library supersampling, and emits a
  deterministic PNG. It has no Git, HEAD, or pre-existing-PNG dependency.
- `go test ./provider -run 'Test(RepositorySetupDeclaresPortableJavaBuildTools|LogoRendererDerivesOutputFromSVG)' -count=1` — PASS.
- `mise exec -- make build_sdks` — PASS (normal project setup, no command-line
  tool injection).
- `make check_codegen` — PASS after the generated PNG was updated and staged;
  schema and SDK working-tree diff is clean.
- `go test ./provider/... ./internal/... -count=1` — PASS.
- `env -i HOME="$HOME" PATH="$PATH" GOPATH="$GOPATH" GOMODCACHE="$GOMODCACHE" go test -race ./provider/... ./internal/... -count=1` — PASS.
- `make docs_check` — PASS (44 docs tests and 2 built-site tests).
- `git diff --check` — PASS.
