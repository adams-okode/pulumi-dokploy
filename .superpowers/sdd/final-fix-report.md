# Final broad-review fix report

Implemented on `fix/live-acceptance-findings` without live calls or `.env` use.

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

- `go test ./provider -run 'TestSSHKeyCreate(Recovers|ReturnsOriginal|ReturnsAmbiguous)' -count=1` — PASS.
- `go test ./provider -run 'TestWorkloadCallPathsEmitOnlyStructuralDiagnostics' -count=1` — PASS.
- `go test ./provider -count=1` — PASS.
- `go test ./... -count=1` — PASS.
- `go test -race ./provider/... ./internal/... -count=1` — PASS.
- `make docs_check` — PASS (Astro check: 0 errors/warnings/hints; 44 docs tests and 2 built-site tests passed).
- `git diff --check` — PASS.

## Earlier concurrent main review evidence

- Backup snapshots retain every non-empty `backupId`, mark malformed observations
  invalid, and never adopt them on a later complete observation.
- Target-backup discovery failures use the fixed operation-level message
  `backup.create could not read target backups`; cancellation and deadline
  wrapping remains intact for `errors.Is`.
- The focused backup regression suite, short provider/internal suite, provider/
  internal race suite, and whitespace check passed. Live Dokploy acceptance was
  not run for those backup fixes.

## RED

- Added focused structural-diagnostic coverage for live failures.
- The workflow contract initially failed because the new protected variables
  were absent from the Tier 2 contract and YAML indentation was invalid.

## GREEN

- `go test ./provider -run 'Test(RequireNoErrorUsesStructuralDiagnosticOnly|SuccessfulCreateRegistersCleanupAndExplicitDeleteReleasesIt|ExplicitDeleteRetainsOwnershipWhenAbsenceVerificationFails|OwnedWorkflow)$' -count=1`
  passed.
- `go test ./provider ./tests -count=1` passed.
- `git diff --check` passed.

The live control-plane delete helpers now retain fallback ownership until
verified absence and release it exactly once after successful explicit cleanup.
Generic live error assertions now emit structural resource/operation text only.
Tier 1 and Tier 2 receive the configured server scope; Tier 2 receives all
GitLab variables through protected workflow secrets. Documentation and the
workflow contract list these variables without values.

## Verification concern

`go test -race ./provider ./tests -count=1` encountered the pre-existing
timing-sensitive `TestBackupCreateCancellationErrorOmitsTargetID`, which left
one scripted polling request when cancellation happened before the next poll.
The race command otherwise completed the `tests` package; no race report was
emitted.

## Final hardening pass

### RED

- `TestLiveDiagnosticSourceContract` initially found raw Automation API error
  concatenation, output-map formatting, and ID-bearing live assertions.
- The project cleanup-owner regression test initially had no disarmable owner
  boundary; coverage was added for retention on failure and disarm on verified
  success.

### GREEN

- Automation API failures and output assertions now use field-only structural
  diagnostics. Destroy-state and stack-removal validators omit live values.
- Live ID presence/equality assertions use operand-free helpers.
- Database passwords and related environment values are registered before live
  creates.
- `liveProject` now uses a disarmable owner and releases it only after verified
  absence.
- Added the static live-source diagnostic contract test.

Verification:

- `go test ./provider -count=1 -timeout 3m` — PASS
- `go test ./tests -count=1 -timeout 3m` — PASS
- `go test -race ./provider -skip 'BackupCreate' -count=1 -timeout 3m` — PASS
- `go test -race ./tests -count=1 -timeout 3m` — PASS
- `go test ./... -count=1 -timeout 3m` — PASS
- `git diff --check` — PASS

The complete provider race run remains timing-sensitive in existing
`TestBackupCreate_Cancellation`/`TestBackupCreateCancellationErrorOmitsTargetID`
polling tests; it failed due to an expected scripted request remaining and
reported no race detector failure. The complete run excluding those existing
backup cancellation tests passed.

## Final verification

- `go test ./provider -run 'Test(ApplicationSource|ComposeSource|Mount|Domain|Destination|Registry|LiveGate|ClassifyLiveServerHealth|OwnedWorkflow)' -count=1` — PASS.
- `go test ./tests -run 'TestLifecycle|TestPulumiCLI' -count=1` — PASS.
- `go test -short -count=1 ./provider/... ./internal/... ./tests/...` — PASS.
- `go test -race ./provider/... ./internal/... -count=1` — PASS on rerun.
- `go test ./... -count=1` — PASS on rerun with a 300-second timeout.
- `go test ./provider -run '^TestBackupCreate(Cancellation|CancellationErrorOmitsTargetID)$' -count=20` — PASS.
- `git diff --check` — PASS.

`golangci-lint run` could not execute because `golangci-lint` is not installed
in the environment. Live acceptance was not run because explicit credentials
were unavailable. Existing unrelated working-tree changes were left untouched.
