package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestIpRangeRuleAlertBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		// Path to the Terraform configuration for the IP Range Rule Alert resource.
		TerraformDir: "../examples/resources/traceable_ip_range_rule_alert",
		Vars: map[string]interface{}{
			"name":              "tf_ip_range_rule_alert_test",
			"description":       "Testing IP Range Rule Alert resource",
			"event_severity":    "HIGH",
			"rule_action":       "RULE_ACTION_ALERT",
			"environment":       []string{"utkarsh_21", "utkarsh_crapi"},
			"raw_ip_range_data": []string{"192.168.1.0/24", "10.0.0.0/8"},
			"inject_request_headers": []map[string]interface{}{
				{
					"header_key":   "X-Test-Header",
					"header_value": "TestValue",
				},
			},
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":      os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for IP Range Rule Alert test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created IP Range Rule Alert resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_alert")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_ip_range_rule_alert_test", "Resource name should match test input")

	// Update test: change description and event_severity.
	terraformOptions.Vars["description"] = "Updated description"
	terraformOptions.Vars["event_severity"] = "MEDIUM"
	// terraformOptions.Vars["inject_request_headers"] = []map[string]interface{}{
	// 	{
	// 		"header_key":   "X-Updated-Header",
	// 		"header_value": "UpdatedValue",
	// 	},
	// }
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated IP Range Rule Alert resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_alert")
	require.Equal(t, "MEDIUM", output["event_severity"], "Event severity should be updated to MEDIUM")
	require.Equal(t, "Updated description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}

// package test

// import (
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/gruntwork-io/terratest/modules/terraform"
// 	"github.com/stretchr/testify/require"
// )

// func TestIpRangeRuleAlertBasic(t *testing.T) {
// 	terraformOptions := &terraform.Options{
// 		// Path to the Terraform configuration for the IP Range Rule Alert resource
// 		TerraformDir: "../examples/resources/traceable_ip_range_rule_alert",
// 		Vars: map[string]interface{}{
// 			"name":              "tf_ip_range_rule_alert_test",
// 			"description":       "Testing IP Range Rule Alert resource",
// 			"event_severity":    "HIGH",
// 			"rule_action":       "RULE_ACTION_ALERT",
// 			"environment":       []string{"env1", "env2"},
// 			"raw_ip_range_data": []string{"192.168.1.0/24", "10.0.0.0/8"},
// 			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
// 			"platform_url":      os.Getenv("PLATFORM_URL"),
// 		},
// 	}

// 	fmt.Println("Starting terraform init and apply for IP Range Rule Alert test")
// 	terraform.InitAndApply(t, terraformOptions)

// 	fmt.Println("Verifying created IP Range Rule Alert resource")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_alert")
// 	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
// 	require.Contains(t, output["name"], "tf_ip_range_rule_alert_test", "Resource name should match test input")

// 	// Update test: change description and event_severity
// 	terraformOptions.Vars["description"] = "Updated IP Range Rule Alert description"
// 	terraformOptions.Vars["event_severity"] = "MEDIUM"
// 	terraform.Apply(t, terraformOptions)

// 	fmt.Println("Verifying updated IP Range Rule Alert resource")
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_alert")
// 	require.Equal(t, "MEDIUM", output["event_severity"], "Event severity should be updated to MEDIUM")
// 	require.Equal(t, "Updated IP Range Rule Alert description", output["description"], "Description should be updated")
// 	require.NotEmpty(t, output["id"], "Resource ID should still be present")

// 	terraform.Destroy(t, terraformOptions)
// }


// package test

// import (
// 	// "fmt"
// 	"os"
// 	"testing"

// 	"github.com/gruntwork-io/terratest/modules/terraform"
// 	"github.com/traceableai/terraform-provider-traceable/test/logger"
// 	"github.com/stretchr/testify/require"
// )

// func TestIpRangeRuleAlertBasic(t *testing.T) {
// 	terraformOptions := &terraform.Options{
// 		// Path to the Terraform configuration for the IP Range Rule Alert resource
// 		TerraformDir: "../examples/resources/traceable_ip_range_rule_alert",
// 		Vars: map[string]interface{}{
// 			"name":              "tf_ip_range_rule_alert_test",
// 			"description":       "Testing IP Range Rule Alert resource",
// 			"event_severity":    "HIGH",
// 			"rule_action":       "RULE_ACTION_ALERT",
// 			"environment":       []string{"env1", "env2"},
// 			"raw_ip_range_data": []string{"192.168.1.0/24", "10.0.0.0/8"},
// 			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
// 			"platform_url":      os.Getenv("PLATFORM_URL"),
// 		},
// 	}

// 	logger.Log("Starting terraform init and apply for IP Range Rule Alert test")
// 	terraform.InitAndApply(t, terraformOptions)

// 	logger.Log("Verifying created IP Range Rule Alert resource")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_alert")
// 	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
// 	require.Contains(t, output["name"], "tf_ip_range_rule_alert_test", "Resource name should match test input")

// 	// Update test: change description and event_severity
// 	terraformOptions.Vars["description"] = "Updated IP Range Rule Alert description"
// 	terraformOptions.Vars["event_severity"] = "MEDIUM"
// 	terraform.Apply(t, terraformOptions)

// 	logger.Log("Verifying updated IP Range Rule Alert resource")
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_range_rule_alert")
// 	require.Equal(t, "MEDIUM", output["event_severity"], "Event severity should be updated to MEDIUM")
// 	require.Equal(t, "Updated IP Range Rule Alert description", output["description"], "Description should be updated")
// 	require.NotEmpty(t, output["id"], "Resource ID should still be present")

// 	terraform.Destroy(t, terraformOptions)
// }
