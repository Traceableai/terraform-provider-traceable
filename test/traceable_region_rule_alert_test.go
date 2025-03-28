package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestRegionRuleAlertBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		// Path to the Terraform configuration for the Region Rule Alert resource
		TerraformDir: "../examples/resources/traceable_region_rule_alert",
		Vars: map[string]interface{}{
			"name":           "tf_region_rule_alert_test",
			"description":    "Testing Region Rule Alert resource",
			"rule_action":    "ALERT",
			"event_severity": "HIGH",
			"environment":    []string{"utkarsh_crapi"},
			"regions":        []string{"afghanistan", "albania"},
			// "regions":        []string{"", "62037528-0793-56ad-a719-509dab66c515"},
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":      os.Getenv("PLATFORM_URL"),
		},
	}

	fmt.Println("Starting terraform init and apply for Region Rule Alert test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created Region Rule Alert resource")
	output := terraform.OutputMap(t, terraformOptions, "traceable_region_rule_alert")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_region_rule_alert_test", "Resource name should match test input")

	// Update test: change description and event severity
	terraformOptions.Vars["description"] = "Updated Region Rule Alert description"
	terraformOptions.Vars["event_severity"] = "MEDIUM"
	terraform.Apply(t, terraformOptions)

	fmt.Println("Verifying updated Region Rule Alert resource")
	output = terraform.OutputMap(t, terraformOptions, "traceable_region_rule_alert")
	require.Equal(t, "MEDIUM", output["event_severity"], "Event severity should be updated to MEDIUM")
	require.Equal(t, "Updated Region Rule Alert description", output["description"], "Description should be updated")
	require.NotEmpty(t, output["id"], "Resource ID should still be present")

	terraform.Destroy(t, terraformOptions)
}
