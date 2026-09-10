# Task 3 Report: Hybrid Mount Lifecycle Coverage

## Status

Complete. Production provider behavior was not changed.

## TDD evidence

- **RED:** `go test ./provider -run 'TestMountDiffCartesianMatrix|TestMountUpdateBodyAndRedeployMatrix' -count=1` failed after the new matrix was added with the intentionally incorrect `mountPath` expectation (`expected update&replace, actual update`).
- **GREEN:** Restoring the required mutable expectation and implementing the lifecycle extraction made the same command pass.
- **Focused GREEN:** `go test ./provider -run 'TestMountDiffCartesianMatrix|TestMountUpdateBodyAndRedeployMatrix|TestLiveTier2Workloads/Mount' -count=1 -v` passed; deterministic matrices passed and live acceptance skipped because opt-in credentials were absent.

## Changes

- Added deterministic 6-target × 3-type mutable and replacement diff coverage.
- Extracted `runLiveMountLifecycle` with immediate disarmable cleanup ownership, create/read/update/post-update read/ID-only read/diff/delete/absence checks.
- Upgraded PostgreSQL Mount dispatch to the complete lifecycle; MySQL, MariaDB, Redis, and Compose dispatch behavior remains serial dispatch-only as required.
- Preserved the existing one-heavy-operation lease and cleanup ordering.

## Verification

- `go test ./provider -run 'TestHeavyOperation|TestMountTargetDispatch' -count=1` — PASS.
- `go test ./...` — PASS.
- `git diff --check` — PASS.
- No `t.Parallel()` occurrences in `provider/live*.go`.

## Files

- `provider/mount_lifecycle_matrix_test.go`
- `provider/live_workloads_test.go`
- `.superpowers/sdd/mount-task-3-report.md`

## Commits

- `8635260` — `test: complete hybrid mount lifecycle coverage`
- Report commit: recorded after this report was added.

## Concerns

- Live lifecycle verification was not executed against Dokploy because the required explicit acceptance opt-in and credentials were unavailable; the focused command recorded the expected skip.
- The pre-existing deletion of `.superpowers/sdd/task-1-report.md` was not staged or modified by this task.
