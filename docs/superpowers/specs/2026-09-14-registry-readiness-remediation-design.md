# Registry Readiness Remediation Design

## Goal

Resolve every uncovered Pulumi Registry readiness item above low priority while
leaving workflow dispatches, upstream Registry changes, pull requests, and
approvals to a maintainer runbook.

The target release is `v0.2.2`. Repository and provider coordinates remain:

- Repository: `dimeskigj/pulumi-dokploy`
- Provider name: `dokploy`
- Publisher: `dimeskigj`
- Maintainer contact: `contact@dimeski.net`

## Scope

This work updates the Registry overview, provider logo metadata, repository
security and governance contacts, readiness tests, public release evidence, and
maintainer documentation. It covers the six languages currently supported and
documented by the provider: TypeScript, Python, Go, C#, Java, and YAML.

This work does not dispatch GitHub workflows, modify `pulumi/registry`, open a
pull request, publish or republish packages, or represent an external Registry
check as complete.

Lookup functions, exhaustive Pulumi-engine lifecycle acceptance, upstream
Dokploy changes, and unrelated documentation-site work remain out of scope.

## Registry Overview

Rewrite `docs/_index.md` to follow the current Pulumi Registry overview-page
guidance while retaining the community-maintained disclaimer.

The file must begin at byte zero with YAML front matter containing:

- `layout: package`
- `title: Dokploy`, matching the provider schema display name
- A concise `meta_desc` naming Dokploy and Pulumi

The body must not contain a second level-one heading. It will open with a short
description of the provider and link to Dokploy, then contain these top-level
sections in order:

1. `## Installation`
2. `## Example Usage`
3. `## Configuration`

Installation will use Pulumi's language chooser shortcodes with the exact six
language keys `typescript,python,go,csharp,java,yaml`. Each tab will contain the
current package command or dependency declaration. Java will show Maven and
Gradle coordinates. YAML will use `pulumi package add`.

Example Usage will use the same six-language chooser and provide a complete,
minimal Project resource program in each language. The section will also give
the `pulumi config set dokploy:endpoint` and
`pulumi config set --secret dokploy:apiKey` commands needed to run the examples.

Configuration will document every provider configuration parameter by bare
name, requiredness, secrecy, environment-variable fallback, and valid value:

- `endpoint`: required, not secret, with `DOKPLOY_ENDPOINT`
- `apiKey`: required, secret, with `DOKPLOY_API_KEY`

`docs/installation-configuration.md` remains a supplementary page. Its package
coordinates, configuration instructions, and maintenance language must remain
consistent with the overview. It is not a substitute for required overview
content.

## Logo And Attribution

Select an existing whale icon distributed under a permissive license such as
MIT, Apache-2.0, BSD, or CC0. Do not use Docker artwork or a close imitation of
Docker's trademarked whale-and-containers mark.

Store a self-contained SVG in the repository under a stable website asset path.
Record the original project, source URL, author if supplied, exact license,
copyright notice if supplied, and any modifications in an adjacent attribution
file. Preserve notices required by the source license.

Set `schema.Metadata.LogoURL` to the raw GitHub URL for that SVG on the default
branch. Regenerate the schema and all generated SDK metadata through project
generators rather than editing generated files manually. The URL must be public,
HTTPS, and return an SVG without authentication.

The Registry runbook must require confirmation that the upstream preview renders
the logo correctly; repository tests cannot establish this external result.

## Contacts And Governance

Update `SECURITY.md` to direct private vulnerability reports to
`contact@dimeski.net`. It must continue to prohibit disclosure of credentials in
public issues and must not imply that GitHub issues are a private reporting
channel.

Update the Code of Conduct enforcement section to use `contact@dimeski.net`
instead of Pulumi's address. The text must identify the recipient as this
project's maintainer rather than Pulumi's project team.

Add `.github/CODEOWNERS` with `@dimeskigj` owning the repository. This documents
local ownership but does not replace the Pulumi Registry repository's required
CODEOWNER approval.

## Public Release Verification

Perform read-only availability probes for exact release `v0.2.2` against:

- GitHub Releases provider archives and `checksums.txt`
- npm package `@dimeskigj/pulumi-dokploy` version `0.2.2`
- PyPI project `pulumi-dokploy` version `0.2.2`
- NuGet package `Dimeskigj.Pulumi.Dokploy` version `0.2.2`
- Maven Central artifact `net.dimeski.pulumi:dokploy:0.2.2`
- Go module `github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy@v0.2.2`

Each probe must be attributable to one ecosystem and must not install packages,
publish artifacts, access repository secrets, or mutate external state. Record
the date, endpoint or command, and result in the readiness ledger. Availability
proves publication only; it does not replace clean-cache installation and
runtime verification by `.github/workflows/release-smoke.yml`.

## Maintainer Runbook

Add a Registry publication runbook under `docs/` containing sequential manual
steps and explicit evidence to retain:

1. Dispatch `release-smoke` with input `0.2.2` and require every provider and
   SDK job to pass.
2. Fork or check out `pulumi/registry` and add this exact package-list entry:

   ```json
   {
     "repoSlug": "dimeskigj/pulumi-dokploy",
     "schemaFile": "provider/cmd/pulumi-resource-dokploy/schema.json"
   }
   ```

3. Add a new publisher mapping whose key is the maintainer-approved public
   display name and whose value is `dimeskigj`. The runbook must clearly require
   replacing the descriptive display-name marker with the maintainer's chosen
   value before submission; it must not invent a personal name.
4. Open the Registry pull request and wait for the automated fact-sheet.
5. Resolve flagged checks and use `/check` to rerun them.
6. Ask a Pulumi maintainer for `/preview`; inspect all six language tabs,
   generated resource pages, configuration text, links, and logo rendering.
7. Obtain Registry CODEOWNER approval and merge through the upstream process.
8. Verify the public Registry page after deployment.

The runbook uses unchecked boxes for external actions. It must distinguish
provider-repository CODEOWNERS from upstream Registry CODEOWNER approval and
must not claim that a Registry PR exists.

## Tests

Extend provider contract tests to validate repository-owned behavior:

- Required front matter starts at byte zero and matches schema display name.
- No body H1 exists.
- Required sections exist in the required order.
- Installation and Example Usage choosers use the exact Pulumi delimiters and
  have balanced tabs for all six language keys.
- Every language contains its expected package coordinate and a complete
  Project example.
- Endpoint and secret API-key commands are present.
- Both configuration parameters document requiredness, secrecy, and environment
  fallbacks.
- Schema and generated metadata contain the exact public `logoUrl`.
- The logo asset and attribution record exist and identify a permitted license.
- Security and conduct files use `contact@dimeski.net` and no longer direct
  reports to Pulumi.
- `.github/CODEOWNERS` assigns `@dimeskigj`.
- The runbook contains the exact upstream package entry, publisher mapping
  instructions, smoke dispatch, fact-sheet, preview, logo check, approval, and
  post-deployment verification without marking them complete.

Run focused tests first, then provider tests, documentation checks, schema and
SDK drift checks, SDK builds, and `git diff --check`. Use the pinned project
toolchain where available and record any unavailable command and equivalent
verification honestly.

## Readiness Ledger

Update `docs/provider-registry-readiness.md` after verification. The ledger must
separate three states:

- Repository-complete: docs, logo metadata, contacts, CODEOWNERS, tests, and
  runbook are present and verified locally.
- Publicly available: exact `v0.2.2` provider and SDK coordinates responded to
  read-only probes.
- External pending: release-smoke dispatch, Registry package entry, publisher
  display-name decision and mapping, fact-sheet, preview, logo rendering,
  Registry CODEOWNER approval, merge, and public page verification.

The provider must not be called Registry-ready until every external pending item
has evidence. A failed ecosystem probe remains a release blocker and must be
reported without changing package versions or publishing credentials as part of
this work.

## Security And Error Handling

No verification command may print credentials or contact secrets. Public probes
must use anonymous registry endpoints. Documentation examples use invalid
example endpoints and placeholder API keys only. Logo selection must reject
ambiguous licensing, missing required attribution, or trademark-confusing
artwork rather than guessing.

Verification failures must identify the affected category: Registry docs,
metadata or asset, governance/contact, package ecosystem, or external maintainer
step. Static tests and public availability probes must never be described as
evidence that an external workflow, preview, approval, or deployment passed.
