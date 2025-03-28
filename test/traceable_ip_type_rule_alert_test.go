package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	// "github.com/traceableai/terraform-provider-traceable/test/logger"
	"github.com/stretchr/testify/require"
)

func TestIpTypeRuleAlertBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		// Path to your Terraform configuration for the IP Type Rule Alert resource
		TerraformDir: "../examples/resources/traceable_ip_type_rule_alert",
		Vars: map[string]interface{}{
			"name":           "tf_ip_type_rule_alert_test",
			"description":    "Automation test for IP Type Rule Alert resource",
			"rule_action":    "ALERT", // default is ALERT
			"event_severity": "HIGH",  // Options: LOW, MEDIUM, HIGH, CRITICAL
			"environment":    []string{"utkarsh_crapi", "utkarsh_tr"},
			"ip_types":       []string{"PUBLIC_PROXY", "BOT"}, // Must be one or more from the allowed set: ANONYMOUS VPN, HOSTING PROVIDER, PUBLIC PROXY, TOR EXIT NODE, BOT
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

	fmt.Println("Starting terraform init and apply for IP Type Rule Alert test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created IP Type Rule Alert resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_ip_type_rule_alert")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_ip_type_rule_alert_test", "Resource name should contain test string")

	// Update test: Change event_severity and update headers.
	terraformOptions.Vars["event_severity"] = "MEDIUM"
	terraformOptions.Vars["inject_request_headers"] = []map[string]interface{}{
		{
			"header_key":   "X-Updated-Header",
			"header_value": "UpdatedValue",
		},
	}
	terraform.Apply(t, terraformOptions)

	
	fmt.Println("Verifying updated IP Type Rule Alert resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_ip_type_rule_alert")
	require.Equal(t, "MEDIUM", output["event_severity"], "Event severity should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}
