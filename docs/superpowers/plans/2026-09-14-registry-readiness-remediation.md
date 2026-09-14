# Registry Readiness Remediation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Resolve every repository-owned Pulumi Registry readiness gap above low priority and provide an exact maintainer runbook for the external publication steps.

**Architecture:** Keep handwritten Registry content and governance files in the repository, generate schema and SDK metadata from `provider.Provider()`, and enforce the result with focused Go contract tests. Treat anonymous public-package probes as dated publication evidence only; keep workflow dispatch, upstream Registry changes, preview, approval, and deployment explicitly manual.

**Tech Stack:** Go 1.26.6, Pulumi Go Provider v1.6.0, Pulumi CLI 3.259.0, Markdown with Pulumi Hugo chooser shortcodes, SVG, GitHub Actions, npm/PyPI/NuGet/Maven/Go public APIs.

## Global Constraints

- Target release is exactly `v0.2.2`; package versions used by probes are exactly `0.2.2`.
- Repository remains `dimeskigj/pulumi-dokploy`, provider name remains `dokploy`, and publisher remains `dimeskigj`.
- Maintainer contact is `contact@dimeski.net`.
- Registry overview supports exactly TypeScript, Python, Go, C#, Java, and YAML; do not add HCL in this work.
- Select an existing whale icon under MIT, Apache-2.0, BSD, or CC0 terms; preserve required notices and reject Docker artwork or close trademark imitation.
- Do not dispatch workflows, modify `pulumi/registry`, open pull requests, publish packages, or claim external checks passed.
- Do not add lookup functions or broaden Pulumi-engine acceptance.
- Generated schema and SDK files must be regenerated, never manually edited.
- Public probes must be anonymous, read-only, separately attributable, and must not install packages.
- Never place real API keys, endpoints, package credentials, or secrets in docs, tests, commands, or logs.
- Existing unrelated worktree changes must remain untouched.

---

## File Structure

- `provider/registry_docs_test.go`: owns Registry overview, logo attribution, contact, CODEOWNERS, and maintainer-runbook contracts.
- `docs/_index.md`: complete Pulumi Registry landing page for six supported languages.
- `docs/installation-configuration.md`: supplementary installation and configuration page kept consistent with the overview.
- `website/public/logo.svg`: canonical public whale asset used by the Registry.
- `website/public/logo-LICENSE.md`: source, copyright, license, and modification record for the whale asset.
- `provider/provider.go`: source of `logoUrl` schema metadata.
- `provider/schema_test.go`: source and generated logo metadata contracts.
- `provider/cmd/pulumi-resource-dokploy/schema.json` and `sdk/`: generator-owned outputs.
- `SECURITY.md`: private vulnerability reporting instructions.
- `CODE-OF-CONDUCT.md`: project-specific enforcement contact.
- `.github/CODEOWNERS`: local repository ownership.
- `docs/registry-publication-runbook.md`: unchecked maintainer procedure for all external actions.
- `docs/provider-registry-readiness.md`: current readiness ledger and dated public probe evidence.

### Task 1: Registry Overview Documentation

**Files:**
- Modify: `provider/registry_docs_test.go:20-74`
- Modify: `docs/_index.md:1-45`
- Modify: `docs/installation-configuration.md:8-63`

**Interfaces:**
- Consumes: existing package coordinates and `Project` constructors from generated SDKs.
- Produces: a six-language Registry overview whose exact headings, shortcodes, package coordinates, examples, and configuration markers are consumed by Task 5's runbook and Task 6's readiness ledger.

- [ ] **Step 1: Replace the broad documentation test with exact overview structure tests**

In `provider/registry_docs_test.go`, retain `readProjectFile` and license tests, and replace the existing overview assertions with these helpers and tests:

```go
var registryLanguages = []string{"typescript", "python", "go", "csharp", "java", "yaml"}

func TestRegistryOverviewStructure(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	require.True(t, strings.HasPrefix(index, "---\n"))
	parts := strings.SplitN(index, "---\n", 3)
	require.Len(t, parts, 3)
	frontMatter := map[string]string{}
	require.NoError(t, yaml.Unmarshal([]byte(parts[1]), &frontMatter))
	require.Equal(t, "package", frontMatter["layout"])
	require.Equal(t, "Dokploy", frontMatter["title"])
	require.Contains(t, frontMatter["meta_desc"], "Dokploy")
	require.Contains(t, frontMatter["meta_desc"], "Pulumi")
	require.NotRegexp(t, regexp.MustCompile(`(?m)^# `), parts[2])

	headings := []string{"## Installation", "## Example Usage", "## Configuration"}
	previous := -1
	for _, heading := range headings {
		position := strings.Index(index, heading)
		require.Greater(t, position, previous, heading)
		previous = position
	}
}

func TestRegistryOverviewLanguageChoosers(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	chooser := `{{< chooser language "typescript,python,go,csharp,java,yaml" >}}`
	require.Equal(t, 2, strings.Count(index, chooser))
	require.Equal(t, 2, strings.Count(index, "{{< /chooser >}}"))
	for _, language := range registryLanguages {
		open := "{{% choosable language " + language + " %}}"
		require.Equal(t, 2, strings.Count(index, open), language)
	}
	require.Equal(t, 12, strings.Count(index, "{{% /choosable %}}"))
}

func TestRegistryOverviewCoordinatesExamplesAndConfiguration(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	for _, marker := range []string{
		"npm install @dimeskigj/pulumi-dokploy",
		"pip install pulumi-dokploy",
		"go get github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy",
		"dotnet add package Dimeskigj.Pulumi.Dokploy",
		"<groupId>net.dimeski.pulumi</groupId>",
		"implementation 'net.dimeski.pulumi:dokploy:",
		"pulumi package add github.com/dimeskigj/pulumi-dokploy dokploy",
		"new dokploy.Project", "pulumi_dokploy.Project", "dokploy.NewProject",
		"new Project", "new Project(\"example\"", "type: dokploy:index:Project",
		"pulumi config set dokploy:endpoint https://dokploy.example.invalid",
		"pulumi config set --secret dokploy:apiKey your-api-key",
		"`endpoint` (Required, Not secret)", "`apiKey` (Required, Secret)",
		"DOKPLOY_ENDPOINT", "DOKPLOY_API_KEY", "community-maintained",
	} {
		require.Contains(t, index, marker)
	}
	require.NotContains(t, index, "official Dokploy")
	require.NotContains(t, index, "official Pulumi")
}

func TestRegistryDocumentationCoordinatesStayConsistent(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	installation := readProjectFile(t, "../docs/installation-configuration.md")
	for _, marker := range []string{
		"@dimeskigj/pulumi-dokploy", "pulumi-dokploy",
		"Dimeskigj.Pulumi.Dokploy", "net.dimeski.pulumi:dokploy",
		"github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy",
		"dokploy:endpoint", "dokploy:apiKey", "DOKPLOY_ENDPOINT", "DOKPLOY_API_KEY",
	} {
		require.Contains(t, index, marker)
		require.Contains(t, installation, marker)
	}
}
```

Keep `TestRegistryFrontMatterAndSupportLinks`, but remove its duplicate index front-matter assertions already covered above. Keep `TestRegistryInstallationDetails`, `TestRegistryExampleDoesNotExposeSecrets`, and `TestRegistryLicense`, updating the example test in Step 4.

- [ ] **Step 2: Run focused tests to verify RED**

Run: `go test ./provider -run '^TestRegistry(Overview|DocumentationCoordinates)' -count=1`

Expected: FAIL because the overview has a body H1, no choosers, no `## Example Usage`, and only a TypeScript example.

- [ ] **Step 3: Rewrite the Registry overview**

Replace `docs/_index.md` with valid front matter and this exact shape:

```markdown
---
layout: package
title: Dokploy
meta_desc: Use the Pulumi Dokploy provider to manage Dokploy infrastructure and workloads.
description: Pulumi provider for managing Dokploy resources.
---

The Dokploy provider for Pulumi lets you manage [Dokploy](https://dokploy.com/)
projects, applications, databases, domains, backups, and related resources as
part of Pulumi programs. This is a **community-maintained** provider; neither
Dokploy nor Pulumi maintains this package.

## Installation

{{< chooser language "typescript,python,go,csharp,java,yaml" >}}
{{% choosable language typescript %}}

```bash
npm install @dimeskigj/pulumi-dokploy
```

{{% /choosable %}}
{{% choosable language python %}}

```bash
pip install pulumi-dokploy
```

{{% /choosable %}}
{{% choosable language go %}}

```bash
go get github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy
```

{{% /choosable %}}
{{% choosable language csharp %}}

```bash
dotnet add package Dimeskigj.Pulumi.Dokploy
```

{{% /choosable %}}
{{% choosable language java %}}

Maven:

```xml
<dependency>
  <groupId>net.dimeski.pulumi</groupId>
  <artifactId>dokploy</artifactId>
  <version>0.2.2</version>
</dependency>
```

Gradle:

```groovy
implementation 'net.dimeski.pulumi:dokploy:0.2.2'
```

{{% /choosable %}}
{{% choosable language yaml %}}

```bash
pulumi package add github.com/dimeskigj/pulumi-dokploy dokploy
```

{{% /choosable %}}
{{< /chooser >}}

## Example Usage
```

Continue with a second chooser using the same six keys. Add one complete Project program per tab using these language-specific constructors and exports:

```text
TypeScript: import @pulumi/pulumi and @dimeskigj/pulumi-dokploy; new dokploy.Project; export projectId = project.id.
Python: import pulumi and pulumi_dokploy; pulumi_dokploy.Project; pulumi.export("project_id", project.id).
Go: package main; import Pulumi and the Dokploy SDK; pulumi.Run; dokploy.NewProject; ctx.Export("projectId", project.ID()).
C#: Deployment.RunAsync; using Pulumi and Dimeskigj.Pulumi.Dokploy; new Project with ProjectArgs; return project.Id.
Java: Pulumi.run; import Context, Pulumi, Project, and ProjectArgs; new Project("example", ProjectArgs.builder()...build()); ctx.export("projectId", project.id()).
YAML: name, runtime yaml, a dokploy:index:Project resource, and projectId output.
```

Every example creates project name `registry-example` and description `Managed by Pulumi`; none embeds provider configuration or API keys. After the chooser, add:

```markdown
Configure the Dokploy endpoint and API key before running an example:

```bash
pulumi config set dokploy:endpoint https://dokploy.example.invalid
pulumi config set --secret dokploy:apiKey your-api-key
```

## Configuration

- `endpoint` (Required, Not secret) - The HTTP or HTTPS URL of the Dokploy instance. May also be set with `DOKPLOY_ENDPOINT`.
- `apiKey` (Required, Secret) - The API key used to authenticate to Dokploy. Set it with `pulumi config set --secret`; it may also be supplied with `DOKPLOY_API_KEY`.

Pulumi can acquire the matching provider plugin automatically from the GitHub
release metadata. See the [installation and configuration guide](./installation-configuration/)
for additional secret-handling guidance.

Support and source code are available in the [repository](https://github.com/dimeskigj/pulumi-dokploy).
Report non-sensitive problems in [GitHub issues](https://github.com/dimeskigj/pulumi-dokploy/issues)
and see the [contributing guide](https://github.com/dimeskigj/pulumi-dokploy/blob/main/CONTRIBUTING.md).
```

Use only ASCII punctuation in added prose. Confirm generated SDK constructor names against checked-in SDKs while writing each complete example.

- [ ] **Step 4: Align the supplementary page and secret-example test**

In `docs/installation-configuration.md`, change the Python command from `pip install pulumi_dokploy` to `pip install pulumi-dokploy`, retain the import-coordinate explanation, and ensure the package-add command remains `pulumi package add github.com/dimeskigj/pulumi-dokploy dokploy`.

Update `TestRegistryExampleDoesNotExposeSecrets` to inspect all content before the configuration command and reject credential configuration inside examples:

```go
func TestRegistryExamplesDoNotExposeSecrets(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	configuration := strings.Index(index, "Configure the Dokploy endpoint and API key")
	require.NotEqual(t, -1, configuration)
	examples := index[:configuration]
	require.NotContains(t, examples, "apiKey")
	require.NotContains(t, examples, "DOKPLOY_API_KEY")
	require.NotContains(t, examples, "your-api-key")
}
```

- [ ] **Step 5: Run focused tests to verify GREEN**

Run: `go test ./provider -run '^TestRegistry(Overview|DocumentationCoordinates|Examples)' -count=1`

Expected: PASS.

Run: `make docs_check`

Expected: PASS, including Markdown/Hugo shortcode checks and website build. Record existing unrelated npm advisories separately if emitted.

- [ ] **Step 6: Commit**

```bash
git add docs/_index.md docs/installation-configuration.md provider/registry_docs_test.go
git commit -m "docs: align registry overview guidelines"
```

### Task 2: Licensed Whale Logo And Metadata

**Files:**
- Create: `website/public/logo.svg`
- Create: `website/public/logo-LICENSE.md`
- Modify: `provider/schema_test.go:93-128`
- Modify: `provider/registry_docs_test.go`
- Modify: `provider/provider.go:20-30`
- Regenerate: `provider/cmd/pulumi-resource-dokploy/schema.json`
- Regenerate: `sdk/go`
- Regenerate: `sdk/nodejs`
- Regenerate: `sdk/python`
- Regenerate: `sdk/dotnet`
- Regenerate: `sdk/java`

**Interfaces:**
- Consumes: a selected whale SVG with a verified permissive upstream license and required notices.
- Produces: `const registryLogoURL = "https://raw.githubusercontent.com/dimeskigj/pulumi-dokploy/main/website/public/logo.svg"` in tests and matching `schema.Metadata.LogoURL` plus generated schema metadata.

- [ ] **Step 1: Select and record an eligible source before copying artwork**

Choose an existing whale SVG only after verifying its upstream source page and license text. Accept MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause, or CC0. Reject artwork if any of these are true:

```text
- The license is absent, custom, noncommercial, or ambiguous.
- The source lacks required copyright/author information.
- The icon depicts a whale carrying containers or is visually close to Docker's mark.
- The SVG embeds scripts, external URLs, raster images, fonts, or metadata containing personal data.
```

Record the selected source URL, project, author/copyright if supplied, SPDX license identifier, and whether the SVG will be recolored or simplified. This exact evidence will populate `website/public/logo-LICENSE.md`; do not proceed with an unverified source.

- [ ] **Step 2: Write failing logo contract tests**

Add to `provider/schema_test.go`:

```go
const registryLogoURL = "https://raw.githubusercontent.com/dimeskigj/pulumi-dokploy/main/website/public/logo.svg"
```

Change `TestSchemaPublishingMetadata` and `TestGeneratedPublishingMetadata` from `require.Empty(...LogoURL)` to:

```go
require.Equal(t, registryLogoURL, spec.LogoURL)
```

and:

```go
require.Equal(t, registryLogoURL, schemaMetadata.LogoURL)
```

Add to `provider/registry_docs_test.go`:

```go
func TestRegistryLogoAssetAndAttribution(t *testing.T) {
	logo := readProjectFile(t, "../website/public/logo.svg")
	attribution := readProjectFile(t, "../website/public/logo-LICENSE.md")
	require.Contains(t, logo, "<svg")
	require.Contains(t, logo, "viewBox=")
	for _, forbidden := range []string{"<script", "javascript:", "http://", "https://", "data:image/", "<image", "<foreignObject"} {
		require.NotContains(t, logo, forbidden)
	}
	for _, marker := range []string{"## Source", "## License", "## Modifications", "SPDX-License-Identifier:"} {
		require.Contains(t, attribution, marker)
	}
	require.Regexp(t, regexp.MustCompile(`SPDX-License-Identifier: (MIT|Apache-2\.0|BSD-2-Clause|BSD-3-Clause|CC0-1\.0)`), attribution)
}
```

- [ ] **Step 3: Run focused tests to verify RED**

Run: `go test ./provider -run 'Test(SchemaPublishingMetadata|GeneratedPublishingMetadata|RegistryLogoAssetAndAttribution)$' -count=1`

Expected: FAIL because source/generated metadata have no logo URL and the logo files do not exist.

- [ ] **Step 4: Add the sanitized SVG and attribution record**

Create `website/public/logo.svg` from the selected source. Keep only static SVG geometry, `viewBox`, accessibility title/description, and project-selected colors. Remove scripts, external references, embedded raster data, `<foreignObject>`, fonts, editor metadata, and dimensions that prevent responsive rendering.

Create `website/public/logo-LICENSE.md` with this concrete structure, substituting the verified facts recorded in Step 1 for the uppercase field labels. The committed file must contain none of these field labels:

```markdown
# Whale Logo Attribution

## Source

- Project: VERIFIED_UPSTREAM_PROJECT_NAME
- Original asset: VERIFIED_HTTPS_SOURCE_URL
- Author or copyright: VERBATIM_UPSTREAM_NOTICE_OR_CC0_NO_NOTICE_STATEMENT

## License

SPDX-License-Identifier: VERIFIED_PERMITTED_SPDX_IDENTIFIER

VERBATIM_LICENSE_ATTRIBUTION_OR_NOTICE

## Modifications

The SVG was sanitized for static web use. Record every visual change, or state
that no visual changes were made.
```

- [ ] **Step 5: Set provider logo metadata**

In `provider/provider.go`, add the field adjacent to `PluginDownloadURL`:

```go
LogoURL:           "https://raw.githubusercontent.com/dimeskigj/pulumi-dokploy/main/website/public/logo.svg",
```

- [ ] **Step 6: Verify source tests before generation**

Run: `go test ./provider -run 'Test(SchemaPublishingMetadata|RegistryLogoAssetAndAttribution)$' -count=1`

Expected: PASS.

Run: `go test ./provider -run '^TestGeneratedPublishingMetadata$' -count=1`

Expected: FAIL because checked-in generated schema still lacks `logoUrl`.

- [ ] **Step 7: Regenerate schema and all SDKs**

Run: `make VERSION_GENERIC=0.0.1-alpha.0+dev codegen`

Expected: PASS; generated schema gains `logoUrl`, and SDK outputs are regenerated through Pulumi CLI 3.259.0.

- [ ] **Step 8: Verify generated metadata and drift**

Run: `go test ./provider -run 'Test(SchemaPublishingMetadata|GeneratedPublishingMetadata|RegistryLogoAssetAndAttribution)$' -count=1`

Expected: PASS.

Run: `make check_codegen`

Expected: PASS with no additional generated diff.

Run: `make build_sdks`

Expected: PASS for Go, Python, Node.js, .NET, and Java.

- [ ] **Step 9: Commit**

```bash
git add website/public/logo.svg website/public/logo-LICENSE.md provider/provider.go provider/schema_test.go provider/registry_docs_test.go provider/cmd/pulumi-resource-dokploy/schema.json sdk
git commit -m "feat: add registry package logo"
```

### Task 3: Maintainer Contacts And Repository Ownership

**Files:**
- Modify: `provider/registry_docs_test.go`
- Modify: `SECURITY.md:1-8`
- Modify: `CODE-OF-CONDUCT.md:62-69`
- Create: `.github/CODEOWNERS`

**Interfaces:**
- Consumes: maintainer email `contact@dimeski.net` and GitHub owner `@dimeskigj`.
- Produces: private reporting instructions and repository ownership consumed by Task 5's distinction between local and upstream CODEOWNER approval.

- [ ] **Step 1: Write failing governance contract tests**

Add to `provider/registry_docs_test.go`:

```go
func TestRegistryMaintainerContactsAndOwnership(t *testing.T) {
	security := readProjectFile(t, "../SECURITY.md")
	conduct := readProjectFile(t, "../CODE-OF-CONDUCT.md")
	codeowners := readProjectFile(t, "../.github/CODEOWNERS")

	for name, content := range map[string]string{"SECURITY.md": security, "CODE-OF-CONDUCT.md": conduct} {
		require.Contains(t, content, "contact@dimeski.net", name)
		require.NotContains(t, content, "code-of-conduct@pulumi.com", name)
	}
	require.Contains(t, security, "do not report security vulnerabilities in public issues")
	require.Contains(t, security, "privately")
	require.Equal(t, "* @dimeskigj\n", codeowners)
}
```

- [ ] **Step 2: Run focused test to verify RED**

Run: `go test ./provider -run '^TestRegistryMaintainerContactsAndOwnership$' -count=1`

Expected: FAIL because neither policy uses the maintainer email and `.github/CODEOWNERS` is absent.

- [ ] **Step 3: Correct security reporting instructions**

Replace the opening instructions in `SECURITY.md` with:

```markdown
# Security Policy

Please do not report security vulnerabilities in public issues. Email the
maintainer privately at [contact@dimeski.net](mailto:contact@dimeski.net) with a
description, reproduction steps, and affected versions. The maintainer will
acknowledge reports and coordinate a fix and disclosure timeline.

Never include Dokploy API keys, SSH private keys, passwords, or other credentials
in issues, pull requests, examples, logs, or test fixtures.
```

- [ ] **Step 4: Correct Code of Conduct enforcement contact**

Replace `CODE-OF-CONDUCT.md:64-69` with:

```markdown
Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported privately to the project maintainer at
[contact@dimeski.net](mailto:contact@dimeski.net). All complaints will be
reviewed and investigated, and the maintainer will respond as appropriate to the
circumstances. The maintainer will keep the reporter's identity confidential.
Further enforcement policies may be posted separately.
```

- [ ] **Step 5: Add local ownership**

Create `.github/CODEOWNERS` with exactly:

```text
* @dimeskigj
```

- [ ] **Step 6: Run focused and provider tests**

Run: `go test ./provider -run '^TestRegistryMaintainerContactsAndOwnership$' -count=1`

Expected: PASS.

Run: `go test -short ./provider/... -count=1`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add SECURITY.md CODE-OF-CONDUCT.md .github/CODEOWNERS provider/registry_docs_test.go
git commit -m "docs: establish maintainer contacts"
```

### Task 4: Registry Publication Runbook

**Files:**
- Modify: `provider/registry_docs_test.go`
- Create: `docs/registry-publication-runbook.md`

**Interfaces:**
- Consumes: exact release `0.2.2`, package-list coordinates, local logo URL, maintainer contact, and six-language docs from Tasks 1-3.
- Produces: an unchecked, executable maintainer runbook consumed by the readiness ledger in Task 6.

- [ ] **Step 1: Write a failing runbook contract test**

Add to `provider/registry_docs_test.go`:

```go
func TestRegistryPublicationRunbook(t *testing.T) {
	runbook := readProjectFile(t, "../docs/registry-publication-runbook.md")
	for _, marker := range []string{
		"workflow run release-smoke.yml -f version=0.2.2",
		`"repoSlug": "dimeskigj/pulumi-dokploy"`,
		`"schemaFile": "provider/cmd/pulumi-resource-dokploy/schema.json"`,
		`"dimeskigj"`, "publisher-names.json", "maintainer-approved public display name",
		"/check", "/preview", "fact-sheet", "six language tabs", "logo",
		"Registry CODEOWNER", "public Registry page",
	} {
		require.Contains(t, runbook, marker)
	}
	require.GreaterOrEqual(t, strings.Count(runbook, "- [ ]"), 8)
	require.NotContains(t, runbook, "- [x]")
	require.NotContains(t, runbook, "Registry PR has been opened")
}
```

- [ ] **Step 2: Run focused test to verify RED**

Run: `go test ./provider -run '^TestRegistryPublicationRunbook$' -count=1`

Expected: FAIL because the runbook does not exist.

- [ ] **Step 3: Create the maintainer runbook**

Create `docs/registry-publication-runbook.md` with these sections and commands:

```markdown
# Pulumi Registry Publication Runbook

Target release: `v0.2.2`

This runbook performs external actions that repository tests cannot complete.
Keep every item unchecked until its evidence exists.

## 1. Verify Released Consumers

- [ ] Dispatch the read-only release smoke workflow:

```bash
gh workflow run release-smoke.yml -f version=0.2.2
```

- [ ] Record the run URL and require `provider`, `node`, `python`, `dotnet`,
  `java`, and `go` jobs to pass. A public-registry propagation failure is not a
  pass; rerun only after that exact package version is available.

## 2. Prepare The Registry Change

- [ ] Fork or check out `pulumi/registry` and add this object to
  `community-packages/package-list.json`:

```json
{
  "repoSlug": "dimeskigj/pulumi-dokploy",
  "schemaFile": "provider/cmd/pulumi-resource-dokploy/schema.json"
}
```

- [ ] Add a `publisher-names.json` entry whose key is the
  **maintainer-approved public display name** and whose value is exactly
  `"dimeskigj"`. The file is
  `tools/resourcedocsgen/pkg/publishers/publisher-names.json`. Decide and insert
  the public display name before submission; do not submit descriptive marker
  text or invent another person's name.

- [ ] Run the current lint/check commands documented by `pulumi/registry` and
  retain their output.

## 3. Open And Validate The Registry PR

- [ ] Open the upstream pull request and wait for its automated fact-sheet.
- [ ] Resolve every fact-sheet finding in the appropriate repository, then
  comment `/check` to rerun checks.
- [ ] Ask a Pulumi maintainer to comment `/preview`.
- [ ] Inspect the preview: all six language tabs, installation commands,
  examples, configuration, generated resource pages, links, and logo must render
  correctly.
- [ ] Obtain approval from a Pulumi Registry CODEOWNER. The local
  `.github/CODEOWNERS` file does not satisfy this upstream approval.

## 4. Verify Publication

- [ ] Merge only through the upstream maintainer process.
- [ ] After deployment, open the public Registry page and verify discovery,
  publisher display name, logo, overview, generated API documentation, and
  installation links.
- [ ] Update `docs/provider-registry-readiness.md` with the smoke run, pull
  request, approval, merge, and public page evidence.
```

- [ ] **Step 4: Run focused test to verify GREEN**

Run: `go test ./provider -run '^TestRegistryPublicationRunbook$' -count=1`

Expected: PASS.

Run: `git diff --check`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add docs/registry-publication-runbook.md provider/registry_docs_test.go
git commit -m "docs: add registry publication runbook"
```

### Task 5: Public v0.2.2 Package Probes

**Files:**
- Modify: `docs/provider-registry-readiness.md`

**Interfaces:**
- Consumes: anonymous public endpoints for release `v0.2.2` and exact SDK versions.
- Produces: dated, ecosystem-specific availability evidence for Task 6; does not produce scripts or claim clean-cache runtime success.

- [ ] **Step 1: Probe the GitHub release and provider artifacts**

Run:

```bash
gh api repos/dimeskigj/pulumi-dokploy/releases/tags/v0.2.2 --jq '{tag_name,draft,prerelease,assets:[.assets[].name]}'
```

Expected: `tag_name` is `v0.2.2`, `draft` and `prerelease` are false, and assets include `checksums.txt` plus Linux, macOS, and Windows provider archives and SBOMs.

- [ ] **Step 2: Probe npm without installing**

Run: `npm view @dimeskigj/pulumi-dokploy@0.2.2 version --json`

Expected: output is `"0.2.2"`.

- [ ] **Step 3: Probe PyPI without installing**

Run: `curl --fail --silent --show-error https://pypi.org/pypi/pulumi-dokploy/0.2.2/json`

Expected: HTTP success and JSON containing `"version":"0.2.2"` under `info`.

- [ ] **Step 4: Probe NuGet without installing**

Run: `curl --fail --silent --show-error https://api.nuget.org/v3-flatcontainer/dimeskigj.pulumi.dokploy/0.2.2/dimeskigj.pulumi.dokploy.nuspec`

Expected: HTTP success and XML containing `<version>0.2.2</version>`.

- [ ] **Step 5: Probe Maven Central without installing**

Run: `curl --fail --silent --show-error https://repo1.maven.org/maven2/net/dimeski/pulumi/dokploy/0.2.2/dokploy-0.2.2.pom`

Expected: HTTP success and POM containing group `net.dimeski.pulumi`, artifact `dokploy`, and version `0.2.2`.

- [ ] **Step 6: Probe the Go proxy without installing**

Run: `curl --fail --silent --show-error https://proxy.golang.org/github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy/@v/v0.2.2.info`

Expected: HTTP success and JSON containing `"Version":"v0.2.2"`.

- [ ] **Step 7: Record attributable results**

Add a `## Public v0.2.2 availability probes` section to `docs/provider-registry-readiness.md`. Use the actual execution date and one bullet per ecosystem in this exact format:

```markdown
- `<ecosystem>` - PASS/FAIL - `<endpoint or command>` - `<concise observed result>`.
```

If any probe fails, record `FAIL` and the HTTP/status error verbatim without secrets. Do not change versions, publish packages, or mark that ecosystem available. Add this statement after the list:

```markdown
These anonymous probes establish public package availability only. They do not
replace the clean-cache installation and runtime evidence from a successful
`release-smoke` dispatch.
```

- [ ] **Step 8: Inspect the evidence diff**

Run: `git diff --check`

Expected: PASS.

Run: `git diff -- docs/provider-registry-readiness.md`

Expected: only factual probe evidence has been added at this stage; no external workflow, Registry PR, preview, approval, or deployment is marked complete.

Do not commit yet; Task 6 completes and commits the readiness ledger as one coherent update.

### Task 6: Current Readiness Ledger And Final Verification

**Files:**
- Modify: `provider/registry_docs_test.go`
- Modify: `docs/provider-registry-readiness.md`
- Verify: all files changed by Tasks 1-5

**Interfaces:**
- Consumes: repository changes, generated metadata, Task 5 probe evidence, and unchecked runbook.
- Produces: an accurate three-state readiness ledger with final verification evidence.

- [ ] **Step 1: Add a failing readiness-ledger contract test**

Add to `provider/registry_docs_test.go`:

```go
func TestRegistryReadinessLedgerSeparatesEvidenceStates(t *testing.T) {
	ledger := readProjectFile(t, "../docs/provider-registry-readiness.md")
	for _, heading := range []string{
		"## Repository-complete", "## Publicly available", "## External pending",
	} {
		require.Contains(t, ledger, heading)
	}
	for _, marker := range []string{
		"v0.2.2", "docs/_index.md", "logoUrl", "contact@dimeski.net",
		".github/CODEOWNERS", "release-smoke", "community-packages/package-list.json",
		"publisher-names.json", "fact-sheet", "preview", "Registry CODEOWNER",
	} {
		require.Contains(t, ledger, marker)
	}
	require.Contains(t, ledger, "not yet Registry-ready")
	require.NotContains(t, ledger, "corrected release must be published")
}
```

- [ ] **Step 2: Run focused test to verify RED**

Run: `go test ./provider -run '^TestRegistryReadinessLedgerSeparatesEvidenceStates$' -count=1`

Expected: FAIL because the historical ledger does not use the three current evidence states and still says a corrected release is required.

- [ ] **Step 3: Rewrite the readiness ledger around current evidence**

Update `docs/provider-registry-readiness.md` with the current date and these sections:

```markdown
## Repository-complete

- Registry overview follows current six-language structure and passes contract/documentation checks.
- Public whale logo, attribution, and generated `logoUrl` metadata are present.
- Private security and conduct reports use `contact@dimeski.net`.
- `.github/CODEOWNERS` assigns repository ownership to `@dimeskigj`.
- The publication runbook records all external actions without claiming completion.
- Existing User-Agent, plugin download metadata, keywords, license, schema compatibility, checksums, SBOMs, attestations, and release-smoke workflow remain covered.

## Publicly available

- Stable GitHub release `v0.2.2` exists with provider archives and checksums.
- Preserve the exact dated Task 5 probe bullets here.
- State explicitly that anonymous availability is not runtime verification.

## External pending

- Successful `release-smoke` dispatch for `0.2.2` with all six jobs.
- Upstream `community-packages/package-list.json` entry.
- Maintainer-approved publisher display name and upstream `publisher-names.json` mapping to `dimeskigj`.
- Registry fact-sheet, `/check`, `/preview`, six-language and logo preview inspection.
- Pulumi Registry CODEOWNER approval, merge, deployment, and public page verification.
```

Keep an explicit statement that the provider is **not yet Registry-ready** until every external item has evidence. Remove stale statements that a corrected stable release has not been published or that no public release exists. Retain relevant residual risks: release-smoke has not run, Pulumi CLI downloads in that workflow are not checksum-verified, and broad engine acceptance remains deferred and is not a Registry blocker.

- [ ] **Step 4: Run focused Registry tests**

Run: `go test ./provider -run 'Registry|SchemaPublishingMetadata|GeneratedPublishingMetadata' -count=1`

Expected: PASS.

- [ ] **Step 5: Format and inspect**

Run: `gofmt -w provider/registry_docs_test.go provider/schema_test.go provider/provider.go`

Expected: files are formatted without errors.

Run: `git diff --check`

Expected: PASS.

- [ ] **Step 6: Run provider and race suites**

Run: `make test_provider`

Expected: PASS.

Run: `make test_race`

Expected: PASS. If the historical backup deadline race flake recurs, rerun the exact failing test once, record both outputs in the ledger, and do not broaden this work into an unrelated provider fix.

- [ ] **Step 7: Run generated, documentation, and SDK checks**

Run: `make check_codegen`

Expected: PASS with no generated drift.

Run: `make check_openapi`

Expected: PASS with no OpenAPI/client drift.

Run: `make docs_check`

Expected: PASS; record unrelated existing advisories or warnings without claiming they were fixed.

Run: `make build_sdks`

Expected: PASS for all five generated SDKs.

- [ ] **Step 8: Run lint, vulnerability, and license checks**

Run: `make lint`

Run: `make govulncheck`

Run: `make license`

Expected: all PASS. If `mise` or a pinned tool is unavailable, use the exact pinned equivalent from `.mise.toml`/`Makefile` where possible and record the substitution in the ledger.

- [ ] **Step 9: Verify no external action was performed or claimed**

Run: `git diff -- docs/registry-publication-runbook.md docs/provider-registry-readiness.md`

Expected: runbook tasks remain unchecked; the ledger lists release-smoke dispatch, upstream entries, fact-sheet, preview, approval, merge, and deployment under `External pending`.

Run: `git status --short`

Expected: only intended Task 6 files are uncommitted; unrelated pre-existing changes, if any, are not staged or modified.

- [ ] **Step 10: Commit the final ledger**

```bash
git add docs/provider-registry-readiness.md provider/registry_docs_test.go
git commit -m "docs: refresh registry readiness evidence"
```

- [ ] **Step 11: Confirm branch verification state**

Run: `git status --short`

Expected: clean worktree unless unrelated pre-existing changes were present and intentionally untouched.

Run: `git log --oneline -6`

Expected: separate commits exist for overview docs, logo metadata, contacts/ownership, runbook, and readiness evidence.
