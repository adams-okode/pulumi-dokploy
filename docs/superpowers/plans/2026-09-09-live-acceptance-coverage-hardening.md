# Live Acceptance Coverage Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Close all reviewed live acceptance gaps with high-priority source and Mount coverage first, while retaining safe execution on a low-powered Dokploy server.

**Architecture:** Use deterministic scripted-server matrices for exhaustive field and dispatch contracts, plus targeted serial live lifecycles for representative end-to-end behavior. Extend the shared harness only for narrow server-health classification, retain prerequisite gates for external integrations, and strengthen workflow and Automation API evidence without manipulating Pulumi state.

**Tech Stack:** Go 1.26.6, `testing`, `testify/require`, `pulumi-go-provider/infer`, Pulumi Automation API 3.259.0, generated Dokploy client, GitHub Actions YAML.

## Global Constraints

- Follow `docs/superpowers/specs/2026-09-09-live-acceptance-coverage-hardening-design.md`.
- Live execution requires `DOKPLOY_ACCEPTANCE=1`, `DOKPLOY_ENDPOINT`, and `DOKPLOY_API_KEY`.
- Test code never reads `.env`.
- No live test or subtest calls `t.Parallel()`.
- Only one heavy create, deploy, redeploy, image pull, or database startup may be active at a time.
- Every created resource receives immediate bounded fallback cleanup and explicit delete paths verify eventual absence.
- Cleanup failure or confirmed server-health failure writes the non-secret stop marker and prevents later heavy work.
- Ordinary validation, decoding, provider-contract, and API failures do not stop independent later tiers.
- Backup and VolumeBackup definitions stay disabled and no execute endpoint is called.
- Registry, GitLab, and custom-certificate live cases remain gated by dedicated prerequisites.
- Live Domain tests never request Let's Encrypt or another public certificate.
- Credentials, integration identifiers, endpoint values, source secrets, SSH material, database passwords, Mount content, and raw request/response bodies never appear in diagnostics.
- Register every new sensitive fixture value with `registerLiveSecrets` before issuing an operation that could expose it.
- Provider production behavior is not changed merely to make a new acceptance expectation pass; classify and report confirmed defects separately.

---

### Task 1: Complete Application Source Reconstruction

**Files:**
- Modify: `provider/application_source_test.go`
- Modify: `provider/live_workloads_test.go:460-493`
- Test: `provider/application_source_test.go`

**Interfaces:**
- Consumes: `configureApplicationSource`, `Application.Read`, `Application.Diff`, `cleanupDirectApplication`, `liveRegistryApplicationSource`, and `liveGitLabApplicationSource`.
- Produces: `assertLiveApplicationSource(t *testing.T, label string, want, got ApplicationSource)` and complete deterministic/live source reconstruction coverage.

- [ ] **Step 1: Write failing deterministic source reconstruction tests**

Add table-driven scripted-server cases for Git, Docker registry, and GitLab. For each source, script the provider endpoints used by an ID-only `Application.Read`, then compare every observable source field through a field-name-only helper. Include Git URL, branch, build path, SSH key ID, submodules, watch paths, build type, Dockerfile or publish directory where applicable; Docker image, registry URL, username, and preserved password; and all GitLab integration/project/repository/build fields.

```go
func assertApplicationSourceFields(t *testing.T, want, got ApplicationSource) {
	t.Helper()
	require.Equal(t, want.Type, got.Type)
	switch want.Type {
	case SourceGit:
		requireLiveEqual(t, "application.source.git.url", want.Git.URL, got.Git.URL)
		requireLiveEqual(t, "application.source.git.branch", want.Git.Branch, got.Git.Branch)
		requireLiveEqual(t, "application.source.git.build", want.Git.Build, got.Git.Build)
	case SourceDocker:
		requireLiveEqual(t, "application.source.docker.image", want.Docker.Image, got.Docker.Image)
		requireLiveEqual(t, "application.source.docker.registryUrl", want.Docker.RegistryURL, got.Docker.RegistryURL)
		requireLiveEqual(t, "application.source.docker.username", want.Docker.Username, got.Docker.Username)
	case SourceGitLab:
		requireLiveEqual(t, "application.source.gitlab.projectId", want.GitLab.ProjectID, got.GitLab.ProjectID)
		requireLiveEqual(t, "application.source.gitlab.repository", want.GitLab.Repository, got.GitLab.Repository)
		requireLiveEqual(t, "application.source.gitlab.build", want.GitLab.Build, got.GitLab.Build)
	}
}
```

Use exact response fields from the generated client and existing source-fetch functions; do not add a second parser in test code.

- [ ] **Step 2: Run the source tests and verify RED**

Run: `go test ./provider -run 'TestApplicationSourceIDOnlyReadReconstructsAllFields|TestApplicationSourceWriteOnlySecretsArePreserved' -count=1`

Expected: FAIL because existing tests do not provide complete scripted reads or the field helper.

- [ ] **Step 3: Add complete live source assertions**

Replace the metadata-only Application source loop with provider-managed lifecycle cases. Generic Git always runs; registry and GitLab skip when their existing prerequisite builders return an empty type. Each case creates a bare Application, configures its source, reads with prior state, performs an ID-only read, compares all observable fields, checks a same-type metadata update when supported, and explicitly deletes with verified absence.

```go
func assertLiveApplicationSource(t *testing.T, label string, want, got ApplicationSource) {
	t.Helper()
	if want.Type != got.Type {
		t.Fatalf("live Application %s source type did not match", label)
	}
	assertApplicationSourceFields(t, want, got)
}
```

Call `registerLiveSecrets` for registry credentials and GitLab identifiers before create/configure calls. Never format either source object in an assertion.

- [ ] **Step 4: Verify deterministic and skip-safe behavior**

Run: `go test ./provider -run 'TestApplicationSource|TestLiveTier2Workloads/SourceVariants' -count=1 -v`

Expected: deterministic source tests PASS; live cases SKIP without explicit acceptance credentials.

- [ ] **Step 5: Commit**

```bash
git add provider/application_source_test.go provider/live_workloads_test.go
git commit -m "test: verify application source reconstruction"
```

### Task 2: Complete Compose Source Reconstruction

**Files:**
- Modify: `provider/compose_source_test.go`
- Modify: `provider/live_workloads_test.go:494-536`
- Test: `provider/compose_source_test.go`

**Interfaces:**
- Consumes: `configureComposeSource`, `fetchComposeSource`, `Compose.Read`, `cleanupDirectCompose`, `gitLabComposeSource`.
- Produces: `assertComposeSourceFields(t *testing.T, want, got ComposeSource)` and complete Git/GitLab reconstruction assertions.

- [ ] **Step 1: Write failing deterministic Compose source tests**

Add scripted tests for Git and GitLab normal reads and ID-only reads. Assert URL, branch, compose path, and every GitLab integration/project/owner/namespace/repository field.

```go
func assertComposeSourceFields(t *testing.T, want, got ComposeSource) {
	t.Helper()
	require.Equal(t, want.Type, got.Type)
	if want.Git != nil {
		requireLiveEqual(t, "compose.source.git.url", want.Git.URL, got.Git.URL)
		requireLiveEqual(t, "compose.source.git.branch", want.Git.Branch, got.Git.Branch)
		requireLiveEqual(t, "compose.source.git.composePath", want.Git.ComposePath, got.Git.ComposePath)
	}
	if want.GitLab != nil {
		requireLiveEqual(t, "compose.source.gitlab.integrationId", want.GitLab.IntegrationID, got.GitLab.IntegrationID)
		requireLiveEqual(t, "compose.source.gitlab.projectId", want.GitLab.ProjectID, got.GitLab.ProjectID)
		requireLiveEqual(t, "compose.source.gitlab.repository", want.GitLab.Repository, got.GitLab.Repository)
		requireLiveEqual(t, "compose.source.gitlab.composePath", want.GitLab.ComposePath, got.GitLab.ComposePath)
	}
}
```

- [ ] **Step 2: Run the tests and verify RED**

Run: `go test ./provider -run 'TestComposeSourceIDOnlyReadReconstructsAllFields|TestComposeSourceNormalReadReconstructsAllFields' -count=1`

Expected: FAIL until complete source fixtures and assertions exist.

- [ ] **Step 3: Replace metadata-only live Compose source cases**

Make generic Git a complete provider-managed live case with configure, fetch, prior-state read, ID-only read, same-type metadata update where supported, explicit delete, and post-delete absence. Apply the same field assertions to GitLab when prerequisites exist. Raw Compose remains the deployed dependency fixture and must not be duplicated.

Register repository and integration fixture values as secrets before operations. Missing GitLab prerequisites remain an explicit skip.

- [ ] **Step 4: Verify GREEN**

Run: `go test ./provider -run 'TestComposeSource|TestLiveTier2Workloads/SourceVariants' -count=1 -v`

Expected: deterministic tests PASS and live tests SKIP without opt-in.

- [ ] **Step 5: Commit**

```bash
git add provider/compose_source_test.go provider/live_workloads_test.go
git commit -m "test: verify compose source reconstruction"
```

### Task 3: Finish Hybrid Mount Lifecycle Coverage

**Files:**
- Modify: `provider/mount_lifecycle_matrix_test.go`
- Modify: `provider/live_workloads_test.go:276-458,708-792`
- Test: `provider/mount_lifecycle_matrix_test.go`

**Interfaces:**
- Consumes: `mountLifecycleTargetCases`, `setMountTarget`, `mountReplacement`, `createDispatchDatabase`, `deleteAndVerifyLiveOwned`.
- Produces: `runLiveMountLifecycle(t *testing.T, ctx context.Context, api *client.Client, inputs MountArgs, targetReplacement MountArgs)` and exhaustive deterministic update/diff matrices.

- [ ] **Step 1: Add failing deterministic target/type diff matrix**

Extend the existing six-target by three-type matrix to assert both mutable update and replacement behavior. For each case, change `mountPath` and the type-specific mutable field, then assert `p.Update`; separately change type and target, then assert `p.UpdateReplace` on `type` and the relevant target properties.

```go
func TestMountDiffCartesianMatrix(t *testing.T) {
	for _, target := range mountLifecycleTargetCases() {
		for _, mountType := range []string{mountTypeBind, mountTypeVolume, mountTypeFile} {
			t.Run(target.name+"/"+mountType, func(t *testing.T) {
				base := mountArgsForMatrix(mountType, target.idField, target.id)
				mutable := base
				mutable.MountPath += "-updated"
				diff, err := (Mount{}).Diff(t.Context(), infer.DiffRequest[MountArgs, MountState]{Inputs: mutable, State: MountState{MountArgs: base}})
				require.NoError(t, err)
				require.Equal(t, p.Update, diff.DetailedDiff["mountPath"].Kind)
			})
		}
	}
}
```

- [ ] **Step 2: Run the matrix and verify RED**

Run: `go test ./provider -run 'TestMountDiffCartesianMatrix|TestMountUpdateBodyAndRedeployMatrix' -count=1`

Expected: FAIL until all update and replacement expectations are represented.

- [ ] **Step 3: Extract one complete live Mount lifecycle helper**

Move the existing Application/Compose lifecycle sequence into `runLiveMountLifecycle`. The helper must register cleanup immediately, perform create/read/update/post-update read/ID-only import/type diff/target diff/delete/absence, and use only field-name diagnostics.

```go
func runLiveMountLifecycle(t *testing.T, ctx context.Context, api *client.Client, inputs, targetReplacement MountArgs) {
	t.Helper()
	r := Mount{client: fixedClient(api)}
	created, err := r.Create(ctx, infer.CreateRequest[MountArgs]{Inputs: inputs})
	// Register cleanup before asserting err, then execute the existing lifecycle.
	// Each assertion identifies only a Mount field, never an input value.
}
```

Use the complete implementation sequence already present in `Mounts/application` and `Mounts/compose`; do not leave the helper body as comments.

- [ ] **Step 4: Upgrade PostgreSQL dispatch to the complete live lifecycle**

In `MountDispatch/postgres`, call `runLiveMountLifecycle` with a bind Mount and a safe replacement target. Keep MySQL, MariaDB, and Redis as serial create/read/import/delete dispatch checks. Keep Compose dispatch unchanged because Compose already has full coverage above.

Ensure the PostgreSQL fixture remains leased until Mount absence is verified, then clean and release it before MySQL starts.

- [ ] **Step 5: Verify GREEN and serial safety**

Run: `go test ./provider -run 'TestMountDiffCartesianMatrix|TestMountUpdateBodyAndRedeployMatrix|TestLiveTier2Workloads/Mount' -count=1 -v`

Expected: deterministic matrices PASS; live tier SKIPS without opt-in; no `t.Parallel()` exists in live files.

Run: `go test ./provider -run 'TestHeavyOperation|TestMountTargetDispatch' -count=1`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add provider/mount_lifecycle_matrix_test.go provider/live_workloads_test.go
git commit -m "test: complete hybrid mount lifecycle coverage"
```

### Task 4: Expand Control-Plane Mutable And Secret Coverage

**Files:**
- Modify: `provider/live_control_plane_test.go:29-248`
- Modify: `provider/destination_test.go`
- Modify: `provider/registry_test.go`
- Test: `provider/destination_test.go`
- Test: `provider/registry_test.go`

**Interfaces:**
- Consumes: `requireLiveEqual`, `registerLiveSecrets`, resource `Read`, `Update`, and `Diff` methods.
- Produces: complete field-level control-plane update assertions and deterministic write-only-state tests.

- [ ] **Step 1: Add failing Destination and Registry secret-state tests**

Add scripted tests proving that reads with prior state preserve write-only secret fields, while ID-only imports reconstruct only API-observable fields. Cover Destination `SecretAccessKey` and Registry `Password`, plus optional `AdditionalFlags`, `ImagePrefix`, and `ServerID` fields.

```go
func TestRegistryReadPreservesWriteOnlyPasswordFromPriorState(t *testing.T) {
	prior := RegistryState{RegistryArgs: RegistryArgs{Password: "password-sentinel"}}
	// Script registry.one without a password, call Read with prior, and assert
	// through requireLiveEqual so a failure never formats the sentinel.
}
```

- [ ] **Step 2: Run focused tests and verify RED**

Run: `go test ./provider -run 'Test(Destination|Registry)Read(Preserves|Reconstructs)' -count=1`

Expected: FAIL until the complete fixtures and assertions are added.

- [ ] **Step 3: Expand live mutable updates safely**

Update Project name and description and assert both after read. Update Tag name and color. For Destination, update all API-safe fields in one request: name, provider, access key, secret access key, bucket, region, endpoint, and additional flags; update server scope only when a dedicated `DOKPLOY_ACCEPTANCE_SERVER_ID` is configured. For Registry, update name, username, password, URL, and image prefix; update server scope under the same optional prerequisite.

Register every changed credential and endpoint value before calling Update. Use `requireLiveEqual` for credential assertions and ordinary equality only for non-sensitive fields.

- [ ] **Step 4: Verify GREEN**

Run: `go test ./provider -run 'Test(Destination|Registry)Read|TestLiveTier1ControlPlane' -count=1 -v`

Expected: deterministic tests PASS and live tier SKIPS without opt-in.

- [ ] **Step 5: Commit**

```bash
git add provider/live_control_plane_test.go provider/destination_test.go provider/registry_test.go
git commit -m "test: expand control plane lifecycle coverage"
```

### Task 5: Cover Safe Domain TLS And Routing Fields

**Files:**
- Modify: `provider/domain_test.go`
- Modify: `provider/live_workloads_test.go:192-274`
- Modify: `.github/workflows/run-acceptance-tests.yml:100-113`
- Test: `provider/domain_test.go`

**Interfaces:**
- Consumes: `Domain.Check`, `Domain.Read`, `Domain.Update`, `Domain.Diff`, `registerLiveSecrets`.
- Produces: prerequisite-gated `DOKPLOY_CUSTOM_CERT_RESOLVER` live coverage and deterministic custom-certificate coverage.

- [ ] **Step 1: Add failing deterministic custom-certificate tests**

Add scripted create/update/read tests asserting `certificateType: custom`, the resolver field, HTTPS, and strip path. Include validation that custom requires a non-empty resolver and that Let's Encrypt is not used by live fixtures.

```go
func TestDomainCustomCertificateRoundTrip(t *testing.T) {
	resolver := "resolver-sentinel"
	inputs := DomainArgs{ApplicationID: stringPtr("a1"), Host: "example.invalid", HTTPS: true, CertificateType: CertificateCustom, CustomCertResolver: &resolver, StripPath: true, Enabled: true}
	// Script create and read; compare each field with requireLiveEqual.
}
```

- [ ] **Step 2: Run the tests and verify RED**

Run: `go test ./provider -run 'TestDomainCustomCertificateRoundTrip|TestDomainSafeRoutingFieldsRoundTrip' -count=1`

Expected: FAIL until complete request and read fixtures exist.

- [ ] **Step 3: Extend live Domain updates**

For both Application and Compose targets, update HTTPS and strip path independently while certificate type remains `none`, then verify both through normal and ID-only reads. Add a nested `custom-certificate` case that runs only when `DOKPLOY_CUSTOM_CERT_RESOLVER` is non-empty; register the resolver as a secret, set certificate type `custom`, verify resolver reconstruction, and clean up normally.

When absent, use `t.Skip("DOKPLOY_CUSTOM_CERT_RESOLVER is not configured")`. Never select `CertificateLetsencrypt` in live code.

- [ ] **Step 4: Pass the prerequisite into CI and verify GREEN**

Add `DOKPLOY_CUSTOM_CERT_RESOLVER: ${{ secrets.DOKPLOY_CUSTOM_CERT_RESOLVER }}` only to Tier 2's environment.

Run: `go test ./provider -run 'TestDomain|TestLiveTier2Workloads/Domain' -count=1 -v`

Expected: deterministic tests PASS; live cases SKIP without acceptance opt-in.

- [ ] **Step 5: Commit**

```bash
git add provider/domain_test.go provider/live_workloads_test.go .github/workflows/run-acceptance-tests.yml
git commit -m "test: cover safe domain certificate behavior"
```

### Task 6: Add Narrow Server-Health Stop Classification

**Files:**
- Modify: `provider/live_harness_test.go`
- Modify: `provider/live_harness_unit_test.go`
- Modify: `provider/live_workloads_test.go`
- Modify: `provider/live_databases_test.go`
- Modify: `provider/live_backups_test.go`
- Test: `provider/live_harness_unit_test.go`

**Interfaces:**
- Produces: `classifyLiveServerHealthFailure(err error) bool` and `verifyLiveServerHealth(ctx context.Context, probe func(context.Context) error) error`.
- Consumes: `recordServerHealthFailure`, `beginLiveHeavyOperation`, and existing bounded contexts.

- [ ] **Step 1: Write failing classifier and probe tests**

Table-test API errors and context failures. Accept only transport/unavailable/capacity classifications such as HTTP 502, 503, 504, safe codes `SERVICE_UNAVAILABLE`, `SERVER_UNHEALTHY`, or `CAPACITY_EXHAUSTED`, and a bounded probe timeout. Reject 400 validation errors, 404 absence, decoding errors, and ordinary operation timeouts unless the follow-up probe also fails.

```go
func TestClassifyLiveServerHealthFailure(t *testing.T) {
	tests := []struct { name string; err error; want bool }{
		{"unavailable", &client.APIError{StatusCode: 503, Code: "SERVICE_UNAVAILABLE"}, true},
		{"validation", &client.APIError{StatusCode: 400, Code: "VALIDATION_ERROR"}, false},
		{"not found", &client.APIError{StatusCode: 404, Code: "NOT_FOUND"}, false},
	}
	for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { require.Equal(t, tt.want, classifyLiveServerHealthFailure(tt.err)) }) }
}
```

- [ ] **Step 2: Run tests and verify RED**

Run: `go test ./provider -run 'TestClassifyLiveServerHealthFailure|TestVerifyLiveServerHealth' -count=1`

Expected: compilation failure because the classifier and probe do not exist.

- [ ] **Step 3: Implement bounded health handling**

Implement the classifier without inspecting or returning raw response bodies. Implement `verifyLiveServerHealth` as a thin bounded probe wrapper. At heavy-tier boundaries, probe a low-cost authenticated endpoint already used by the provider. After a potentially health-related heavy failure, probe once; call `recordServerHealthFailure` only when classification or the probe confirms unsafe state.

Do not convert ordinary resource failures into skips. Preserve the original test failure after recording the stop marker.

- [ ] **Step 4: Verify stop behavior and safety**

Run: `go test ./provider -run 'Test(ClassifyLiveServerHealthFailure|VerifyLiveServerHealth|ServerHealthFailure|OrdinaryLiveResult)' -count=1`

Expected: PASS, with stop marker created only for confirmed health failures.

Run: `env -u DOKPLOY_ACCEPTANCE go test ./provider -run 'TestLiveTier(2|3|4)' -count=1 -v`

Expected: all live tiers SKIP and no probe is attempted.

- [ ] **Step 5: Commit**

```bash
git add provider/live_harness_test.go provider/live_harness_unit_test.go provider/live_workloads_test.go provider/live_databases_test.go provider/live_backups_test.go
git commit -m "test: stop acceptance after server health failures"
```

### Task 7: Preserve Workflow Artifacts On Failure

**Files:**
- Modify: `.github/workflows/run-acceptance-tests.yml:174-195`
- Modify: `provider/registry_metadata_test.go:840-910,1002-1100`
- Test: `provider/registry_metadata_test.go`

**Interfaces:**
- Consumes: acceptance step outcomes and existing workflow parser helpers.
- Produces: unconditional safe archive/upload behavior after aggregated failures.

- [ ] **Step 1: Add failing workflow contract assertions**

Require `Tar provider binaries` and `Upload artifacts` to use `if: ${{ always() }}`. Require archive creation to test for the provider binary and to create a non-secret status artifact when unavailable. Require upload to warn rather than fail when no archive exists, preserving the original acceptance result.

```go
tarStep := steps[findStepIndex(steps, "Tar provider binaries")].(map[string]any)
require.Equal(t, "${{ always() }}", tarStep["if"])
uploadStep := steps[findStepIndex(steps, "Upload artifacts")].(map[string]any)
require.Equal(t, "${{ always() }}", uploadStep["if"])
require.Equal(t, "warn", uploadStep["with"].(map[string]any)["if-no-files-found"])
```

- [ ] **Step 2: Run contract test and verify RED**

Run: `go test ./provider -run TestOwnedWorkflow -count=1`

Expected: FAIL because artifact steps are currently suppressed after the reporting step fails.

- [ ] **Step 3: Update the workflow**

Set both steps to `if: ${{ always() }}`. Make archive creation conditional on an executable `bin/pulumi-resource-dokploy`; write only a fixed `provider binary unavailable` marker when absent. Configure upload with `if-no-files-found: warn` and include the archive/status path without printing directory contents.

- [ ] **Step 4: Verify GREEN**

Run: `go test ./provider -run 'TestOwnedWorkflow|TestRegistryMetadata' -count=1`

Expected: PASS; final aggregation still fails on failed/cancelled test outcomes.

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/run-acceptance-tests.yml provider/registry_metadata_test.go
git commit -m "ci: retain acceptance artifacts after failures"
```

### Task 8: Strengthen Pulumi Smoke Evidence

**Files:**
- Modify: `tests/acceptance_test.go`
- Modify: `tests/acceptance_program_test.go`
- Test: `tests/acceptance_test.go`

**Interfaces:**
- Produces: `assertPreviewSummary(t *testing.T, summary auto.PreviewResult, phase string)` and `assertUpdateSummary(t *testing.T, summary auto.UpResult, phase string)` or equivalent helpers using the exact Automation API result types present in the pinned Pulumi SDK.
- Consumes: `auto.Stack.Preview`, `auto.Stack.Up`, `auto.Stack.Export`, `auto.Workspace.ListStacks`, and existing output assertions.

- [ ] **Step 1: Inspect pinned Automation API result fields and write failing pure summary tests**

Use the SDK types in `go.mod` to identify exact change-summary fields. Add pure tests with synthetic summaries proving revision one requires four managed creates and revision two allows updates but rejects replace/delete operations.

```go
func TestLifecycleSummaryRejectsReplacement(t *testing.T) {
	err := validateLifecycleChanges("revision two", map[string]int{"update": 3, "replace": 1})
	if err == nil || !strings.Contains(err.Error(), "replacement") {
		t.Fatalf("validateLifecycleChanges() = %v, want replacement error", err)
	}
}
```

- [ ] **Step 2: Run summary tests and verify RED**

Run: `go test ./tests -run 'TestLifecycleSummary' -count=1`

Expected: compilation failure because the summary validator does not exist.

- [ ] **Step 3: Assert live preview/up summaries**

Capture both Preview and Up results. Revision one must include the expected four custom resources without replacement/deletion. Revision two must include updates and no replacement/deletion. Keep stable-ID output assertions. Normalize Pulumi operation names in one helper rather than spreading string assumptions through the test.

- [ ] **Step 4: Add supported destroy and removal evidence**

Change cleanup callbacks to capture destroy result state where supported. After destroy, export stack state and assert no `dokploy:index:*` resources remain. After `RemoveStack`, call `ListStacks` and assert `stackName` is absent. Keep independent two-minute contexts and do not read backend files directly.

Add deterministic tests for stack-name absence and destroy-state filtering.

- [ ] **Step 5: Verify GREEN**

Run: `go test ./tests -run 'TestLifecycleSummary|TestLifecycleSmokeProgram|TestLifecycleSmokeCleanup' -count=1`

Expected: PASS without requiring Pulumi CLI or live credentials.

Run: `go test ./tests -run TestAccLifecycleSmoke -count=1 -v`

Expected: SKIP without `DOKPLOY_ACCEPTANCE=1`.

- [ ] **Step 6: Commit**

```bash
git add tests/acceptance_test.go tests/acceptance_program_test.go
git commit -m "test: strengthen Pulumi lifecycle evidence"
```

### Task 9: Static Verification And Review

**Files:**
- Review: `provider/application_source_test.go`
- Review: `provider/compose_source_test.go`
- Review: `provider/mount_lifecycle_matrix_test.go`
- Review: `provider/live_control_plane_test.go`
- Review: `provider/live_workloads_test.go`
- Review: `provider/live_databases_test.go`
- Review: `provider/live_backups_test.go`
- Review: `provider/live_harness_test.go`
- Review: `provider/live_harness_unit_test.go`
- Review: `provider/domain_test.go`
- Review: `provider/destination_test.go`
- Review: `provider/registry_test.go`
- Review: `tests/acceptance_test.go`
- Review: `tests/acceptance_program_test.go`
- Review: `.github/workflows/run-acceptance-tests.yml`

**Interfaces:**
- Consumes: Tasks 1-8.
- Produces: reviewed suite ready for cautious live execution.

- [ ] **Step 1: Format changed Go files**

Run: `gofmt -w provider/application_source_test.go provider/compose_source_test.go provider/mount_lifecycle_matrix_test.go provider/live_control_plane_test.go provider/live_workloads_test.go provider/live_databases_test.go provider/live_backups_test.go provider/live_harness_test.go provider/live_harness_unit_test.go provider/domain_test.go provider/destination_test.go provider/registry_test.go tests/acceptance_test.go tests/acceptance_program_test.go`

Expected: command exits zero.

- [ ] **Step 2: Run focused regression tests**

Run: `go test ./provider -run 'Test(ApplicationSource|ComposeSource|Mount|Domain|Destination|Registry|LiveGate|ClassifyLiveServerHealth|OwnedWorkflow)' -count=1`

Expected: PASS.

Run: `go test ./tests -run 'TestLifecycle|TestPulumiCLI' -count=1`

Expected: PASS; live smoke skips without opt-in.

- [ ] **Step 3: Run full static checks**

Run: `go test -short -count=1 ./provider/... ./internal/... ./tests/...`

Expected: PASS with live tests skipped.

Run: `go test -race ./provider/... ./internal/...`

Expected: PASS.

Run: `golangci-lint run`

Expected: PASS.

Run: `git diff --check`

Expected: PASS.

- [ ] **Step 4: Request code review**

Dispatch a reviewer against `docs/superpowers/specs/2026-09-09-live-acceptance-coverage-hardening-design.md`. Require findings on source field completeness, Mount hybrid coverage, cleanup ownership, health-stop false positives, secret-safe diagnostics, workflow failure semantics, and Pulumi summary assumptions.

- [ ] **Step 5: Apply accepted findings with TDD**

For each Critical or Important finding, add the smallest failing deterministic test, run it to observe failure, implement the minimal correction, and rerun the focused test plus all commands from Steps 2-3.

- [ ] **Step 6: Commit review fixes**

```bash
git add provider tests .github/workflows/run-acceptance-tests.yml
git commit -m "test: harden acceptance coverage review fixes"
```

### Task 10: Cautious Live Verification And Reporting

**Files:**
- Modify: `docs/bugs/2026-09-05-live-acceptance-run.md`
- Create conditionally: `docs/bugs/2026-09-09-<resource>-<operation>.md`

**Interfaces:**
- Consumes: reviewed test suite, built local provider, pinned Pulumi CLI, protected environment credentials.
- Produces: sanitized live evidence and reports for confirmed provider defects only.

- [ ] **Step 1: Build and preflight without printing secrets**

Run: `make provider`

Expected: `bin/pulumi-resource-dokploy` exists and is executable.

Run a shell preflight that sources protected local credentials, sets `DOKPLOY_ACCEPTANCE=1`, verifies required values and the optional prerequisite names without printing values, removes the configured stop marker, and prints only `acceptance prerequisites checked`.

- [ ] **Step 2: Run high-priority Tier 2 coverage first**

Run: `PATH="$PWD/bin:$PATH" mise exec -- go test ./provider -run '^TestLiveTier2Workloads$' -parallel=1 -count=1 -v`

Expected: Application/Compose source cases and Mount cases pass or produce sanitized failures; registry, GitLab, and custom resolver cases explicitly skip when prerequisites are absent; cleanup succeeds and the stop marker remains absent.

- [ ] **Step 3: Run focused Tier 1 coverage**

Run: `PATH="$PWD/bin:$PATH" mise exec -- go test ./provider -run '^TestLiveTier1ControlPlane$' -parallel=1 -count=1 -v`

Expected: configured cases pass or fail with sanitized diagnostics; all created resources are absent afterward.

- [ ] **Step 4: Run shared-harness regression tiers when safe**

If the stop marker remains absent, run:

```bash
PATH="$PWD/bin:$PATH" mise exec -- go test ./provider -run '^TestLiveTier3Databases$' -parallel=1 -count=1 -v
PATH="$PWD/bin:$PATH" mise exec -- go test ./provider -run '^TestLiveTier4Backups$' -parallel=1 -count=1 -v
```

Expected: serial execution; no backup execution; all cleanup verified. Stop immediately if the marker appears.

- [ ] **Step 5: Run the Pulumi smoke**

Run: `PATH="$PWD/bin:$PATH" DOKPLOY_ACCEPTANCE=1 mise exec -- go test ./tests -run '^TestAccLifecycleSmoke$' -parallel=1 -count=1 -v`

Expected: preview/up summaries match the lifecycle contract, IDs remain stable, and destroy/removal evidence passes.

- [ ] **Step 6: Classify failures and update reports**

Reproduce each failure once with the smallest matching subtest after confirming cleanup and server health. Record provider defects, server defects, and environment limitations separately. Create standalone bug reports only for reproduced provider defects, using structural diagnostics and no raw payloads or sensitive values.

Update `docs/bugs/2026-09-05-live-acceptance-run.md` with the new date/revision section, exact commands, result, duration, prerequisite skips, hybrid Mount limitation, cleanup evidence, and stop-marker status.

- [ ] **Step 7: Final verification and commit evidence**

Run: `go test -short -count=1 ./provider/... ./internal/... ./tests/...`

Expected: PASS.

Run: `git diff --check`

Expected: PASS.

```bash
git add docs/bugs
git commit -m "docs: record hardened acceptance coverage results"
```
