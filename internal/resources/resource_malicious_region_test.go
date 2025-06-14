package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
)

func TestAccMaliciousRegionResourceDefault(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccMaliciousRegionResourceConfigDefault("tf_automation_region", "ALERT"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_region.test", "name", "tf_automation_region"),
					resource.TestCheckResourceAttr("traceable_malicious_region.test", "action", "ALERT"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_region.test", "regions.*", "NU"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_region.test", "regions.*", "NF"),
				), //checking with state
			},

			{
				ResourceName:      "traceable_malicious_region.test",
				ImportState:       true,
				ImportStateId:     "tf_automation_region",
				ImportStateVerify: true,
			},

			{
				Config: testAccMaliciousRegionResourceConfigDefault("tf_automation_region", "BLOCK"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_region.test", "name", "tf_automation_region"),
					resource.TestCheckResourceAttr("traceable_malicious_region.test", "action", "BLOCK"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_region.test", "regions.*", "NU"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_region.test", "regions.*", "NF"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "traceable_malicious_region.test",
				ImportState:       true,
				ImportStateId:     "tf_automation_region",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMaliciousRegionResourceConfigDefault(name string, action string) string {
	return fmt.Sprintf(`
resource "traceable_malicious_region" "test" {
  name = "%s"
  description = "revamp"
  enabled = true
  event_severity = "LOW"
  duration = "PT1M"
  action = "%s"
  regions = ["NU","NF"]
  environments = ["env1","env2"]
}
`, name, action)
}
