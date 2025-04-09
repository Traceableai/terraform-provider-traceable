package main

import "fmt"

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/traceableai/terraform-provider-traceable/internal/provider"
)

var (
	//comes from goreleaser
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
	fmt.Println(version)

}
