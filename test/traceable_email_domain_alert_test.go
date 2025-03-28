package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestEmailDomainBasic(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_email_domain_alert",
		Vars: map[string]interface{}{
			"name":                    "tf_email_domain_alert_without_headers",
			"description":             "Testing Email Domain Alert without header injection",
			"rule_action":             "ALERT",
			"event_severity":          "HIGH",
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

	fmt.Println("Starting terraform init and apply for Email Domain Alert without headers test")
	terraform.InitAndApply(t, terraformOptions)

	fmt.Println("Verifying created Email Domain Alert resource without headers")
	output := terraform.OutputMap(t, terraformOptions, "traceable_email_domain_alert")
	require.NotEmpty(t, output["id"], "Resource ID should not be empty")
	require.Contains(t, output["name"], "tf_email_domain_alert_without_headers", "Resource name should match test input")

	terraform.Destroy(t, terraformOptions)
}
