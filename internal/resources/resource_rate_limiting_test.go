package resources_test

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccRateLimitiningResourceDefault(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRateLimitiningResourceConfigDefault("rate_limit_T1", "ALERT"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "name", "rate_limit_T1"),
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.action_type", "ALERT"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NU"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NF"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_rate_limiting.test",
				ImportState:       true,
				ImportStateId:     "rate_limit_T1",
				ImportStateVerify: true,
			},
			{
				Config: testAccRateLimitiningResourceConfigDefault("rate_limit_T1", "ALERT"), //change this to block when shreyansh fix it
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "name", "rate_limit_T1"),
					resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.action_type", "ALERT"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NU"),
					resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NF"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRateLimitiningResourceConfigDefault(name string, action string) string {
	name = name + utils.GenerateRandomString(8)
	return fmt.Sprintf(acctest.RATE_LIMIT_CREATE, name, action)
}
