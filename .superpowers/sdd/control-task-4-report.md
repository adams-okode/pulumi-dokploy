# Task 4 Report: Control-Plane Mutable And Secret Coverage

## RED/GREEN evidence

- RED: no production behavior was required. The existing provider contract
  already preserved write-only values from prior state and reconstructed the
  API-observable fields on ID-only reads; the pre-change focused suite passed.
  Adding a deliberately artificial failing expectation would not test a real
  requirement.
- GREEN: the new deterministic Destination and Registry reconstruction tests
  pass. They use `requireLiveEqual` for secret and endpoint assertions and do
  not expose sentinel values in failure formatting.
- Live control-plane execution is prerequisite-gated and skipped without
  `DOKPLOY_ACCEPTANCE`, `DOKPLOY_ENDPOINT`, and `DOKPLOY_API_KEY`.

## Files

- `provider/destination_test.go`: prior-state secret preservation and ID-only
  observable-field reconstruction coverage.
- `provider/registry_test.go`: prior-state password preservation and ID-only
  observable-field reconstruction coverage.
- `provider/live_control_plane_test.go`: field-level Project, Tag,
  Destination, and Registry update assertions; optional server scope is used
  only when `DOKPLOY_ACCEPTANCE_SERVER_ID` is configured; changed sensitive
  fixtures are registered.

## Verification

- `go test ./provider -run 'Test(Destination|Registry)Read(Preserves|Reconstructs)' -count=1` — PASS
- `go test ./provider -run 'Test(Destination|Registry)Read|TestLiveTier1ControlPlane' -count=1 -v` — PASS; live tier SKIP (missing opt-in)
- `go test ./provider -count=1` — PASS
- `go test -short -count=1 ./provider/... ./internal/... ./tests/...` — PASS
- `go test -race ./provider/... ./internal/...` — PASS
- `go test ./... -count=1` — PASS
- `git diff --check` — PASS

## Concerns

- Live Registry connection fields retain configured valid credentials and URL
  during the update so prerequisite validation remains safe; the update
  request and read assertions still cover every mutable Registry field.
- No live Dokploy credentials were available, so live API behavior was not
  exercised in this worktree.

## Commit

- `ced216a4a807ceb05d1ea48499dc4751ccf6ec07` — implementation and test
  coverage.
- This report is committed in the follow-up documentation commit.
