package resources_test

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccRateLimitiningResourceDefault(t *testing.T) {
	var rule_name="terraform_rate_limit_T1_" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRateLimitiningResourceConfigDefault(rule_name, "LOW","PT60S"), 
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.duration", "PT60S"),
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.event_severity", "LOW"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NU"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NF"),
				),
			},
			{
				ResourceName:      "traceable_rate_limiting.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccRateLimitiningResourceConfigDefault(rule_name, "HIGH","PT5M"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.duration", "PT5M"),
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.event_severity", "HIGH"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NU"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NF"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRateLimitiningResourceConfigDefault(name string, severity string,duration string) string {
	return fmt.Sprintf(acctest.RATE_LIMIT_CREATE, name, severity,duration)
}
