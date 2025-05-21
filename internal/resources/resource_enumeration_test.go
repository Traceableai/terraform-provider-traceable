package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccEnumerationResourceDefault(t *testing.T) {
	var rule_name="terraform_enumeration_T1" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnumerationResourceConfigDefault(rule_name, "PT60S","LOW"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_enumeration.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_enumeration.test", "action.duration", "PT60S"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AF"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AW"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_enumeration.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccEnumerationResourceConfigDefault(rule_name, "PT5M","HIGH"),  
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_enumeration.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_enumeration.test", "action.duration", "PT5M"),
					resource.TestCheckResourceAttr("traceable_enumeration.test", "action.event_severity", "HIGH"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AF"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AW"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEnumerationResourceConfigDefault(name string, duration string,severity string) string {
	return fmt.Sprintf(acctest.ENUMERATION_RESOURCE, name, duration,severity)
}
