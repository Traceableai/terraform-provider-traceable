package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccCustomSignatureResourceDefault(t *testing.T) {
	var rule_name="terraform_custom_signature_" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccCustomSignatureResourceConfigDefault(rule_name, "PT60S","LOW"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_custom_signature.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_custom_signature.test", "action.duration", "PT60S"),
					resource.TestCheckResourceAttr("traceable_custom_signature.test", "action.event_severity", "LOW"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_custom_signature.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccCustomSignatureResourceConfigDefault(rule_name, "PT5M","HIGH"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_custom_signature.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_custom_signature.test", "action.duration", "PT5M"),
					resource.TestCheckResourceAttr("traceable_custom_signature.test", "action.event_severity", "HIGH"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomSignatureResourceConfigDefault(name string, duration string, severity string) string {
	return fmt.Sprintf(acctest.CUSTOM_SIGNATURE_RESOURCE, name, duration, severity)
}
