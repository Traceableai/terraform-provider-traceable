package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)

func TestCustomSignatureAllowRRSingle(t *testing.T) {

  for _,rr :=range customRRSingletc {
		if val, ok := rr[0]["match_category"]; ok && val == "RESPONSE" {
						continue
					}

    t.Run(fmt.Sprint(rr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_allow",
  
          Vars: map[string]interface{}{
     "name":            "tf_automation_custom_allow",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
      "request_payload_single_valued_conditions":rr,
      "request_payload_multi_valued_conditions":customRRMultitc[0],	
      "custom_sec_rule": testdata.CustomSecRuleOptions	,
			"disabled": false,
			"allow_expiry_duration":"PT1H",
      "traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
      "platform_url":      os.Getenv("PLATFORM_URL"),

    },
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
   
    terraform.Destroy(t, terraformOptions)

  })


}

}




func TestCustomSignatureAllowRRMultiple(t *testing.T) {
 

  for _,rr :=range customRRMultitc {
		if val, ok := rr[0]["match_category"]; ok && val == "RESPONSE" {
			continue
		}
    t.Run(fmt.Sprint(rr), func(t *testing.T) {

		
    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_allow",
  
          Vars: map[string]interface{}{
      "name":            "tf_automation_custom_allow",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
      "request_payload_single_valued_conditions":customRRSingletc[0],
      "request_payload_multi_valued_conditions":rr,	
      "custom_sec_rule": testdata.CustomSecRuleOptions	,
			"disabled": false,
			"allow_expiry_duration":"PT1H",
      "traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
      "platform_url":      os.Getenv("PLATFORM_URL"),

    },
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
   
    terraform.Destroy(t, terraformOptions)

  })



}
}