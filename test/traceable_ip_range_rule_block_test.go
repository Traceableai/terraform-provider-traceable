package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	// "github.com/traceableai/terraform-provider-traceable/test/logger"
	"github.com/stretchr/testify/require"
)

func TestIpRangeRuleBlockBasic(t *testing.T) {
	// Define Terraform options with example variables
	terraformOptions := &terraform.Options{
		// Path to your Terraform configuration for IP range rule block
		TerraformDir: "../examples/resources/traceable_ip_range_rule_block",
		Vars: map[string]interface{}{
			"name":               "tf_ip_range_rule_block_test",
			"description":        "Testing IP range rule block",
			"rule_action":        "RULE_ACTION_BLOCK", // Options: RULE_ACTION_BLOCK_ALL_EXCEPT, RULE_ACTION_BLOCK, etc.
			"event_severity":     "HIGH",              // Options: LOW, MEDIUM, HIGH, CRITICAL
			"expiration":         "PT3600S",           // e.g., one hour block
			"environments":        []string{"utkarsh_crapi","utkarsh_tr"},    // Note: schema expects a set, pass a list of strings
			"raw_ip_range_data":  []string{"192.168.0.0/16", "10.0.0.0/8"},
			"traceable_api_key":  os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":       os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for IP Range Rule Block test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created IP Range Rule Block resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_block")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_ip_range_rule_block_test", "Name should contain test string")

	// Update test: change expiration and description
	terraformOptions.Vars["expiration"] = "PT7200S" // update expiration to 2 hours
	terraformOptions.Vars["description"] = "Updated IP range rule block description"
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated IP Range Rule Block resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_block")
	require.Equal(t, "PT7200S", output["expiration"], "Expiration should be updated")
	require.Equal(t, "Updated IP range rule block description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still not be empty")

	terraform.Destroy(t, terraformOptions)
}
