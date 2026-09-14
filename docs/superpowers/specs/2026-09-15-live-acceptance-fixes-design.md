# Live Acceptance Fixes Design

## Goal

Fix the proven PostgreSQL mount-dispatch readiness bug and establish safe,
decisive evidence for the Domain create failures before changing production
Domain behavior.

## Scope

This work has two independently reviewable tracks:

1. Correct resource-specific readiness and stop semantics for mount-dispatch
   fixtures.
2. Add a bounded, sanitized Domain provider-versus-generated-client comparison,
   then make the smallest production correction only if that evidence proves a
   provider mismatch.

The work does not add speculative Domain payload fallbacks, weaken cleanup
requirements, or convert confirmed failures into skips.

## Current Evidence

`TestLiveTier2Workloads/MountDispatch/postgres` creates a PostgreSQL fixture and
passes its ID to `runLiveMountLifecycle`. That helper calls
`readLiveWorkloadTarget`, which supports only Application and Compose. Because
the PostgreSQL mount has no Compose ID, the helper queries `Application.Read`
with a PostgreSQL ID and reports the target missing. The cleanup/health path then
creates a stop marker. This behavior reproduced in full Tier 2 runs and in an
isolated focused run.

Both Application and Compose Domain creates return a sanitized
`4xx/BAD_REQUEST` classification on multiple environments. The provider request
keys and generated request body match the pinned OpenAPI contract. Existing
evidence therefore does not justify changing the production payload.

## Architecture

### Resource-specific dispatch readiness

Represent readiness as behavior owned by each dispatch fixture rather than
inferring the resource type inside the generic mount lifecycle. A dispatch
fixture exposes a bounded readiness function that reads its concrete resource
and returns presence, readiness, and error. Application and Compose use their
existing resource reads. PostgreSQL, MySQL, MariaDB, and Redis use their own
provider reads and require a non-empty matching ID plus `statusDone`.

`runLiveMountLifecycle` receives this readiness behavior explicitly. Ordinary
Application/Compose mount cases retain their existing readers. Dispatch cases
pass the fixture reader. The mount helper does not select an API endpoint from
pointer presence or treat a database ID as an Application ID.

### Domain comparison diagnostic

Add a live-only comparison path that submits equivalent validated Domain inputs
through the provider resource and the generated client. It records only:

- target kind (`application` or `compose`),
- HTTP status class,
- allowlisted API error code,
- sorted request-key names,
- whether a target was present and ready.

It must not record resource IDs, hosts, endpoint values, credentials, request or
response values, or response bodies. Both attempts are serial and bounded. Any
partial Domain returned by either path is registered for immediate verified
cleanup before its result is asserted or compared.

The comparison classifies one of three outcomes:

- **Provider serialization mismatch:** generated direct create succeeds while
  provider create fails, or safe structural request evidence differs. This
  authorizes a minimal provider correction.
- **Server contract rejection:** provider and generated direct create submit the
  same contract shape and receive the same rejection. No provider payload change
  is made; the report identifies the deployed Dokploy contract or migration as
  the next owner.
- **Environment/server-health failure:** transport, timeout, `5xx`, cleanup, or
  health-probe evidence prevents comparison. The stop marker is preserved and no
  later heavy test runs.

## Data Flow

### Dispatch fixture

1. Acquire the heavy-operation lease after a successful health probe.
2. Create and deploy the concrete fixture.
3. Preserve the returned ID and register fixture ownership before assertions.
4. Invoke the fixture's concrete readiness function.
5. Create and verify the Mount only when the target is present and `done`.
6. Delete and verify the Mount before deleting the fixture.
7. Release the fixture lease only after verified fixture absence.

### Domain comparison

1. Create ready Application and Compose targets under existing parent ownership.
2. Construct one validated `DomainArgs` value per target kind.
3. Derive the generated request from `domainCreateBody` and verify key parity.
4. Run provider and generated-client attempts serially with distinct reserved
   hosts.
5. Register cleanup immediately for every returned Domain ID.
6. Compare only sanitized structural outcomes.
7. Apply no production change until the comparison identifies a provider-owned
   difference.

## Error And Stop Semantics

A missing or not-ready dispatch target is an ordinary fixture failure when the
concrete read succeeds. It does not create a stop marker by itself. A marker is
created only for failed verified cleanup or an independently confirmed
server-health failure, preserving the documented acceptance contract.

Read errors retain their safe structural classification. Transport errors,
timeouts followed by failed health probes, and `502`, `503`, or `504` responses
remain server-health stop conditions. IDs and raw API errors never enter test
output.

Domain `4xx` rejection is an ordinary compatibility failure unless cleanup
fails. The direct comparison must not retry alternate payloads, because repeated
creates can produce duplicate resources and conceal the actual contract.

## Testing

### Deterministic tests

- Reproduce the current PostgreSQL bug by proving that the old generic reader
  selects Application for a database fixture.
- Verify explicit readiness dispatch for Application, Compose, PostgreSQL,
  MySQL, MariaDB, and Redis.
- Verify each database readiness function requires a present ID and `done`
  status and propagates read errors.
- Verify missing/not-ready fixture results do not create a stop marker without a
  cleanup or health failure.
- Verify fixture cleanup remains exactly once and keeps the lease through
  verified absence.
- Verify Domain provider/generated request-key parity for Application and
  Compose.
- Verify Domain result classification emits only allowlisted structural fields.
- Verify partial Domain creation is cleaned before a fatal comparison result.
- Verify diagnostic output excludes configured secret and identifier sentinels.

### Verification sequence

1. Run focused new tests and observe them fail before implementation.
2. Implement the minimal dispatch correction and diagnostic path.
3. Run focused tests until they pass.
4. Run the provider short suite.
5. Run the provider race suite.
6. Run lint and repository static checks.
7. With a fresh absent stop marker, run the focused PostgreSQL dispatch live
   test.
8. If its marker remains absent, run focused Domain Application and Compose
   diagnostics.
9. If diagnostics prove a provider mismatch, add a separate minimal production
   correction and regression test, then repeat steps 3 through 8.
10. Run the full serial tier sequence only while the marker remains absent after
    every command.

## Success Criteria

- The focused PostgreSQL dispatch test reads PostgreSQL, not Application, and
  completes mount and fixture cleanup without creating a stop marker.
- MySQL, MariaDB, and Redis dispatches use their corresponding readers; existing
  Application and Compose behavior remains intact.
- Ordinary target absence cannot be mislabeled as a server-health failure.
- Domain comparison produces enough sanitized evidence to assign ownership to
  provider serialization or deployed Dokploy behavior.
- Production Domain behavior changes only when a deterministic or controlled
  live comparison proves the exact mismatch.
- No test output or tracked report contains credentials, endpoint values,
  resource IDs, private hosts, payload values, or response bodies.
- Later heavy tiers never run after a cleanup or confirmed health stop marker.

## Delivery Boundaries

The dispatch correction is independently shippable after deterministic and live
verification. Domain diagnostics are independently shippable as test-only
hardening. Any production Domain correction is a subsequent evidence-gated task
and is not required when comparison demonstrates a server-owned incompatibility.
