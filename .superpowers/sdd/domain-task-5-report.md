# Task 5 Report: Safe Domain TLS and Routing Fields

## Status

Implemented and verified. No production provider behavior was changed.

## Evidence

- Added deterministic scripted custom-certificate create/update/read coverage,
  including resolver validation coverage already exercised by the Domain check
  tests.
- Added deterministic HTTPS and `stripPath` updates with `certificateType:
  none`, normal reads, and ID-only reads.
- Extended serial Tier 2 Application and Compose Domain coverage to update and
  verify HTTPS and strip path independently while retaining certificate `none`.
- Added resolver-gated live custom-certificate cases for both targets. The
  resolver is registered with `registerLiveSecrets`; existing bounded,
  disarmable cleanup owners remain in use.
- Added `DOKPLOY_CUSTOM_CERT_RESOLVER` only to the Tier 2 workflow environment.
- Updated the live acceptance README and documentation contract from planned to
  active wording, plus workflow metadata contract expectations.
- No live fixture selects Let's Encrypt.

## Commands and outcomes

1. `go test ./provider -run 'TestDomainCustomCertificateRoundTrip|TestDomainSafeRoutingFieldsRoundTrip' -count=1`
   - PASS.
2. `go test ./provider -run 'TestDomain|TestLiveTier2Workloads/Domain' -count=1 -v`
   - PASS; live Tier 2 was skipped because acceptance credentials were not configured.
3. `go test ./provider -run 'TestRegistryMetadata|TestOwnedWorkflow|TestAcceptanceWorkflowContractsRejectRepresentativePermissionDrift|TestWorkflowStepsDoNotHaveEmptyEnvMappings|TestExampleTestWorkflowsRunFromExamplesDirectory' -count=1`
   - PASS.
4. `go test ./tests -count=1`
   - PASS.
5. `go test ./... -count=1`
   - PASS.
6. `git diff --check`
   - PASS.

## Concerns

- Live Dokploy execution was not available in this environment, so the new
  resolver-gated cases were not exercised against a server.
- The worktree contained a pre-existing deletion of
  `.superpowers/sdd/task-1-report.md`; it was not modified or staged.

## Review follow-up: Domain diff coverage

- Added `TestDomainDiffCertificateFieldsAreMutable` covering independent
  `certificateType` and `customCertResolver` changes, asserting both are
  documented `Update` diffs rather than replacements, and asserting unchanged
  inputs produce no diff.
- TDD RED evidence: `go test ./provider -run TestDomainDiffCertificateFieldsAreMutable -count=1`
  failed because the temporary first assertion expected `update&replace` while
  the implementation returned `update`.
- TDD GREEN evidence: after correcting that assertion to the documented mutable
  kind, the focused Domain and full test suites passed.

Additional verification:

1. `go test ./provider -run 'TestDomain' -count=1` - PASS.
2. `go test ./provider -run 'TestDomain|TestLiveTier2Workloads/Domain' -count=1 -v`
   - PASS; live acceptance skipped without opt-in credentials.
3. `go test ./... -count=1` - PASS.
4. `git diff --check` - PASS.
