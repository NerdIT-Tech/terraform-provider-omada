# terraform-provider-omada

Terraform provider for the TP-Link Omada Controller API, built on
[tplink-omada-sdk-for-go](https://github.com/NerdIT-Tech/tplink-omada-sdk-for-go).

> **Status:** early days. The provider currently implements the
> `omada_site` resource; more resources and data sources are added as the
> underlying Go SDK grows support for them.

## Requirements

- [Go](https://go.dev/doc/install) (see `go.mod` for the minimum version)
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [just](https://github.com/casey/just) for the development commands below

## Developing

```console
just build   # build the provider binary into ./bin
just test    # run unit tests
just testacc # run acceptance tests, including the Gherkin scenarios in internal/provider/features/, against a real controller (requires OMADA_* env vars)
just lint    # go vet + golangci-lint
just docs    # regenerate docs/ from the schema and templates/
just check   # fmt + lint + docs + test, mirrors CI
```

To exercise the provider locally against a real Terraform configuration
without publishing it, add a
[dev overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)
block to `~/.terraformrc` pointing at `just install`'s `$GOBIN` output.

## Configuration

```hcl
provider "omada" {
  host      = "https://omada.example.com:8043"
  omadac_id = "example-omadac-id"
}
```

`host`, `client_id`, `client_secret`, and `omadac_id` may also be set via the
`OMADA_HOST`, `OMADA_CLIENT_ID`, `OMADA_CLIENT_SECRET`, and `OMADA_OMADAC_ID`
environment variables — recommended so credentials stay out of configuration
files. See [`docs/index.md`](docs/index.md) for the full schema.

## Resources

- `omada_site` — manages a site, the top-level container for a location's
  devices and clients. See
  [`docs/resources/site.md`](docs/resources/site.md).

## Security

See [`SECURITY.md`](SECURITY.md) for the supported-versions policy and how to
report a vulnerability privately.

## License

[MIT](LICENSE)
