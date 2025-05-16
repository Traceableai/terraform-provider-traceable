package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
)

func TestAccEnumerationResourceDefault(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnumerationResourceConfigDefault("enumeration_T1", "ALERT"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_enumeration.test", "name", "enumeration_T1"),
					resource.TestCheckResourceAttr("traceable_enumeration.test", "action.action_type", "ALERT"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AF"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AW"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_enumeration.test",
				ImportState:       true,
				ImportStateId:     "enumeration_T1",
				ImportStateVerify: true,
			},
			{
				Config: testAccEnumerationResourceConfigDefault("enumeration_T1", "ALERT"),  // change this to BLOCK when shreyeansh fixes it
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_enumeration.test", "name", "enumeration_T1"),
					resource.TestCheckResourceAttr("traceable_enumeration.test", "action.action_type", "ALERT"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AF"),
					resource.TestCheckTypeSetElemAttr("traceable_enumeration.test", "sources.regions.regions_ids.*", "AW"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEnumerationResourceConfigDefault(name string, action string) string {

	return fmt.Sprintf(acctest.ENUMERATION_RESOURCE, name, action)
}
