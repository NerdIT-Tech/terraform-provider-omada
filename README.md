# terraform-provider-omada

Terraform provider for the TP-Link Omada Controller API, built on
[tplink-omada-sdk-for-go](https://github.com/NerdIT-Tech/tplink-omada-sdk-for-go).

> **Status:** early scaffold. The provider builds and serves its
> configuration schema, but no resources or data sources are implemented
> yet — they're blocked on the underlying Go SDK, which doesn't have a
> client implementation to wire up yet either.

## Requirements

- [Go](https://go.dev/doc/install) (see `go.mod` for the minimum version)
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [just](https://github.com/casey/just) for the development commands below

## Developing

```console
just build   # build the provider binary into ./bin
just test    # run unit tests
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

## Security

See [`SECURITY.md`](SECURITY.md) for the supported-versions policy and how to
report a vulnerability privately.

## License

[MIT](LICENSE)
