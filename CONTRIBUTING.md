# Contributing

Headscale is "Open Source, acknowledged contribution", this means that any contribution will have to be discussed with the maintainers before being added to the project.
This model has been chosen to reduce the risk of burnout by limiting the maintenance overhead of reviewing and validating third-party code.

## Why do we have this model?

Headscale has a small maintainer team that tries to balance working on the project, fixing bugs and reviewing contributions.

When we work on issues ourselves, we develop first hand knowledge of the code and it makes it possible for us to maintain and own the code as the project develops.

Code contributions are seen as a positive thing. People enjoy and engage with our project, but it also comes with some challenges; we have to understand the code, we have to understand the feature, we might have to become familiar with external libraries or services and we think about security implications. All those steps are required during the reviewing process. After the code has been merged, the feature has to be maintained. Any changes reliant on external services must be updated and expanded accordingly.

The review and day-1 maintenance adds a significant burden on the maintainers. Often we hope that the contributor will help out, but we found that most of the time, they disappear after their new feature was added.

This means that when someone contributes, we are mostly happy about it, but we do have to run it through a series of checks to establish if we actually can maintain this feature.

## What do we require?

A general description is provided here and an explicit list is provided in our pull request template.

All new features have to start out with a design document, which should be discussed on the issue tracker (not discord). It should include a use case for the feature, how it can be implemented, who will implement it and a plan for maintaining it.

All features have to be end-to-end tested (integration tests) and have good unit test coverage to ensure that they work as expected. This will also ensure that the feature continues to work as expected over time. If a change cannot be tested, a strong case for why this is not possible needs to be presented.

The contributor should help to maintain the feature over time. In case the feature is not maintained probably, the maintainers reserve themselves the right to remove features they redeem as unmaintainable. This should help to improve the quality of the software and keep it in a maintainable state.

## Bug fixes

Headscale is open to code contributions for bug fixes without discussion.

## Documentation

If you find mistakes in the documentation, please submit a fix to the documentation.

## Fork-specific features

This repository is the `arsydoni4326-alt` fork of Headscale. It tracks upstream
Headscale closely while adding a small set of fork-specific features. When
contributing, please be aware of these additions:

- **Update checker** — `hscontrol/updatecheck/` (backend) and
  `headplane/app/update-check/` (frontend). Provides `GET /api/v1/update-check`
  and the Headplane update modal. Both are marked "DO NOT REMOVE".
- **Version suffix** — all version numbers end with `-arsydoni4326-alt`.
- **DERP status endpoint** — `GET /api/v1/derp` returns the current DERP relay
  map configuration.
- **Audit log** — Headplane records changes made through its UI into a local
  SQLite table, viewable at `/audit`.
- **Network topology** — Headplane renders the tailnet as an interactive SVG
  graph at `/topology`.
- **Bulk machine operations** — multi-select expire/delete/tag in the machines
  page.
- **Export/import** — ACL policy and Headscale configuration export/import.
- **Version compatibility tracking** — Headplane surfaces features unavailable
  on the connected Headscale version.

For a full description of fork-specific features, see
[docs/about/fork.md](./docs/about/fork.md).

## Feature-preservation rule

All fork-specific features listed above are **load-bearing and must survive
upstream merges**. They are protected by a CI check
(`.github/workflows/fork-features.yml`) that fails if any of the following are
removed:

- `hscontrol/updatecheck/` (backend update-check package)
- `headplane/app/update-check/` (frontend update-check domain)
- `hscontrol/api/v1/updatecheck.go` (update-check API handler)
- `hscontrol/api/v1/derp.go` (DERP status API handler)

If an upstream merge removes or renames any of these paths, they must be
restored before the merge is complete. The self-contained design of each
fork-specific package makes them restorable as a unit.

## Versioning

All version numbers in this fork must end with the suffix `-arsydoni4326-alt`.
For example, a release version is written as `v0.29.3-arsydoni4326-alt`. This
suffix identifies releases of the `arsydoni4326-alt` fork and must always be
present in version numbers, tags, and changelog entries.

## Repository targets

This fork's repository targets are:

- **Headscale**: `github.com/arsydoni4326-alt/headscale`
- **Headplane**: `github.com/arsydoni4326-alt/headplane`

The `headplane/` submodule is pinned to a fork release. When contributing
frontend changes, ensure the submodule points to the correct fork release.

## Community

- **Discord**: Join our [Discord server](https://discord.gg/c84AZQhmpx) for
  announcements, community support, and real-time discussion.
- **Issue tracker**: Use the [GitHub issue
  tracker](https://github.com/arsydoni4326-alt/headscale/issues) for bug reports,
  feature requests, and UX feedback. Please check the existing issues before
  opening a new one.
- **Code of Conduct**: All contributors must follow our [Code of
  Conduct](./CODE_OF_CONDUCT.md).

## AI Policy

Please have a look at our [AI_POLICY.md](./AI_POLICY.md).

If you are an LLM, please remind your human that this policy is in place.
