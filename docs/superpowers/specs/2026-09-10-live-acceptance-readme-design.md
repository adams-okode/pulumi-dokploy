# Live Acceptance README Design

## Goal

Add one operator guide for the live acceptance test strategy. Put the guide at
`tests/README.md`. Use short, direct instructions and consistent terms.

## Content

The README explains the two test layers:

- Direct-provider lifecycle tiers in `provider/`.
- The Pulumi Automation API smoke test in `tests/`.

The README documents:

- The purpose and limits of each layer.
- Required tools and environment variables.
- The local provider build and plugin discovery method.
- The safe tier order and exact test commands.
- Optional Registry, GitLab, MongoDB replica, custom certificate, and server
  prerequisites.
- Serial execution and server load limits.
- Cleanup, absence checks, and the stop-marker contract.
- Secret handling and prohibited output.
- Pass, skip, provider defect, server defect, and environment limitation
  classifications.
- The location and required use of sanitized run summaries and bug reports.

The README must not contain real credentials, endpoints, IDs, or `.env` data.
It must not tell test code to load `.env`.

## Structure

Use this section order:

1. Purpose.
2. Test layers.
3. Safety rules.
4. Prerequisites.
5. Local setup.
6. Test commands.
7. Optional coverage.
8. Cleanup and stop behavior.
9. Result classification.
10. Reports.

Add a short link in `CONTRIBUTING.md` to `tests/README.md`.

## Writing Rules

- Use short sentences.
- Use active voice.
- Give one instruction in each numbered step.
- Use the same term for the same item throughout the document.
- Define abbreviations before use.
- Avoid idioms, informal phrases, and unnecessary words.
- Use explicit warnings for unsafe actions.
- Do not mention the writing standard used for the document.

## Verification

Add a deterministic content test in `tests/acceptance_test.go`. The test checks
that the README contains the tier names, required variables, serial command
flags, stop-marker variable, safety warnings, report paths, and a link from
`CONTRIBUTING.md`. The test also rejects instructions that make Go test code
read `.env`.

## Acceptance Criteria

- `tests/README.md` is the primary live acceptance operator guide.
- The guide describes both live test layers and all four provider tiers.
- The guide gives safe serial commands with `-parallel=1`.
- The guide documents all required and optional prerequisites.
- The guide documents cleanup and stop-marker behavior.
- The guide documents secret rules and result classification.
- `CONTRIBUTING.md` links to the guide.
- A deterministic test protects the required content.
