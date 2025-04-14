package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
)

func TestAccVariantResourceDefault(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccMaliciousIpRangeResourceConfigDefault("tf_automation_ip_range", "ALERT"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_ip_range.test", "name", "tf_automation_ip_range"),
					resource.TestCheckResourceAttr("traceable_malicious_ip_range.test", "action", "ALERT"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_ip_range.test", "ip_range.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_ip_range.test", "ip_range.*", "192.168.1.2"),
				), //checking with state
			},

			{
				ResourceName:      "traceable_malicious_ip_range.test",
				ImportState:       true,
				ImportStateId:     "tf_automation_ip_range",
				ImportStateVerify: true,
			},

			{
				Config: testAccMaliciousIpRangeResourceConfigDefault("tf_automation_ip_range", "ALLOW"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_ip_range.test", "name", "tf_automation_ip_range"),
					resource.TestCheckResourceAttr("traceable_malicious_ip_range.test", "action", "ALLOW"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_ip_range.test", "ip_range.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_ip_range.test", "ip_range.*", "192.168.1.2"),
				),
			},

			{
				Config: testAccMaliciousIpRangeResourceConfigDefault("tf_automation_ip_range", "BLOCK"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_ip_range.test", "name", "tf_automation_ip_range"),
					resource.TestCheckResourceAttr("traceable_malicious_ip_range.test", "action", "BLOCK"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_ip_range.test", "ip_range.*", "192.168.1.1"),
					resource.TestCheckTypeSetElemAttr("traceable_malicious_ip_range.test", "ip_range.*", "192.168.1.2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "traceable_malicious_ip_range.test",
				ImportState:       true,
				ImportStateId:     "tf_automation_ip_range",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMaliciousIpRangeResourceConfigDefault(name string, action string) string {
	return fmt.Sprintf(`
resource "traceable_malicious_ip_range" "test" {
  name = "%s"
  description = "revamp"
  enabled = true
  event_severity = "LOW"
  duration = "PT1M"
  action = "%s"
  ip_range = ["192.168.1.1","192.168.1.2"]
  environments = ["env1","env2"]
}
`, name, action)
}
