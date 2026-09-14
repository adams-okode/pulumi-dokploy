# Pulumi Registry Readiness Ledger

Date: 2026-09-14

Pulumi Dokploy is **not yet Registry-ready**. Repository evidence and public
package availability are separated below from external actions that remain
pending. Every external item must have evidence before Registry readiness can
be claimed.

## Repository-complete

- Registry overview follows the current six-language structure and passes
  contract/documentation checks in `docs/_index.md`.
- Public whale logo, attribution, and generated `logoUrl` metadata are present.
- Private security and conduct reports use `contact@dimeski.net`.
- `.github/CODEOWNERS` assigns repository ownership to `@dimeskigj`.
- The publication runbook records all external actions without claiming
  completion.
- Existing User-Agent, plugin download metadata, keywords, license, schema
  compatibility, checksums, SBOMs, attestations, and `release-smoke` workflow
  remain covered.

## Publicly available

- Stable GitHub release `v0.2.2` exists with provider archives and checksums.
- The dated Task 5 anonymous public-package probes are recorded below.
- Anonymous availability is not runtime verification and does not replace a
  successful clean-cache `release-smoke` dispatch.

## External pending

- [ ] Successful `release-smoke` dispatch for `0.2.2` with all six jobs.
- [ ] Upstream `community-packages/package-list.json` entry.
- [ ] Maintainer-approved publisher display name and upstream
  `publisher-names.json` mapping to `dimeskigj`.
- [ ] Registry `fact-sheet`, `/check`, `/preview`, six-language and logo preview
  inspection.
- [ ] Pulumi Registry CODEOWNER approval, merge, deployment, and public page
  verification.

## Verification evidence

- `go test ./provider -run '^TestRegistryReadinessLedgerSeparatesEvidenceStates$' -count=1`
  passed after the ledger rewrite and pending-state strengthening; the required
  red runs failed first for the missing headings and then for missing unchecked
  markers.
- `go test ./provider -run '^TestBuildDotnetCreatesVersionFileForCleanCheckout$' -count=1`
  passed after the clean-checkout build fix.
- `go test ./provider -run 'Registry|SchemaPublishingMetadata|GeneratedPublishingMetadata' -count=1`,
  `make test_provider`, `env -i PATH="$PATH" HOME="$HOME" go test -race
  ./provider/... ./internal/...`, `PATH="/tmp/pinned-rsvg-bin:$PATH" make
  check_codegen`, `make check_openapi`, `make docs_check`, `mise exec
  golangci-lint@2.9.0 -- golangci-lint run`, `make govulncheck`, and `make
  license` passed.
- `make docs_check` reported three existing high-severity npm advisories and
  the existing missing `website/src/icons` warning; these remain out of scope.
- `make test_race` was also run with the inherited acceptance `.env` and reached
  live Domain operations that returned `BAD_REQUEST`; it was not rerun against
  external services. The clean-environment pinned race equivalent above passed
  with all live tests skipped.
- `mise run setup-svg-renderer` could not install the package because sudo
  requires an interactive password. An exact Ubuntu Noble
  `librsvg2-bin=2.58.0+dfsg-1build1` amd64 package (SHA256
  `84e6dc1615a63d202ae67699f8fb0ead13c8f4fc2664353cbed2fb662a32a4b4`) and
  matching `librsvg2-2`/`libxml2`/`libicu74` runtime libraries were extracted
  without installation; its `rsvg-convert --version` reported exactly 2.58.0.
  With that temporary pinned PATH, `make check_codegen` passed with no drift.
- `make build_sdks` passed all Go, Python, .NET, and Java stages when run with
  mise Temurin 11.0.32.1 and the exact Gradle 7.6 distribution. The .NET
  version file is now created deterministically from `VERSION_GENERIC`.
- `make lint` could not run because the unwrapped `golangci-lint` executable is
  unavailable; the pinned v2.9.0 equivalent above passed.
- The publication runbook remains unchecked; no release-smoke workflow was
  dispatched, and no upstream Registry files or Registry page were modified.

## Residual risks and limitations

- `release-smoke` has not run; its Pulumi CLI downloads are not
  checksum-verified, and Java package/runtime behavior remains unverified.
- Broad engine acceptance remains deferred and is not a Registry blocker.
- Existing website npm advisories and the missing `src/icons` warning remain
  outside this readiness task.

## Public v0.2.2 availability probes

Date: 2026-09-14

- `GitHub release and provider artifacts` - PASS - `https://api.github.com/repos/dimeskigj/pulumi-dokploy/releases/tags/v0.2.2` - HTTP success; JSON reported `tag_name` `v0.2.2`, `draft` `false`, and `prerelease` `false`, with `checksums.txt`, Linux/macOS/Windows provider archives for amd64 and arm64, and a `.sbom.json` asset for each archive.
- `npm` - PASS - `npm view @dimeskigj/pulumi-dokploy@0.2.2 version --json` - observed `"0.2.2"`.
- `PyPI` - PASS - `https://pypi.org/pypi/pulumi-dokploy/0.2.2/json` - HTTP success; JSON `info.version` was `0.2.2`.
- `NuGet` - PASS - `https://api.nuget.org/v3-flatcontainer/dimeskigj.pulumi.dokploy/0.2.2/dimeskigj.pulumi.dokploy.nuspec` - HTTP success; XML contained `<version>0.2.2</version>`.
- `Maven Central` - PASS - `https://repo1.maven.org/maven2/net/dimeski/pulumi/dokploy/0.2.2/dokploy-0.2.2.pom` - HTTP success; POM contained group `net.dimeski.pulumi`, artifact `dokploy`, and version `0.2.2`.
- `Go proxy` - PASS - `https://proxy.golang.org/github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy/@v/v0.2.2.info` - HTTP success; JSON `Version` was `v0.2.2`.

These anonymous probes establish public package availability only. They do not
replace the clean-cache installation and runtime evidence from a successful
`release-smoke` dispatch.
