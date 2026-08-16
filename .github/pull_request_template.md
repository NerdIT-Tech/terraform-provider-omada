<!-- PR title must be a Conventional Commit line, e.g. `fix(provider): validate host URL scheme` — CI lints it. -->

## What & why

<!-- Short description of the change and the problem it solves. Link related issues. -->

## Checklist

- [ ] `just check` (fmt, lint, docs, test) passes locally
- [ ] New or changed exported surface has unit tests
- [ ] If the provider schema changed, `docs/` was regenerated via `just docs` and committed
- [ ] `CHANGELOG.md` is untouched (release-please manages it)

<!-- If this PR intentionally skips tests or docs, say why: -->
