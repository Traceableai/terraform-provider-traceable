package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	// "github.com/stretchr/testify/require"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)

func TestCustomSignatureBlockRRSingle(t *testing.T) {

  for _,rr :=range customRRSingletc {
		if val, ok := rr[0]["match_category"]; ok && val == "RESPONSE" {
						continue
					}

    t.Run(fmt.Sprint(rr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_block",
  
          Vars: map[string]interface{}{
     "name":            "tf_automation_custom_block",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
      "request_payload_single_valued_conditions":rr,
      "request_payload_multi_valued_conditions":customRRMultitc[0],	
      "custom_sec_rule": testdata.CustomSecRuleOptions	,
			"disabled": false,
			"block_expiry_duration":"PT1H",
			"alert_severity":testdata.SeverityOptions[0],
      "traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
      "platform_url":      os.Getenv("PLATFORM_URL"),

    },
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
    // output := terraform.OutputMap(t, terraformOptions, "traceable_custom_signature_alert")

    // require.Equal(t,testdata.SeverityOptions[0],output["alert_severity"])
    // // require.NotEmpty(t,output["id"])
    // // terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
    // // terraform.Apply(t, terraformOptions)
    // // output = terraform.OutputMap(t, terraformOptions, "traceable_custom_signature_alert")
    // // require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
    // // require.NotEmpty(t,output["id"])
    terraform.Destroy(t, terraformOptions)

  })


}

}




func TestCustomSignatureBlokRRMultiple(t *testing.T) {
 

  for _,rr :=range customRRMultitc {
		if val, ok := rr[0]["match_category"]; ok && val == "RESPONSE" {
			continue
		}
    t.Run(fmt.Sprint(rr), func(t *testing.T) {

		
			terraformOptions := &terraform.Options{
				TerraformDir: "../examples/resources/traceable_custom_signature_block",
		
						Vars: map[string]interface{}{
			 "name":            "tf_automation_custom_block",
				"environments":    []string{"utkarsh_21"},
				"description":         "hii nice day",
				"request_payload_single_valued_conditions":customRRSingletc[0],
				"request_payload_multi_valued_conditions":rr,	
				"custom_sec_rule": testdata.CustomSecRuleOptions	,
				"disabled": false,
				"block_expiry_duration":"PT1H",
				"alert_severity":testdata.SeverityOptions[0],
				"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
				"platform_url":      os.Getenv("PLATFORM_URL"),
	
			},
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
    // output := terraform.OutputMap(t, terraformOptions, "traceable_custom_signature_allow")

    // require.Equal(t,"hii nice day",output["description"])
    // // require.NotEmpty(t,output["id"])
    // // terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
    // // terraform.Apply(t, terraformOptions)
    // // output = terraform.OutputMap(t, terraformOptions, "traceable_custom_signature_alert")
    // // require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
    // // require.NotEmpty(t,output["id"])
    terraform.Destroy(t, terraformOptions)

  })



}
}