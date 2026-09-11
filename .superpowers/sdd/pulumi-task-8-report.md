# Task 8 report: strengthen Pulumi lifecycle evidence

## Status

Implemented and verified in the `live-acceptance-coverage-hardening` worktree.

## Changes

- Added aggregate lifecycle validators and tests for expected mutating phases,
  intentional reads/no-ops, and replacement/deletion rejection.
- Captured and asserted Automation API preview and update summaries for both live
  revisions while retaining stable output/ID assertions.
- Added supported export-state filtering after destroy and `ListStacks` absence
  evidence after workspace removal.
- Kept cleanup bounded and independent, with diagnostic sanitization unchanged.

## Verification

- `go test ./tests -run 'TestLifecycleAggregate|TestLifecycleSmokeCleanup' -count=1`
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
aggregate operation categories and exported deployment resource shape; no live
credentials or Pulumi CLI were available for execution in this environment.

## Quality finding follow-up

- Aggregate `ChangeSummary` is now treated honestly: it cannot identify custom
  resources, so validation uses the strongest supported proxy (at least one
  phase-appropriate create/update), allows intentional `read`/`same`/`noop`,
  and rejects replace/delete/unsupported mutating operations. Exact four-custom
  resource coverage is proved by mock resource capture and stable live IDs,
  not by aggregate counts. This limitation is documented in the validator.
- Preview and update adapter tests now cover deterministic conversion, allowed
  reads/no-ops, sorted diagnostics, and missing update summaries.
- Cleanup orchestration now performs export validation after destroy, then runs
  export validation in its own bounded context even when destroy fails, then
  runs `RemoveStack` and `ListStacks` in further independent contexts. All four
  errors remain separately observable.

Follow-up verification:

- `go test ./tests -run 'TestLifecycleAggregate|TestLifecycleSmokeProgram|TestLifecycleSmokeCleanup' -count=1` — PASS
- `go test ./tests -run TestAccLifecycleSmoke -count=1 -v` — PASS; skipped without live acceptance configuration

## Review correction

The aggregate-operation design was corrected to avoid attributing provider or
default-resource counts to custom resources. The validator now checks only
phase-appropriate aggregate activity and forbidden operations; mock capture
asserts the exact four custom resources and stable IDs. Destroy, export, remove,
and list failures are returned independently, and export is attempted with its
own bounded context even after destroy fails.
