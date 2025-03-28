package test
import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)




// func TestDlpRequestDataTypes(t *testing.T) {


//   for _,rule :=range testdata.EnumerationRuleType{
// 	for _,datatype :=range dlpdataTypeConditionstc {

// 		t.Run(fmt.Sprint(rule,datatype), func(t *testing.T) {


// 		terraformOptions := &terraform.Options{
// 			TerraformDir: "../examples/resources/traceable_dlp_request_based",
	
// 					Vars: map[string]interface{}{
// 		  "rule_type":       rule,
// 			"name":            "tf_automation_test ",
// 			"environments":    []string{"bot-protection-demo","utkarsh_21"},
// 			"alert_severity":  testdata.SeverityOptions[0],
// 			"enabled":          testdata.Enabled[0],
// 			"ip_address":         ipaddresstc[0],
// 			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
// 			"platform_url":      os.Getenv("PLATFORM_URL"),
// 			"data_set_name":    testdata.DataSetName[0],
// 			"data_types_conditions": datatype,
// 			"url_regexes":testdata.UrlRegexes[1],
// 			"service_name": "nginx-gateway",


// 		},
// 	}
	
// 		logger.Log("Starting terraform init and apply")
// 		terraform.InitAndApply(t, terraformOptions)

// 		logger.Log("Verifying created resources")
// 		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 		require.Equal(t, testdata.SeverityOptions[0],output["alert_severity"])
// 		require.NotEmpty(t,output["id"])
// 		terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
// 		terraform.Apply(t, terraformOptions)
// 		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 		require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
// 		require.NotEmpty(t,output["id"])
// 		terraform.Destroy(t, terraformOptions)


// 	})



// }
// 	}

// }

// func TestDlpIpAddress(t *testing.T) {
// 	addressesList := [][]string{{"1.1.1.1", "1.1.1.1/32"},{"1.1.1.1"}}


//   for _,rule :=range testdata.EnumerationRuleType{
// 	for _,ipaddress :=range addressesList {

// 		t.Run(fmt.Sprint(rule,ipaddress), func(t *testing.T) {


// 		terraformOptions := &terraform.Options{
// 			TerraformDir: "../examples/resources/traceable_dlp_request_based",
	
// 					Vars: map[string]interface{}{
// 						"rule_type":       rule,
// 						"name":            "tf_automation_test ",
// 						"environments":    []string{"bot-protection-demo","utkarsh_21"},
// 						"alert_severity":  testdata.SeverityOptions[0],
// 						"enabled":          testdata.Enabled[0],
// 						"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
// 						"platform_url":      os.Getenv("PLATFORM_URL"),
// 						"data_set_name":    testdata.DataSetName[0],
// 						"data_types_conditions": dlpdataTypeConditionstc[0],
// 						"url_regexes":testdata.UrlRegexes[1],
// 						"service_name": "nginx-gateway",
// 						"ip_address": ipaddress,

// 		},
// 	}
	
// 	logger.Log("Starting terraform init and apply")
// 	terraform.InitAndApply(t, terraformOptions)

// 	logger.Log("Verifying created resources")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t, testdata.SeverityOptions[0],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
// 	terraform.Apply(t, terraformOptions)
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraform.Destroy(t, terraformOptions)

// 	})

// }
// 	}

// }





// func TestDlpIpRegions(t *testing.T) {
// 	regionsList := [][]string{{"AX", "DZ"},{"AX"}}

//   for _,rule :=range testdata.EnumerationRuleType{
// 	for _,regions :=range regionsList {

// 		t.Run(fmt.Sprint(rule,regions), func(t *testing.T) {


// 		terraformOptions := &terraform.Options{
// 			TerraformDir: "../examples/resources/traceable_dlp_request_based",
	
// 					Vars: map[string]interface{}{
// 						"rule_type":       rule,
// 						"name":            "tf_automation_test ",
// 						"environments":    []string{"bot-protection-demo","utkarsh_21"},
// 						"alert_severity":  testdata.SeverityOptions[0],
// 						"enabled":          testdata.Enabled[0],
// 						"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
// 						"platform_url":      os.Getenv("PLATFORM_URL"),
// 						"data_set_name":    testdata.DataSetName[0],
// 						"data_types_conditions": dlpdataTypeConditionstc[0],
// 						"url_regexes":testdata.UrlRegexes[1],
// 						"service_name": "nginx-gateway",
// 						"regions": regions,

// 		},
// 	}
	
// 	logger.Log("Starting terraform init and apply")
// 	terraform.InitAndApply(t, terraformOptions)

// 	logger.Log("Verifying created resources")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t, testdata.SeverityOptions[0],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
// 	terraform.Apply(t, terraformOptions)
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraform.Destroy(t, terraformOptions)

// 	})

// }
// 	}

// }

// func TestDlpIpLocationType(t *testing.T) {
// 	locationList := [][]string{{"BOT"},{"TOR_EXIT_NODE"},{"SCANNER"},{"PUBLIC_PROXY"},{"HOSTING_PROVIDER"},{"ANONYMOUS_VPN"},{"BOT","SCANNER"}}
//   for _,rule :=range testdata.EnumerationRuleType{
// 	for _,loc :=range locationList {

// 		t.Run(fmt.Sprint(rule,loc), func(t *testing.T) {


// 		terraformOptions := &terraform.Options{
// 			TerraformDir: "../examples/resources/traceable_dlp_request_based",
	
// 					Vars: map[string]interface{}{
// 						"rule_type":       rule,
// 						"name":            "tf_automation_test ",
// 						"environments":    []string{"bot-protection-demo","utkarsh_21"},
// 						"alert_severity":  testdata.SeverityOptions[0],
// 						"enabled":          testdata.Enabled[0],
// 						"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
// 						"platform_url":      os.Getenv("PLATFORM_URL"),
// 						"data_set_name":    testdata.DataSetName[0],
// 						"data_types_conditions": dlpdataTypeConditionstc[0],
// 						"url_regexes":testdata.UrlRegexes[1],
// 						"service_name": "nginx-gateway",
// 						"ip_location_type": loc,

// 		},
// 	}
	
// 	logger.Log("Starting terraform init and apply")
// 	terraform.InitAndApply(t, terraformOptions)

// 	logger.Log("Verifying created resources")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t, testdata.SeverityOptions[0],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
// 	terraform.Apply(t, terraformOptions)
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraform.Destroy(t, terraformOptions)

// 	})

// }
// 	}

// }



// func TestEnumerationAlertRRSingle(t *testing.T) {

// 	for _,rule :=range testdata.EnumerationRuleType{
// for _,rr := range dlpRRSingletc {

// 	t.Run(fmt.Sprint(rule,rr), func(t *testing.T) {


// 	terraformOptions := &terraform.Options{
// 		TerraformDir: "../examples/resources/traceable_dlp_request_based",

// 		Vars: map[string]interface{}{
// 			"rule_type":       rule,
// 			"name":            "tf_automation_test ",
// 			"environments":    []string{"bot-protection-demo","utkarsh_21"},
// 			"alert_severity":  testdata.SeverityOptions[0],
// 			"enabled":          testdata.Enabled[0],
// 			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
// 			"platform_url":      os.Getenv("PLATFORM_URL"),
// 			"data_set_name":    testdata.DataSetName[0],
// 			"data_types_conditions": dlpdataTypeConditionstc[0],
// 			"url_regexes":testdata.UrlRegexes[1],
// 			"service_name": "nginx-gateway",
// 			"request_payload_single_valued_conditions": rr,
			

// },
// }



// logger.Log("Starting terraform init and apply")
// 	terraform.InitAndApply(t, terraformOptions)

// 	logger.Log("Verifying created resources")
// 	output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t, testdata.SeverityOptions[0],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
// 	terraform.Apply(t, terraformOptions)
// 	output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
// 	require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
// 	require.NotEmpty(t,output["id"])
// 	terraform.Destroy(t, terraformOptions)
// })

// }

// }
// }


func TestDlpRRMulti(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{
for _,rr := range rrMultiValuedConditionstc {

	t.Run(fmt.Sprint(rule,rr), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_dlp_request_based",

		Vars: map[string]interface{}{
			"rule_type":       rule,
			"name":            "tf_automation_test ",
			"environments":    []string{"bot-protection-demo","utkarsh_21"},
			"alert_severity":  testdata.SeverityOptions[0],
			"enabled":          testdata.Enabled[0],
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_types_conditions": dlpdataTypeConditionstc[0],
			"url_regexes":testdata.UrlRegexes[1],
			"service_name": "nginx-gateway",
			"request_payload_multi_valued_conditions": rr,
			

},
}



logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
	require.Equal(t, testdata.SeverityOptions[0],output["alert_severity"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["alert_severity"] = testdata.SeverityOptions[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_request_based")
	require.Equal(t,testdata.SeverityOptions[1],output["alert_severity"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)
})
break

}

}
}


