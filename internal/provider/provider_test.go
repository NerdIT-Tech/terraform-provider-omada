package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate the provider during
// acceptance testing. The factory function is called for each Terraform CLI
// command executed to create a provider server that the CLI can connect to
// and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"omada": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccPreCheck validates the necessary environment variables are set
// before running acceptance tests that talk to a real Omada Controller.
func testAccPreCheck(t *testing.T) {
	for _, env := range []string{envHost, envClientID, envClientSecret, envOmadacID} {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}
