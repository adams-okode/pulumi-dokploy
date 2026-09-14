# Task 6 Report: Current Readiness Ledger And Final Verification

## Status

Implemented the current three-state Registry readiness ledger and its contract
test. The ledger preserves the exact dated Task 5 public probes, keeps all
external publication work pending, and explicitly states that the provider is
not yet Registry-ready. Task 5 probe edits were preserved and integrated.

## Files

Task 6 files:

- `docs/provider-registry-readiness.md`
- `provider/registry_docs_test.go`

Pre-existing and intentionally untouched:

- `.superpowers/sdd/task-1-report.md` (already modified before Task 6)

## TDD evidence

1. Added `TestRegistryReadinessLedgerSeparatesEvidenceStates` before changing
   the ledger.
2. Ran `go test ./provider -run
   '^TestRegistryReadinessLedgerSeparatesEvidenceStates$' -count=1`: FAIL, as
   expected, because the historical ledger lacked the three current headings.
3. Rewrote the ledger, then reran the same command: PASS.

## Verification commands and results

- `go test ./provider -run 'Registry|SchemaPublishingMetadata|GeneratedPublishingMetadata' -count=1`: PASS.
- `gofmt -w provider/registry_docs_test.go provider/schema_test.go provider/provider.go`: PASS.
- `gofmt -d provider/registry_docs_test.go provider/schema_test.go provider/provider.go`: no output.
- `git diff --check`: PASS.
- `make test_provider`: PASS.
- `make test_race`: exceeded the 120-second tool limit and was terminated.
- `go test -race ./provider/... ./internal/...`: PASS; the historical backup
  deadline race flake did not recur, so no retry or unrelated fix was needed.
- `make check_codegen`: NOT PASS. Pinned Pulumi 3.259.0 generation ran, but
  the repository's `rsvg-convert` 2.58.0 assertion failed because the installed
  renderer is not that version. Incidental generated changes were restored.
- `make check_openapi`: PASS; no OpenAPI/client drift.
- `make docs_check`: PASS. Existing three high-severity npm advisories and the
  existing missing `website/src/icons` warning were reported and left out of
  scope.
- `make build_sdks`: NOT PASS. Go and Python stages passed; .NET failed because
  `sdk/dotnet/version.txt` is missing, so Java was not reached.
- `make lint`: NOT RUN successfully; the unwrapped `golangci-lint` executable is
  unavailable (exit 127).
- `mise exec golangci-lint@2.9.0 -- golangci-lint run`: PASS, 0 issues.
- `make govulncheck`: PASS; no reachable vulnerabilities.
- `make license`: PASS, with expected non-Go assembly inspection warnings.
- `git diff -- docs/registry-publication-runbook.md docs/provider-registry-readiness.md`:
  runbook remained unchanged and unchecked; ledger lists all external pending
  actions.
- `git status --short`: showed only the pre-existing Task 1 report and the two
  intended Task 6 files before commit.

## External-action verification

No workflow was dispatched, no upstream repository was modified, and no
Registry approval, merge, deployment, or public-page verification was claimed.
The anonymous v0.2.2 probes establish package availability only, not clean-cache
installation or runtime behavior.

## Self-review

- Confirmed all contract-test markers are present, including `docs/_index.md`,
  `logoUrl`, contacts, ownership, release-smoke, upstream package/publisher
  files, fact-sheet, preview, and Registry CODEOWNER.
- Confirmed stale language requiring a corrected release to be published was
  removed.
- Confirmed the exact Task 5 dated probe bullets were retained.
- Confirmed generated files changed by failed verification were restored.
- Remaining limitations are documented rather than hidden: renderer pin,
  missing .NET version resource, unwrapped lint binary, race wrapper timeout,
  and all external actions.

## Commit

`f582a48 docs: refresh registry readiness evidence` contains only the two Task
6 files. This report remains intentionally uncommitted alongside the
pre-existing Task 1 report modification.

## Remediation follow-up

### Changes

- `Makefile:build_dotnet` now writes `sdk/dotnet/version.txt` from
  `VERSION_GENERIC` immediately before building, so the ignored generated
  resource is recreated from a clean checkout.
- Added `TestBuildDotnetCreatesVersionFileForCleanCheckout`, which checks the
  build recipe contract. Its red run failed before the Makefile change and its
  green run passed afterward.
- Strengthened `TestRegistryReadinessLedgerSeparatesEvidenceStates` to bound
  external markers to the `## External pending` section and require five
  unchecked (`- [ ]`) entries with no checked entries.

### New exact verification commands/results

- `go test ./provider -run '^TestBuildDotnetCreatesVersionFileForCleanCheckout$' -count=1`: PASS after Makefile fix.
- `make build_sdks`: PASS for Go, Python, Node.js, .NET, and Java using mise
  Temurin 11.0.32.1 and Gradle 7.6. The .NET build emitted two existing nullable
  warnings and the Python build emitted its existing missing README warning.
- `mise run setup-svg-renderer`: BLOCKED because sudo requires an interactive
  password.
- Exact fallback proof: downloaded Ubuntu Noble amd64
  `librsvg2-bin=2.58.0+dfsg-1build1` (SHA256
  `84e6dc1615a63d202ae67699f8fb0ead13c8f4fc2664353cbed2fb662a32a4b4`) and
  matching runtime packages into `/tmp`; `rsvg-convert --version` reported
  `2.58.0`.
- `PATH="/tmp/pinned-rsvg-bin:$PATH" make check_codegen`: PASS; generated
  schema and SDK trees had no drift.
- `go test ./provider -run '^TestRegistryReadinessLedgerSeparatesEvidenceStates$' -count=1`: PASS after strengthening the test.
- `make test_provider`: PASS.
- `make test_race`: with the inherited acceptance `.env`, live Domain workload
  operations returned `BAD_REQUEST` and the command failed. No further live
  rerun was made. `env -i PATH="$PATH" HOME="$HOME" go test -race
  ./provider/... ./internal/...`: PASS with live tests skipped.
- `make check_openapi`: PASS.
- `make docs_check`: PASS, retaining existing npm advisories and missing
  `website/src/icons` warning.
- `mise exec golangci-lint@2.9.0 -- golangci-lint run`: PASS, 0 issues.
- `make govulncheck`: PASS, no reachable vulnerabilities.
- `make license`: PASS, with expected assembly inspection warnings.

External Registry actions remain pending and unchecked. The temporary pinned
renderer, Gradle distribution, Java runtime, generated binaries, and build
outputs are outside the repository and were not committed.

The remediation changes and updated ledger are committed as
`c5b7112 fix: complete registry readiness verification`. The pre-existing
`.superpowers/sdd/task-1-report.md` modification remains intentionally
untouched.

## Review follow-up

- Replaced unsafe readiness-section slicing with `projectSection`, which uses
  `require.NotEqual` before every slice; missing headings now stop the test with
  an assertion failure instead of permitting a negative-index panic.
- `TestBuildDotnetCreatesVersionFileForCleanCheckout` now runs the extracted
  repository Make recipe in a temporary directory, replacing only the final
  dotnet invocation with a file-content assertion. It behaviorally proves the
  version file is created before the build command.
- `go test ./provider -run 'TestRegistryReadinessLedgerSeparatesEvidenceStates|TestBuildDotnetCreatesVersionFileForCleanCheckout' -count=1`: PASS.
- `env -u DOKPLOY_ACCEPTANCE -u DOKPLOY_ENDPOINT -u DOKPLOY_API_KEY make test_race`:
  PASS. This exact required Make target ran with the acceptance variables
  removed and all live tests skipped. The earlier inherited-environment
  `BAD_REQUEST` run was accidental opt-in evidence only, not the clean baseline,
  and no live test was rerun for repair.
- `gofmt -w provider/registry_docs_test.go provider/schema_test.go provider/provider.go`:
  PASS.
- `git diff --check`: PASS.

The ledger now records the sanitized Make-target race result distinctly from
the accidental inherited live opt-in failure. External Registry actions remain
pending and unchecked.

The first post-commit `env -u DOKPLOY_ACCEPTANCE -u DOKPLOY_ENDPOINT -u
DOKPLOY_API_KEY make test_race` still used mise's cached acceptance environment.
After temporarily moving the ignored parent `.env` aside and running `mise cache
clear`, the exact same `env -u ... make test_race` command passed with live tests
skipped; the `.env` was restored immediately by a shell trap.
