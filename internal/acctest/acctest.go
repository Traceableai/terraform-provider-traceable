package acctest

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/traceableai/terraform-provider-traceable/internal/provider"
	"os"
	"testing"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"traceable": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("API_TOKEN"); v == "" {
		t.Fatal("API_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("PLATFORM_URL"); v == "" {
		t.Fatal("PLATFORM_URL must be set for acceptance tests")
	}

}
