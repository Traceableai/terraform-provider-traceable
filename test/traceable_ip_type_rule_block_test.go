package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestIpTypeRuleBlockBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_ip_type_rule_block",
		Vars: map[string]interface{}{
			"name":              "tf_ip_type_rule_block_test",
			"description":       "Testing IP Type Rule Block resource",
			"rule_action":       "BLOCK",
			"event_severity":    "HIGH",
			"expiration":        "PT3600S",
			"environment":       []string{"utkarsh_21", "utkarsh_crapi"},
			"ip_types":          []string{"PUBLIC_PROXY", "TOR_EXIT_NODE"},
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":      os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for IP Type Rule Block test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created IP Type Rule Block resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_type_rule_block")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_ip_type_rule_block_test", "Resource name should match test input")

	// Update test: change expiration and description.
	terraformOptions.Vars["expiration"] = "PT7200S"
	terraformOptions.Vars["description"] = "Updated IP Type Rule Block description"
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated IP Type Rule Block resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_type_rule_block")
	require.Equal(t, "PT7200S", output["expiration"], "Expiration should be updated to 2 hours")
	require.Equal(t, "Updated IP Type Rule Block description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}

// package test

// import (
// 	// "fmt"
// 	"os"
// 	"testing"

// 	"github.com/gruntwork-io/terratest/modules/terraform"
// 	"github.com/traceableai/terraform-provider-traceable/test/logger"
// 	"github.com/stretchr/testify/require"
// )

// func TestIpTypeRuleBlockBasic(t *testing.T) {
// 	terraformOptions := &terraform.Options{
// 		// Path to the Terraform configuration for the IP Type Rule Block resource
// 		TerraformDir: "../examples/resources/traceable_ip_type_rule_block",
// 		Vars: map[string]interface{}{
// 			"name":           "tf_ip_type_rule_block_test",
// 			"description":    "Testing IP Type Rule Block resource",
// 			"rule_action":    "BLOCK", // default action is BLOCK
// 			"event_severity": "HIGH",  // Options: LOW, MEDIUM, HIGH, CRITICAL
// 			"expiration":     "PT3600S", // e.g., one hour block period
// 			"environment":    []string{"env1", "env2"},
// 			"ip_types":       []string{"PUBLIC PROXY", "TOR EXIT NODE"}, // Must be one or more valid IP types
// 			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
// 			"platform_url":      os.Getenv("PLATFORM_URL"),
// 		},
// 	}

// 	logger.Log("Starting terraform init and apply for IP Type Rule Block test")
// 	terraform.InitAndApply(t, terraformOptions)

// 	logger.Log("Verifying created IP Type Rule Block resource")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_type_rule_block")
// 	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
// 	require.Contains(t, output["name"], "tf_ip_type_rule_block_test", "Resource name should match the test input")

// 	// Update test: change expiration and description
// 	terraformOptions.Vars["expiration"] = "PT7200S" // update to 2 hours
// 	terraformOptions.Vars["description"] = "Updated IP Type Rule Block description"
// 	terraform.Apply(t, terraformOptions)

// 	logger.Log("Verifying updated IP Type Rule Block resource")
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_type_rule_block")
// 	require.Equal(t, "PT7200S", output["expiration"], "Expiration should be updated to 2 hours")
// 	require.Equal(t, "Updated IP Type Rule Block description", output["description"], "Description should be updated")
// 	require.NotEmpty(t, output["id"], "Resource ID should still be present")

// 	terraform.Destroy(t, terraformOptions)
// }
