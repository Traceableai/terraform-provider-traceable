package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccDLPUserBasedResourceDefault(t *testing.T) {
	var rule_name="terraform_dlp_user_based_" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDLPUserBasedResourceConfigDefault(rule_name, "PT60S","LOW"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_data_loss_prevention_user_based.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_data_loss_prevention_user_based.test", "action.duration", "PT60S"),
					resource.TestCheckResourceAttr("traceable_data_loss_prevention_user_based.test", "action.event_severity", "LOW"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_data_loss_prevention_user_based.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccDLPUserBasedResourceConfigDefault(rule_name, "PT5M","HIGH"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_data_loss_prevention_user_based.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_data_loss_prevention_user_based.test", "action.duration", "PT5M"),
					resource.TestCheckResourceAttr("traceable_data_loss_prevention_user_based.test", "action.event_severity", "HIGH"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDLPUserBasedResourceConfigDefault(name string, duration string, severity string) string {
	return fmt.Sprintf(acctest.DLP_USER_BASED, name, duration, severity)
}
