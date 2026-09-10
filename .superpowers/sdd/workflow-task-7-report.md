# Task 7 Report: Preserve Workflow Artifacts On Failure

Status: complete

Implementation commit: `b869095` (`ci: retain acceptance artifacts after failures`)

Changes:

- Made `Tar provider binaries` and `Upload artifacts` unconditional with `always()`.
- Made archiving conditional on an executable provider binary.
- Added a fixed `provider binary unavailable` status marker when the binary is absent.
- Configured artifact upload to warn when no files exist and upload the archive/status paths without directory or secret output.
- Added workflow contract assertions covering these guarantees.

Verification:

- `go test ./provider -run TestOwnedWorkflow -count=1` — PASS
- `go test ./provider -run 'TestOwnedWorkflow|TestRegistryMetadata' -count=1` — PASS
- `go test ./...` — PASS
- `gofmt -w provider/registry_metadata_test.go` — PASS
- `git diff --check` — PASS

Concerns: live acceptance execution was not run locally because it requires protected Dokploy credentials. The worktree retains pre-existing unrelated changes (`.superpowers/sdd/task-1-report.md` deletion and `.superpowers/sdd-tools/` untracked content).
