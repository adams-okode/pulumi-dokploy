# Live Acceptance README Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a clear operator guide for the live acceptance strategy and protect its required content with a deterministic test.

**Architecture:** Put the operator guide next to the Pulumi acceptance tests at `tests/README.md`. Link to it from `CONTRIBUTING.md`, and use a Go content test to protect required commands, safety rules, variables, test layers, and report locations.

**Tech Stack:** Markdown, Go 1.26.6, `testing`, standard library file and string functions.

## Global Constraints

- Follow `docs/superpowers/specs/2026-09-10-live-acceptance-readme-design.md`.
- Use short sentences and active voice.
- Give one instruction in each numbered step.
- Use the same term for the same item throughout the document.
- Define abbreviations before use.
- Avoid idioms, informal phrases, and unnecessary words.
- Use explicit warnings for unsafe actions.
- Do not mention the writing standard used for the document.
- Do not include real credentials, endpoints, resource IDs, or `.env` data.
- Do not tell Go test code to read `.env`.

---

### Task 1: Write The Live Acceptance Operator Guide

**Files:**
- Create: `tests/README.md`
- Modify: `tests/acceptance_test.go`
- Modify: `CONTRIBUTING.md:9-10`
- Test: `tests/acceptance_test.go`

**Interfaces:**
- Consumes: live tier names and commands from `.github/workflows/run-acceptance-tests.yml`, safety rules from `provider/live_harness_test.go`, and report rules from `docs/bugs/README.md`.
- Produces: `TestLiveAcceptanceReadmeContract(t *testing.T)` and the primary operator guide at `tests/README.md`.

- [ ] **Step 1: Write the failing README contract test**

Add a test that reads `README.md` from the `tests` package directory and `../CONTRIBUTING.md`. Require the two layers, all four provider tiers, required variables, serial flags, stop marker, optional prerequisites, warnings, and report paths. Reject language that tells Go test code to load `.env`.

```go
func TestLiveAcceptanceReadmeContract(t *testing.T) {
	readmeBytes, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	readme := string(readmeBytes)
	for _, required := range []string{
		"Direct-provider lifecycle tests",
		"Pulumi Automation API smoke test",
		"TestLiveTier1ControlPlane",
		"TestLiveTier2Workloads",
		"TestLiveTier3Databases",
		"TestLiveTier4Backups",
		"TestAccLifecycleSmoke",
		"DOKPLOY_ACCEPTANCE=1",
		"DOKPLOY_ENDPOINT",
		"DOKPLOY_API_KEY",
		"DOKPLOY_ACCEPTANCE_STOP_FILE",
		"-parallel=1",
		"DOKPLOY_REGISTRY_URL",
		"DOKPLOY_GITLAB_INTEGRATION_ID",
		"DOKPLOY_ACCEPTANCE_ALLOW_REPLICAS",
		"DOKPLOY_CUSTOM_CERT_RESOLVER",
		"docs/bugs/README.md",
		"docs/bugs/2026-09-05-live-acceptance-run.md",
	} {
		if !strings.Contains(readme, required) {
			t.Errorf("live acceptance README is missing %q", required)
		}
	}
	for _, prohibited := range []string{"Go test code loads .env", "tests read .env", "parse .env"} {
		if strings.Contains(readme, prohibited) {
			t.Errorf("live acceptance README contains prohibited instruction %q", prohibited)
		}
	}
	contributingBytes, err := os.ReadFile("../CONTRIBUTING.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contributingBytes), "tests/README.md") {
		t.Error("CONTRIBUTING.md does not link to the live acceptance guide")
	}
}
```

Add `os` and `strings` to the existing import list.

- [ ] **Step 2: Run the test and verify RED**

Run: `go test ./tests -run TestLiveAcceptanceReadmeContract -count=1`

Expected: FAIL because `tests/README.md` does not exist.

- [ ] **Step 3: Write `tests/README.md`**

Use this section order and include the listed content:

```markdown
# Live acceptance tests

## Purpose

State that these tests use a real Dokploy server and can create and delete resources.
Add a warning not to use a production server.

## Test layers

Define direct-provider lifecycle tests and the Pulumi Automation API smoke test.
List the responsibility of each layer.

## Safety rules

Require serial execution, one heavy operation, disabled backups, reserved test hosts,
bounded cleanup, no pre-existing resource changes, and no secret output.

## Prerequisites

List Go, the pinned Pulumi CLI, the local provider binary, required Dokploy variables,
and the stop-marker variable.

## Local setup

Give separate numbered steps for `mise install`, `make provider`, PATH configuration,
credential export, opt-in, and stop-marker preparation. State that the shell can source
a protected file, but test code must not read `.env`.

## Test commands

Give exact commands for Tier 1, Tier 2, Tier 3, Tier 4, and the Pulumi smoke. Each command
must use `-parallel=1 -count=1 -v`. Tell the operator to run them in order and inspect the
stop marker after each tier.

## Optional coverage

List Registry, GitLab, MongoDB replica, custom certificate resolver, and optional server
scope variables. State the skip behavior when a prerequisite is absent.

## Cleanup and stop behavior

Explain fallback cleanup, explicit absence checks, independent cleanup contexts, and
the conditions that create `DOKPLOY_ACCEPTANCE_STOP_FILE`. Tell the operator not to run
a later heavy tier when the marker exists.

## Result classification

Define pass, skip, provider defect, Dokploy server defect, and test environment limitation.
State that an ordinary resource failure does not stop an independent tier.

## Reports

Link to `docs/bugs/README.md` and `docs/bugs/2026-09-05-live-acceptance-run.md`. Require
sanitized evidence, durations, skip reasons, and cleanup status.
```

Write complete prose for every section. Do not leave the instructional notes in the final README.

- [ ] **Step 4: Link the guide from `CONTRIBUTING.md`**

Replace the existing acceptance sentence with:

```markdown
Pull requests must explain user-visible behavior and include regression tests.
Maintainers run live acceptance tests manually with protected repository secrets.
See [`tests/README.md`](tests/README.md) for the test strategy and safe commands.
```

- [ ] **Step 5: Run focused tests and verify GREEN**

Run: `go test ./tests -run 'TestLiveAcceptanceReadmeContract|TestLifecycleSmokeProgram|TestRequirePulumiCLI' -count=1`

Expected: PASS.

Run: `git diff --check`

Expected: PASS.

- [ ] **Step 6: Review the language**

Check each sentence in `tests/README.md`. Confirm that it has one clear meaning, uses active voice when possible, and uses the terms `live acceptance`, `provider tier`, `Pulumi smoke`, `stop marker`, and `cleanup` consistently. Confirm that the document does not name the writing standard.

- [ ] **Step 7: Commit**

```bash
git add tests/README.md tests/acceptance_test.go CONTRIBUTING.md
git commit -m "docs: explain live acceptance testing"
```
