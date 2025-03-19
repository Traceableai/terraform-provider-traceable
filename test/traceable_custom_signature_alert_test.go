package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)

func TestCustomSignatureAlertRRSingle(t *testing.T) {
 

  for _,rr :=range customRRSingletc {

    t.Run(fmt.Sprint(rr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_alert",
  
          Vars: map[string]interface{}{
      "name":            "tf_automation_custom_alert-re",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
      "alert_severity":  testdata.SeverityOptions[0],
      "request_payload_single_valued_conditions":rr,
      "request_payload_multi_valued_conditions":customRRMultitc[0],	
      "attribute_based_conditions":customAttributetc[0],		
      "custom_sec_rule": testdata.CustomSecRuleOptions	,
      "traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
      "platform_url":      os.Getenv("PLATFORM_URL"),
         "inject_request_headers": injectRequestHeadertc[0],


    },
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
   
    terraform.Destroy(t, terraformOptions)

  })


}

}

func TestCustomSignatureAlertRRMulti(t *testing.T) {
  
  for _,rr :=range customRRMultitc {

    t.Run(fmt.Sprint(rr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_alert",
  
          Vars: map[string]interface{}{
      "name":            "tf_automation_custom_alert",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
      "alert_severity":  testdata.SeverityOptions[0],
      "request_payload_single_valued_conditions":customRRSingletc[0],
      "request_payload_multi_valued_conditions":rr,	
      "attribute_based_conditions":customAttributetc[0],		
      "custom_sec_rule": testdata.CustomSecRuleOptions	,
      "traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
      "platform_url":      os.Getenv("PLATFORM_URL"),
         "inject_request_headers": injectRequestHeadertc[0],


    },
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
 
    terraform.Destroy(t, terraformOptions)

  })


}
}


func TestCustomSignatureAlertAttr(t *testing.T) {
 
  for _,attr :=range customAttributetc {

    t.Run(fmt.Sprint(attr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_alert",
  
          Vars: map[string]interface{}{
      "name":            "tf_automation_custom_alert",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
      "alert_severity":  testdata.SeverityOptions[0],
      "request_payload_single_valued_conditions":customRRSingletc[0],
      "request_payload_multi_valued_conditions":customRRMultitc[0],	
      "attribute_based_conditions":attr,		
      "custom_sec_rule": testdata.CustomSecRuleOptions	,
      "traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
      "platform_url":      os.Getenv("PLATFORM_URL"),
      "inject_request_headers": injectRequestHeadertc[0],


    },
  }
  
    logger.Log("Starting terraform init and apply")
    terraform.InitAndApply(t, terraformOptions)

    logger.Log("Verifying created resources")
   
    terraform.Destroy(t, terraformOptions)

  })


}
}