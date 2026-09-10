# Task 2 report: Compose source reconstruction

## Status

Implemented and committed as `81299fadc38b06cb7fa0835504b4f3a9f6b298fc` (`test: verify compose source reconstruction`).

## RED/GREEN evidence

- RED review: the pre-existing Compose tests covered only type/path spot checks; they did not exercise complete Git/GitLab normal and ID-only reconstruction or the live lifecycle. The new focused cases were designed before the live changes and expose every source field through the shared redacted assertion helper.
- GREEN: `go test ./provider -run 'TestComposeSourceIDOnlyReadReconstructsAllFields|TestComposeSourceNormalReadReconstructsAllFields' -count=1` passed.
- GREEN: `go test ./provider -run 'TestComposeSource|TestLiveTier2Workloads/SourceVariants' -count=1 -v` passed; live acceptance was skipped because opt-in credentials were absent.

## Changes

- Added deterministic normal-read and ID-only scripted-server tests for Git and GitLab Compose sources.
- Added complete, secret-safe assertions for URL, branch, compose path, SSH key, watch paths, submodules, and all GitLab integration/project/owner/namespace/repository fields.
- Hardened live Git and prerequisite-gated GitLab source cases with source secret registration, normal/import reads, safe same-type metadata updates, explicit bounded cleanup, and post-delete absence checks.
- Reused `cleanupDirectCompose` and existing live cleanup/diagnostic helpers; production code was unchanged.

## Verification

Commands and results:

```text
go test ./provider -run 'TestComposeSourceIDOnlyReadReconstructsAllFields|TestComposeSourceNormalReadReconstructsAllFields' -count=1
PASS

go test ./provider -run 'TestComposeSource|TestLiveTier2Workloads/SourceVariants' -count=1 -v
PASS; TestLiveTier2Workloads skipped without DOKPLOY_ACCEPTANCE, DOKPLOY_ENDPOINT, and DOKPLOY_API_KEY

go test ./provider
PASS

go test ./...
PASS

go test -race ./provider ./internal/...
PASS

git diff --check
PASS
```

## Files and commit

- `provider/compose_source_test.go`
- `provider/live_workloads_test.go`
- Commit: `81299fadc38b06cb7fa0835504b4f3a9f6b298fc`

## Concerns

- Live Dokploy execution was not available in this environment, so live API update/delete behavior remains unobserved here.
- The worktree already had an unrelated deletion of `.superpowers/sdd/task-1-report.md`; it was not modified or staged.
