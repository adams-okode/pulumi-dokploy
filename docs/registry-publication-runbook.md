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
  examples, configuration, generated resource pages, links, and logo must
  render correctly.
- [ ] Obtain approval from a Pulumi Registry CODEOWNER. The local
  `.github/CODEOWNERS` file does not satisfy this upstream approval.

## 4. Verify Publication

- [ ] Merge only through the upstream maintainer process.
- [ ] After deployment, open the public Registry page and verify discovery,
  publisher display name, logo, overview, generated API documentation, and
  installation links.
- [ ] Update `docs/provider-registry-readiness.md` with the smoke run, pull
  request, approval, merge, and public page evidence.
