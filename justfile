default_provider_name := "omada"

# Build the provider binary.
build:
    go build -o bin/terraform-provider-{{default_provider_name}} .

# Install the provider binary into the local Terraform plugin cache for dev-overrides use.
install:
    go install .

# Run unit tests.
test:
    go test -count=1 ./...

# Run acceptance tests against a real Omada Controller (requires OMADA_* env vars).
testacc:
    TF_ACC=1 go test -count=1 -v -timeout 120m ./...

# Format Go sources.
fmt:
    gofmt -s -w .

# Run static analysis and linting.
lint:
    go vet ./...
    golangci-lint run ./...

# Regenerate provider documentation from schema and templates/.
docs:
    go tool tfplugindocs generate --provider-name {{default_provider_name}}

# Tidy and verify go.mod/go.sum.
tidy:
    go mod tidy
    go mod verify

# Run the full pre-commit check: format, lint, docs, test.
check: fmt lint docs test
    git diff --exit-code -- docs || (echo "docs/ is out of date, run 'just docs' and commit the result" && exit 1)
