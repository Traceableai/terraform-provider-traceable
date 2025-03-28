package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	// "github.com/traceableai/terraform-provider-traceable/test/logger"
	"github.com/stretchr/testify/require"
)

func TestEmailDomainBlockBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		// Path to the Terraform configuration for the email domain block resource
		TerraformDir: "../examples/resources/traceable_email_domain_block",
		Vars: map[string]interface{}{
			"name":                    "tf_email_domain_block_test",
			"description":             "Testing email domain block policy",
			"rule_action":             "BLOCK",
			"event_severity":          "CRITICAL", // Options: LOW, MEDIUM, HIGH, CRITICAL
			"expiration":              "PT3600S",  // 1 hour block
			"environment":             []string{"utkarsh_crapi"},
			"data_leaked_email":       true,
			"disposable_email_domain": false,
			"email_domains":           []string{"example.com", "test.com"},
			"email_regexes":           []string{".*example.*", ".*test.*"},
			"email_fraud_score":       "HIGH",      
			"traceable_api_key":       os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":            os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for Email Domain Block test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created Email Domain Block resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_email_domain_block")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_email_domain_block_test", "Resource name should contain test string")

	// Update test: change description and expiration
	terraformOptions.Vars["description"] = "Updated email domain block description"
	terraformOptions.Vars["expiration"] = "PT7200S" // update to 2 hours
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated Email Domain Block resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_email_domain_block")
	require.Equal(t, "PT7200S", output["expiration"], "Expiration should be updated")
	require.Equal(t, "Updated email domain block description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}
