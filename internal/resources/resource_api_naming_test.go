package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccApiNamingResourceDefault(t *testing.T) {
	var rule_name="terraform_api_naming" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccApiNamingResourceConfigDefault(rule_name, false), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_api_naming.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_api_naming.test", "disabled", "false"),
				), //checking with state
			},
			{
				Config: testAccApiNamingResourceConfigDefault(rule_name, true),	
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_api_naming.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_api_naming.test", "disabled", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApiNamingResourceConfigDefault(name string, disabled bool) string {
	return fmt.Sprintf(acctest.API_NAMING, name, disabled)
}
