package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	// "github.com/traceableai/terraform-provider-traceable/test/logger"
	"github.com/stretchr/testify/require"
)

func TestIpRangeRuleAllowBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		// Path to your Terraform configuration for the IP Range Rule Allow resource.
		TerraformDir: "../examples/resources/traceable_ip_range_rule_allow",
		Vars: map[string]interface{}{
			"name":              "tf_ip_range_rule_allow_test",
			"description":       "Testing IP Range Rule Allow resource",
			"rule_action":       "RULE_ACTION_ALLOW", // default action
			"expiration":        "PT3600S",            // e.g. one hour allow period
			"environment":       []string{"utkarsh_crapi,utkarsh_tr"},
			"raw_ip_range_data": []string{"192.168.1.0/24", "10.0.0.0/8"},
			// "inject_request_headers": []map[string]interface{}{
			// 	{
			// 		"header_key":   "X-Test-Header",
			// 		"header_value": "TestValue",
			// 	},
			// },
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":      os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for IP Range Rule Allow test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created IP Range Rule Allow resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_allow")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_ip_range_rule_allow_test", "Resource name should match test input")

	// Update test: change expiration and description.
	terraformOptions.Vars["expiration"] = "PT7200S" // update to 2 hours
	terraformOptions.Vars["description"] = "Updated IP Range Rule Allow description"
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated IP Range Rule Allow resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_allow")
	require.Equal(t, "PT7200S", output["expiration"], "Expiration should be updated to 2 hours")
	require.Equal(t, "Updated IP Range Rule Allow description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}
