package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/traceableai/terraform-provider-traceable/internal/acctest"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func TestAccEmailDomainDefault(t *testing.T) {
	var rule_name="terraform_email_domain_" + utils.GenerateRandomString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccEmailDomainConfigDefault(rule_name, "LOW","PT60S"), // initalize the resource
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_email_domain.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_malicious_email_domain.test", "duration", "PT60S"),
					resource.TestCheckResourceAttr("traceable_malicious_email_domain.test", "event_severity", "LOW"),
				), //checking with state
			},
			{
				ResourceName:      "traceable_malicious_email_domain.test",
				ImportState:       true,
				ImportStateId:     rule_name,
			},
			{
				Config: testAccEmailDomainConfigDefault(rule_name, "HIGH","PT5M"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("traceable_malicious_email_domain.test", "name", rule_name),
					resource.TestCheckResourceAttr("traceable_malicious_email_domain.test", "duration", "PT5M"),
					resource.TestCheckResourceAttr("traceable_malicious_email_domain.test", "event_severity", "HIGH"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccEmailDomainConfigDefault(name string, severity string, duration string) string {
	return fmt.Sprintf(acctest.EMAIL_DOMAIN, name, severity, duration)
}
