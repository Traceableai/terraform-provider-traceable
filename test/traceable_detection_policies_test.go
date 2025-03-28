package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)


var ruleConfigTestCases = []map[string]interface{}{
	{
		"environment": "utkarsh_crapi",
		"waap_config": []map[string]interface{}{
			{
				"rule_id": "XSS",
				"rule_config": []map[string]interface{}{
					{
						"disabled": false,
					},
				},
			},
		},
	},
	{
		"environment": "utkarsh_tr",
		"waap_config": []map[string]interface{}{
			{
				"rule_id": "RFI", 
				"rule_config": []map[string]interface{}{
					{
						"disabled": false,
					},
				},
			},
		},
	},
}

var subRuleConfigTestCases = []map[string]interface{}{
    	{
		"environment": "utkarsh_crapi",
		"waap_config": []map[string]interface{}{
			{
				"rule_id": "XSS",
				"subrule_config": []map[string]interface{}{
					{
						"sub_rule_id":    "crs_9410170",
						"sub_rule_action": "BLOCK",
					},
				},
			},
		},
	},
    {
        		"environment": "utkarsh_tr",
        		"waap_config": []map[string]interface{}{
        			{
        				"rule_id": "RFI", 
        				"subrule_config": []map[string]interface{}{
        					{
        						"sub_rule_id":    "crs_931120",
        						"sub_rule_action": "MONITOR",
        					},
        				},
        			},
        		},
        	},
}


func TestDetectionConfigRule(t *testing.T) {
	for i, testCase := range ruleConfigTestCases {
		t.Run(fmt.Sprintf("RuleConfigTestCase-%d", i), func(t *testing.T) {
			terraformOptions := &terraform.Options{
				TerraformDir: "../examples/resources/traceable_detection_policies",
				Vars: map[string]interface{}{
					"environment":       testCase["environment"],
					"waap_config":       testCase["waap_config"],
					"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"), 
					"platform_url":      os.Getenv("PLATFORM_URL"),      
				},
			}
      
			terraform.InitAndApply(t, terraformOptions)

			output := terraform.OutputMap(t, terraformOptions, "traceable_detection_policies")
			require.NotEmpty(t, output["id"])
			require.Equal(t, testCase["environment"], output["environment"]) 

			terraform.Destroy(t, terraformOptions)
		})
	}
}

func TestDetectionConfigSubRule(t *testing.T) {
	for i, testCase := range subRuleConfigTestCases {
		t.Run(fmt.Sprintf("SubRuleConfigTestCase-%d", i), func(t *testing.T) {
			terraformOptions := &terraform.Options{
				TerraformDir: "../examples/resources/traceable_detection_policies",
				Vars: map[string]interface{}{
					"environment":       testCase["environment"],
					"waap_config":       testCase["waap_config"],
					"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
					"platform_url":      os.Getenv("PLATFORM_URL"),      
				},
			}

			terraform.InitAndApply(t, terraformOptions)

			output := terraform.OutputMap(t, terraformOptions, "traceable_detection_policies")
			require.NotEmpty(t, output["id"])
            require.Equal(t, testCase["environment"], output["environment"])
			terraform.Destroy(t, terraformOptions)
		})
	}
}