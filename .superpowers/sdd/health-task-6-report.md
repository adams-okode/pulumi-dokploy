# Task 6: Narrow Server-Health Stop Classification

Status: complete

Implemented a bounded, sanitized server-health classifier and authenticated
low-cost project probe. Heavy tiers probe before acquiring a lease; explicit
unavailable/capacity create failures stop immediately, while operation
timeouts require a failed follow-up probe. Validation, decode, not-found, and
ordinary operation timeout failures do not stop later tiers. Existing opt-in,
serial lease, cleanup, and stop-marker behavior remain intact.

Verification:

- `go test ./provider -run 'TestClassifyLiveServerHealthFailure|TestVerifyLiveServerHealth' -count=1` — PASS
- `go test ./provider -run 'Test(ClassifyLiveServerHealthFailure|VerifyLiveServerHealth|ServerHealthFailure|OrdinaryLiveResult)' -count=1` — PASS
- `env -u DOKPLOY_ACCEPTANCE go test ./provider -run 'TestLiveTier(2|3|4)' -count=1 -v` — PASS; all three tiers skipped before probing
- `go test ./provider -count=1` — PASS
- `go test ./... -count=1` — PASS
- `git diff --check` — PASS

Concerns: live Dokploy execution was not performed because acceptance opt-in
and credentials were not available. The probe intentionally uses an invalid
project ID and treats responsive 4xx responses as healthy; it records only
fixed structural diagnostics and never includes response bodies.
