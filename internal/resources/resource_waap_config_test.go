package resources_test

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
)

func TestAccWaapConfigResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccWaapConfigResourceConfigDefault("hitech-3"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_waap_config.test", "environment", "hitech-3"),
				), //checking with state
			},
			{
				Config: testAccWaapConfigResourceConfigDefault("fintech-2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_waap_config.test", "environment", "fintech-2"),	
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWaapConfigResourceConfigDefault(env string) string {
	return fmt.Sprintf(acctest.WAAP_CONFIG, env)
}
