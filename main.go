package main

import (

	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/traceableai/terraform-provider-traceable/provider"
)
var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {

	var debug bool
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/traceableai/traceable",
		Debug:   debug,
	}
	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
