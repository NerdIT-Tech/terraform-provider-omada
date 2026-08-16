// Package main provides the entrypoint for the Terraform provider plugin
// binary. Terraform launches this binary and communicates with it over the
// plugin protocol; the provider implementation itself lives in
// internal/provider.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/NerdIT-Tech/terraform-provider-omada/internal/provider"
)

// version is set via -ldflags at build time, e.g.:
//
//	go build -ldflags "-X main.version=$(VERSION)"
//
// goreleaser sets this automatically from the release tag.
var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// Address must match the registry source used in provider
		// requirements, e.g.:
		//   terraform {
		//     required_providers {
		//       omada = {
		//         source = "NerdIT-Tech/omada"
		//       }
		//     }
		//   }
		Address: "registry.terraform.io/NerdIT-Tech/omada",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
