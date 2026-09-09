# Live Acceptance Coverage Hardening Design

## Goal

Close the identified live acceptance gaps while preserving the suite's safety
constraints for a low-powered Dokploy server. Prioritize complete source-state
reconstruction and Mount lifecycle confidence, then address the remaining
control-plane, Domain, server-health, workflow-artifact, and Pulumi-smoke gaps.

## Scope

This work covers all findings from the 2026-09-09 live acceptance review:

1. Application source variants do not verify complete reconstructed state.
2. Database-target Mount routes do not receive complete lifecycle coverage.
3. Control-plane mutable and write-only fields have uneven live coverage.
4. Compose source variants do not verify complete lifecycle reconstruction.
5. Domain tests omit safe HTTPS, strip-path, and custom-certificate behavior.
6. Live test paths do not classify server-health failures or propagate stops.
7. Acceptance workflow artifacts are skipped after a reported test failure.
8. Pulumi smoke tests do not assert operation summaries or removal evidence.

Provider behavior changes are outside this scope. If stronger acceptance tests
expose a provider defect, the test run will classify and report it separately.

## Strategy

Use layered hybrid coverage:

- Live tests prove representative end-to-end behavior against Dokploy.
- Deterministic scripted tests verify every source field, request shape, state
  reconstruction rule, and Mount type/target diff combination.
- Expensive live matrices remain intentionally reduced where repeating them
  would create disproportionate deployment load.
- External integrations remain prerequisite-gated rather than becoming
  mandatory acceptance infrastructure.

This strategy strengthens confidence without replacing the suite's serial,
bounded, fail-safe execution model.

## Application Source Coverage

Generic Git becomes the always-available complete live source case. The test
creates a provider-managed Application with a Git source, then verifies all
observable source fields after the initial read and an ID-only import-style
read. Assertions include URL, branch, build type, and every configured build
option that Dokploy returns. It performs a supported source metadata update if
the provider contract permits one without changing source type, verifies the
persisted state, and finishes with delete and post-delete absence checks.

Registry and GitLab sources use the same field-level assertions when their
dedicated environment variables are configured. Missing credentials or an
unmanaged GitLab integration produce explicit prerequisite skips. These cases
must not print credentials, integration identifiers, repository metadata, or
server response bodies.

Deterministic scripted-server tests cover every source-specific field for Git,
Docker/registry, and GitLab regardless of live prerequisites. They verify create
and configuration request shapes, normal reads, ID-only import reconstruction,
secret preservation where the API cannot return a write-only value, and diff
behavior for mutable and replacement-only source changes.

## Compose Source Coverage

Generic Git becomes the always-available complete Compose source case. It
verifies URL, branch, compose path or related source options when configured,
normal read, ID-only import reconstruction, supported metadata update, and
delete/absence behavior.

GitLab remains prerequisite-gated and receives field-level assertions for the
integration, project, owner, namespace, repository, branch, and compose path
values that are part of the provider contract. Deterministic tests cover these
fields when live credentials are absent. Raw Compose remains the representative
deployed workload used by Domain and Mount tests.

Source assertions compare sensitive values through redacted helpers. Test
failure output identifies only the resource, source type, and field name.

## Mount Coverage

Application and Compose retain their complete live bind, volume, and file Mount
lifecycles. Each lifecycle covers create, read, mutable update, post-update
read, ID-only import, type replacement diff, target replacement diff, delete,
and post-delete absence.

PostgreSQL becomes the representative database target for a complete live
Mount lifecycle. It uses one bind Mount and performs update, import, replacement
diff, delete, and absence checks before deleting the database fixture.

MySQL, MariaDB, and Redis retain serial live dispatch cases covering create,
read, ID-only import, delete, and absence. They do not repeat update-triggered
redeploys. Deterministic tests instead exercise mutable updates and replacement
diffs for every supported target across bind, volume, and file types.

The deterministic Cartesian matrix covers Application, Compose, PostgreSQL,
MySQL, MariaDB, and Redis. MongoDB and LibSQL remain documented exclusions
because they are not supported Mount targets. Any reduction from a full live
matrix is stated in test comments and the run summary.

## Control-Plane Coverage

Project and Tag live tests update and verify each safely mutable field rather
than relying on one representative value. Where multiple fields can be changed
in one safe request, the test may group them while retaining field-level read
assertions.

Destination live coverage updates and verifies provider, bucket, region,
endpoint, additional flags, and optional server scope where safe and available.
Access-key and secret-key behavior is verified without printing values. ID-only
reads must preserve write-only credentials from prior state only when the
provider contract requires prior state; true imports assert only fields that
Dokploy can reconstruct.

Registry live coverage updates and verifies name, username, URL, image prefix,
optional server scope, and password preservation. Registry network validation
and server prerequisites remain gated by dedicated registry configuration.
Deterministic tests cover write-only password behavior and optional fields
without relying on an external registry.

No test invents a replacement requirement for mutable control-plane fields.
Diff assertions follow the resource's documented contract.

## Domain Coverage

The existing Application and Compose target cases continue to use reserved
`.example.invalid` hosts and never request public certificates. Their mutable
update adds independent assertions for HTTPS and strip-path behavior while
using certificate type `none`.

Custom-certificate behavior is live-tested only when a dedicated custom
resolver environment variable is configured. The case updates or creates a
Domain using certificate type `custom`, verifies the resolver through read and
ID-only import, and performs normal cleanup. When the resolver is absent, the
case records an explicit prerequisite skip.

Deterministic tests always cover custom certificate request construction,
validation, read reconstruction, and diff behavior. Let's Encrypt remains
excluded from live tests because it could issue a public certificate request.

## Server-Health Stop Classification

Add a narrow classifier for failures that establish unsafe server state, such
as a failed bounded health probe, repeated heavy-operation timeout followed by
a failed health probe, or an explicit capacity/unavailable response recognized
by status and safe error code. Ordinary provider validation, decoding, and API
contract failures must not trigger a stop.

Heavy tiers perform a bounded health check at their existing safe boundaries:
before acquiring the first heavy-operation lease and after a heavy operation
fails in a way that may indicate server instability. A confirmed unhealthy
result calls `recordServerHealthFailure`, writes the non-secret stop marker, and
prevents later heavy tiers from starting.

Deterministic tests prove positive and negative classifications, bounded
contexts, sanitized diagnostics, and stop-marker propagation. Health checks do
not expose endpoint values or response bodies.

## Workflow Artifacts

Acceptance result aggregation continues to fail the job after all safe tiers
have run. Provider archive and artifact upload steps gain `if: ${{ always() }}`
so artifacts remain available after ordinary tier failures or stop-marker
failures. Archive creation must tolerate an unavailable binary and communicate
that state without masking the original acceptance failure.

Workflow contract tests verify step ordering, unconditional artifact handling,
serial tier commands, stop-marker gates, and final failure aggregation.

## Pulumi Smoke Evidence

The Automation API smoke inspects preview and update summaries rather than only
checking returned errors. Revision one must plan and create the expected four
resources. Revision two must report updates without replacement or deletion,
and resource IDs must remain stable.

Cleanup continues to use independent bounded contexts. After destroy, the test
queries stack outputs or export state through supported Automation API methods
to verify that managed resources are absent. After workspace stack removal, it
uses a supported list/select operation to verify that the stack no longer
exists locally. If Automation API cannot provide reliable direct evidence for
one removal phase, the test records the exact supported proxy assertion rather
than manipulating state files.

Deterministic tests cover summary validation and cleanup-result handling without
requiring Pulumi or live Dokploy credentials.

## Safety And Secrets

All existing safety constraints remain mandatory:

- Live execution requires explicit opt-in and both Dokploy credentials.
- Tests and subtests remain serial; no live test calls `t.Parallel()`.
- Only one heavy operation may be active at a time.
- Every created resource receives immediate bounded fallback cleanup.
- Explicit delete paths verify eventual absence.
- A cleanup or confirmed server-health failure stops later heavy work.
- Backup definitions remain disabled and are never executed.
- Test code never reads `.env`.
- Credentials, source secrets, SSH material, database passwords, mount content,
  endpoint values, integration identifiers, IDs with sensitive context, and raw
  request or response bodies never appear in diagnostics.

New fixture values must be registered with the live secret registry before an
operation can expose them through an error.

## Verification

Implementation follows test-driven development. Each gap first receives a
focused deterministic failing test, followed by the smallest production-test
harness or live-test change needed to pass it.

Static verification includes:

- Focused source, Mount, Domain, harness, workflow, and Automation API tests.
- The short provider, internal, and acceptance suites.
- Race tests for provider and internal packages.
- Workflow contract tests.
- Formatting, lint, and `git diff --check`.

Live verification runs tiers serially with `-parallel=1 -count=1 -v`. Tier 2 is
run first for source, Domain, and Mount coverage; focused control-plane and
Pulumi smoke runs follow. Database or backup reruns occur only when changed
shared harness behavior affects them. The run summary records all skips,
durations, cleanup results, and hybrid-matrix limitations.

## Acceptance Criteria

- Application Git live coverage verifies complete observable source state on
  normal read and ID-only import.
- Configured registry and GitLab Application sources receive the same field-level
  assertions; missing prerequisites produce explicit skips.
- Compose Git and configured GitLab sources verify complete observable state and
  import reconstruction.
- PostgreSQL Mounts receive a complete live lifecycle, while every supported
  Mount type/target combination receives deterministic update and replacement
  coverage.
- Control-plane tests verify safely mutable optional fields and write-only
  credential preservation according to each provider contract.
- Domain live tests cover HTTPS and strip path without public certificates, and
  custom resolvers are tested when configured.
- Confirmed server-health failures create the stop marker; ordinary provider
  failures do not.
- Acceptance artifacts upload even when a tier or stop-marker gate fails.
- Pulumi smoke verifies expected create/update summaries, stable IDs, and the
  strongest supported destroy/removal evidence without state manipulation.
- All new diagnostics remain sanitized, execution remains serial, and no new
  unbounded cleanup or polling path is introduced.
