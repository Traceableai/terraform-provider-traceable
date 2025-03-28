package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestRegionRuleBlockBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		// Path to your Terraform configuration for the region rule block resource
		TerraformDir: "../examples/resources/traceable_region_rule_block",
		Vars: map[string]interface{}{
			"name":         "tf_region_rule_block_test_1",
			"description":  "Testing region rule block resource",
			"rule_action":  "BLOCK",  // Allowed values: BLOCK, BLOCK_ALL_EXCEPT
			"event_severity": "HIGH", // LOW, MEDIUM, HIGH, CRITICAL
			"expiration":   "PT3600S", // 1 hour block period
			"environment":  []string{"utkarsh_crapi"}, // variable type is set(string) in tf; pass as list
			"regions":      []string{"afghanistan", "albania"},
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":      os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for Region Rule Block test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created Region Rule Block resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_region_rule_block")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_region_rule_block_test", "Name should match test input")

	// Update test: change expiration and description
	terraformOptions.Vars["expiration"] = "PT7200S" // update to 2 hours
	terraformOptions.Vars["description"] = "Updated region rule block description"
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated Region Rule Block resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_region_rule_block")
	require.Equal(t, "PT7200S", output["expiration"], "Expiration should be updated")
	require.Equal(t, "Updated region rule block description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}
