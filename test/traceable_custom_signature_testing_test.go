package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	// "github.com/stretchr/testify/require"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)





func TestCustomSignatureTestingRRSingle(t *testing.T) {
  for _,rr :=range customRRSingletc {

    t.Run(fmt.Sprint(rr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_testing",
  
          Vars: map[string]interface{}{
            "name":            "tf_automation_custom_testing",
            "environments":    []string{"utkarsh_21"},
            "description":         "hii nice day",
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







func TestCustomSignatureTestingRRMulti(t *testing.T) {
  // fmt.Println("Custom Attribute CC")
  // fmt.Println(customAttributetc)
  // fmt.Println(testdata.MatchOperatorskey["key"])
  // fmt.Println(testdata.MatchOperatorskey["Value"])

  for _,rr :=range customRRMultitc {

    t.Run(fmt.Sprint(rr), func(t *testing.T) {
     

      terraformOptions := &terraform.Options{
        TerraformDir: "../examples/resources/traceable_custom_signature_testing",
    
            Vars: map[string]interface{}{
        "name":            "tf_automation_custom_testing",
        "environments":    []string{"utkarsh_21"},
        "description":         "hii nice day",
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



func TestCustomSignatureTestingAttr(t *testing.T) {
  for _,attr :=range customAttributetc {

    t.Run(fmt.Sprint(attr), func(t *testing.T) {
     

    terraformOptions := &terraform.Options{
      TerraformDir: "../examples/resources/traceable_custom_signature_testing",
  
          Vars: map[string]interface{}{
      "name":            "tf_automation_custom_testing",
      "environments":    []string{"utkarsh_21"},
      "description":         "hii nice day",
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