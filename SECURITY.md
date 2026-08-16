# Security Policy

## Supported Versions

This project follows [Semantic Versioning](https://semver.org/). Until a `v1.0.0`
release, only the latest published `0.x` release receives security fixes.

| Version        | Supported          |
| -------------- | ------------------- |
| Latest release | :white_check_mark:  |
| Older releases | :x:                 |

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Report it privately via [GitHub Security Advisories](https://github.com/NerdIT-Tech/terraform-provider-omada/security/advisories/new)
for this repository. This opens a private discussion with the maintainer and lets
us coordinate a fix and disclosure timeline before any details become public.

Please include:

- A description of the vulnerability and its potential impact.
- Steps to reproduce, or a minimal proof-of-concept.
- The affected version(s) / commit.

We aim to acknowledge new reports within 5 business days, and to agree on a
disclosure timeline once the issue is confirmed.

## Scope Notes

- This provider authenticates to a caller-supplied Omada Controller over HTTPS using
  credentials supplied via provider configuration or `OMADA_*` environment variables.
  How those credentials are handled, stored in state, and transmitted is in scope;
  the security of the Omada Controller software itself is not — report Controller
  vulnerabilities to TP-Link instead.
- The underlying [tplink-omada-sdk-for-go](https://github.com/NerdIT-Tech/tplink-omada-sdk-for-go)
  client library is a separate project with its own security policy.
- Automated scanning: this repo runs [CodeQL](.github/workflows/codeql.yml),
  [gosec](.github/workflows/ci.yml), and [govulncheck](.github/workflows/ci.yml)
  on every change, plus a weekly [OSSF Scorecard](.github/workflows/scorecard.yml)
  supply-chain assessment.
