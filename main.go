package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/traceableai/terraform-provider-traceable/provider"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.Provider}
	plugin.Serve(opts)
}
