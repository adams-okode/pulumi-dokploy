# Task 8 report: strengthen Pulumi lifecycle evidence

## Status

Implemented and verified in the `live-acceptance-coverage-hardening` worktree.

## Changes

- Added pure lifecycle summary validators and tests for four revision-one creates,
  revision-two updates, and replacement/deletion rejection.
- Captured and asserted Automation API preview and update summaries for both live
  revisions while retaining stable output/ID assertions.
- Added supported export-state filtering after destroy and `ListStacks` absence
  evidence after workspace removal.
- Kept cleanup bounded and independent, with diagnostic sanitization unchanged.

## Verification

- `go test ./tests -run 'TestLifecycleSummary|TestLifecycleSmokeCleanup' -count=1`
  — initially failed to compile before implementation because validators were
  undefined (TDD red phase).
- `go test ./tests -run 'TestLifecycleSummary|TestLifecycleSmokeProgram|TestLifecycleSmokeCleanup' -count=1`
  — PASS.
- `go test ./tests -run TestAccLifecycleSmoke -count=1 -v` — PASS; live test
  skipped because `DOKPLOY_ACCEPTANCE=1` was not set.
- `go test ./tests -count=1` — PASS.
- `go test ./...` — PASS.
- `git diff --check` — PASS.

## Concerns

Live assertions require the pinned Pulumi Automation API to report the expected
resource operation summaries and exported deployment resource shape; no live
credentials or Pulumi CLI were available for execution in this environment.

## Quality finding follow-up

- Revision summaries now require the exact stable custom-resource counts (four
  creates for revision one and three updates for revision two), reject extra
  creates and unsupported operations, and allow only `same`/`noop` ambient
  reports in addition. Replacement and deletion diagnostics remain explicit.
- Preview and update adapter tests now cover deterministic conversion, allowed
  no-op reports, and missing update summaries.
- Cleanup orchestration now performs export validation after destroy, then runs
  `RemoveStack` and `ListStacks` in independent contexts regardless of
  post-destroy validation or earlier cleanup errors. Errors remain separately
  observable.

Follow-up verification:

- `go test ./tests -run 'TestLifecycleSummary|TestLifecycleSmokeProgram|TestLifecycleSmokeCleanup' -count=1` — PASS
- `go test ./tests -run TestAccLifecycleSmoke -count=1 -v` — PASS; skipped without live acceptance configuration
