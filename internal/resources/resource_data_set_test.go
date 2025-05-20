package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccDataSetDefault(t *testing.T) {
	var rule_name="terraform_date_set_" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSetConfigDefault(rule_name, "terraform data set"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("internal/resources/resource_rate_limiting_test.go.test", "name", rule_name),
				), //checking with state
			},
			{
				ResourceName:      "traceable_data_set.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccDataSetConfigDefault(rule_name, "terraform data set updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_data_set.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_data_set.test", "description", "terraform data set updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDataSetConfigDefault(name string, desc string) string {
	return fmt.Sprintf(acctest.DATA_SET, name, desc)
}
