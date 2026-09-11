# Task 9 static verification report

## Status

Static verification completed on 2026-09-11 in the requested worktree. No
Critical or Important review finding required a code change. The existing
worktree changes remain untouched; only this report was added by Task 9.

## Commands and results

- `gofmt -w provider/application_source_test.go provider/compose_source_test.go provider/mount_lifecycle_matrix_test.go provider/live_control_plane_test.go provider/live_workloads_test.go provider/live_databases_test.go provider/live_backups_test.go provider/live_harness_test.go provider/live_harness_unit_test.go provider/domain_test.go provider/destination_test.go provider/registry_test.go tests/acceptance_test.go tests/acceptance_program_test.go` — passed.
- `go test ./provider -run 'Test(ApplicationSource|ComposeSource|Mount|Domain|Destination|Registry|LiveGate|ClassifyLiveServerHealth|OwnedWorkflow)' -count=1` — passed (`ok`, 0.195s).
- `go test ./tests -run 'TestLifecycle|TestPulumiCLI' -count=1` — passed (`ok`, 0.031s); live smoke remained skipped without opt-in.
- `go test -short -count=1 ./provider/... ./internal/... ./tests/...` — passed; provider, internal/client, and tests packages passed, generated package had no tests.
- `go test -race ./provider/... ./internal/...` — first attempt reported a remaining scripted request in `TestBackupCreateCancellationErrorOmitsTargetID`; a subsequent exact rerun passed (`ok provider 4.505s`, internal packages passed). A bounded verbose rerun also passed.
- `golangci-lint run` — unavailable: `golangci-lint: command not found`.
- `mise exec -- golangci-lint run` — unavailable: `mise: command not found`.
- `git diff --check` — passed.

## Self-review

Reviewed the changed provider, acceptance, workflow, README, and supporting
spec/plan files since `7a2f99c` for:

- complete Application and Compose source field reconstruction;
- Mount hybrid coverage, including the supported target/type matrix;
- reverse-order cleanup ownership and dependent PostgreSQL Mount cleanup;
- health-stop classification, bounded probes, and false-positive exclusions;
- secret registration and sanitized diagnostics;
- README controlled language, technical accuracy, and `.env` boundary;
- serial workflow gates, final failure aggregation, and unconditional artifact
  handling; and
- Pulumi summary, stable-ID, destroy-export, and stack-removal assumptions.

The deterministic tests cover the reviewed contracts. No unbounded cleanup,
parallel live test, public certificate request, backup execution, credential
printing, or state-file manipulation was found. No provider behavior change was
introduced by the reviewed changes.

## Files and commits

The reviewed change set includes:

`.github/workflows/run-acceptance-tests.yml`, `CONTRIBUTING.md`,
`provider/application_source_test.go`, `provider/compose_source_test.go`,
`provider/destination_test.go`, `provider/domain_test.go`,
`provider/live_backups_test.go`, `provider/live_control_plane_test.go`,
`provider/live_databases_test.go`, `provider/live_harness_test.go`,
`provider/live_harness_unit_test.go`, `provider/live_workloads_test.go`,
`provider/mount_lifecycle_matrix_test.go`, `provider/registry_metadata_test.go`,
`provider/registry_test.go`, `tests/README.md`, `tests/acceptance_test.go`, and
`tests/acceptance_program_test.go`, plus the Task 1-8 reports and the two
acceptance design/plan documents.

Commits reviewed from `7a2f99c` through `bf81a76` include the Task 1-8 source,
workflow, documentation, and test changes. No review-fix commit was needed.

## Concerns and unavailable verification

Live execution was not attempted because the worktree has no acceptance
opt-in, Dokploy credentials, or configured server. `golangci-lint` and `mise`
are not installed. The first race invocation exposed a timing-sensitive
scripted-server failure; the required command passed on rerun, but this
transient behavior should be watched in CI.

## Final-review fixes

The follow-up review findings were addressed without changing provider runtime
behavior: Tier 2 workload and Mount redeploy calls now acquire operation-scoped
heavy leases; cleanup and health diagnostics are structural and allowlisted;
control-plane explicit deletes poll bounded eventual absence before releasing
cleanup ownership; and Application source assertions nil-check each variant
before field assertions.

Exact verification evidence from 2026-09-11:

- `gofmt -w provider/live_harness_test.go provider/live_harness_unit_test.go provider/application_source_test.go provider/live_control_plane_test.go provider/live_workloads_test.go provider/task9_regressions_test.go` — passed.
- `go test ./provider -run 'Test(ApplicationSource|ComposeSource|Mount|Domain|Destination|Registry|LiveGate|ClassifyLiveServerHealth|OwnedWorkflow|Task9|DeleteAndVerifyOnce|VerifiedCleanup|StopMarker)' -count=1` — passed (`ok`, 0.230s).
- `go test -short -count=1 ./provider/... ./internal/... ./tests/...` — passed.
- `go test -race ./provider/... ./internal/...` — passed (`ok provider`, 4.506s; internal packages passed).
- `go test ./... -count=1` — passed (all packages).
- `git diff --check` — passed.
- `golangci-lint` — unavailable (`command -v golangci-lint` returned no path).
- `mise exec -- golangci-lint run` — unavailable (`command -v mise` returned no path).

Post-commit reruns after the final source-variant lease coverage change:

- `go test ./provider -run 'Test(ApplicationSource|ComposeSource|Mount|Domain|Destination|Registry|LiveGate|ClassifyLiveServerHealth|OwnedWorkflow|Task9|DeleteAndVerifyOnce|VerifiedCleanup|StopMarker)' -count=1` — passed (`ok`, 0.240s).
- `go test -short -count=1 ./provider/... ./internal/... ./tests/...` — passed.
- `go test -race ./provider/... ./internal/...` — first rerun exposed the pre-existing timing-sensitive `TestBackupCreate_Cancellation` scripted-request failure; the exact command rerun passed (`ok provider`, 4.617s; internal packages passed).
- `go test ./... -count=1` — passed (all packages).
- `git diff --check` — passed.
