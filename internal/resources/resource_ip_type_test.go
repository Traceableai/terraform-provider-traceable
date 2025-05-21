package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccIpTypeDefault(t *testing.T) {
	var rule_name="terraform_ip_type_" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccIpTypeConfigDefault(rule_name, "LOW","PT60S"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_ip_type.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_malicious_ip_type.test", "duration", "PT60S"),
					resource.TestCheckResourceAttr("traceable_malicious_ip_type.test", "event_severity", "LOW"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_malicious_ip_type.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccIpTypeConfigDefault(rule_name, "HIGH","PT5M"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_ip_type.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_malicious_ip_type.test", "duration", "PT5M"),
					resource.TestCheckResourceAttr("traceable_malicious_ip_type.test", "event_severity", "HIGH"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIpTypeConfigDefault(name string, severity string, duration string) string {
	return fmt.Sprintf(acctest.IP_TYPE_RESOURCE, name, severity, duration)
}
