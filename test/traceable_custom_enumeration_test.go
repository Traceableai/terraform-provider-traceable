package test
import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)




func TestEnumerationAlertIpAddress(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,ipaddress :=range ipaddresstc {

		t.Run(fmt.Sprint(rule,ipaddress), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_enumeration",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_address":      ipaddress,
			"threshold_configs": enumthresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],

		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}
	}

}


func TestEnumerationAlertIpAsn(t *testing.T) {

// fmt.Println(ipAsnConditiontc)

for _,rule :=range testdata.EnumerationRuleType{
for _,ipasn :=range  ipAsnConditiontc {

	t.Run(fmt.Sprint(ipasn), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,			
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_asn":      ipasn,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],	
		"label_id_scope"   : []string{"label1", "label2"},		
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}
logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}
}

}

func TestEnumerationAlertIpConnection(t *testing.T) {

// fmt.Println(ipAsnConditiontc)

for _,rule :=range testdata.EnumerationRuleType{

for _,conn:=range  ipConnectionType {

	t.Run(fmt.Sprint(conn), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_connection_type":      conn,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}
logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}
}

}

func TestEnumerationAlertRegions(t *testing.T) {


	for _,rule :=range testdata.EnumerationRuleType{
for _,region:=range  regionstc {

	t.Run(fmt.Sprint(rule,region), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
	  "rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"regions":      region,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}
logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}
	}

}

func TestEnumerationAlertScanner(t *testing.T) {

  for _,rule :=range testdata.EnumerationRuleType{
for _,scan:=range  scannerTypetc {

	t.Run(fmt.Sprint(rule,scan), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":            rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"request_scanner_type":      scan,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}
logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}
}
}


func TestEnumerationAlertEmailDomain(t *testing.T) {


for _,rule :=range testdata.EnumerationRuleType{

for _,tc:=range  emailDomaintc {

	t.Run(fmt.Sprint(rule,tc), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,			
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"email_domain":      tc,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}
logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}
}
}


func TestEnumerationAlertUserAgent(t *testing.T) {


	for _,rule :=range testdata.EnumerationRuleType{
for _,tc:=range  userAgenttc {

	t.Run(fmt.Sprint(rule,tc), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,	
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"user_agents":      tc,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],
	},
}
logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}
}

}












func TestEnumerationAlertIpOrganisation(t *testing.T) {

fmt.Println(ipOrganisationtc  )

for _,rule :=range testdata.EnumerationRuleType{
for _,iporg :=range ipOrganisationtc  {

	t.Run(fmt.Sprint(rule,iporg), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,				
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_organisation":      iporg,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})


}

}
}

func TestEnumerationAlertIpLocation(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{
for _,iplocation := range ipLocationstc{

	t.Run(fmt.Sprint(rule,iplocation), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_location_type":      iplocation,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}



func TestEnumerationAlertIpReputation(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{

for _,ipreputation := range testdata.IPReputationOptions{

	t.Run(fmt.Sprint(rule,ipreputation), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_reputation":      ipreputation,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}

func TestEnumerationAlertIpAbuse(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{
for _,ipabuse := range testdata.IPAbuseVelocity{

	t.Run(fmt.Sprint(rule,ipabuse), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_abuse_velocity":      ipabuse,
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}



func TestEnumerationAlertthresholdconfigs(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{
for _,config := range enumthresholdConfigtc{

	t.Run(fmt.Sprint(rule,config), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,			
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
		"threshold_configs": config,
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}

func TestEnumerationAlertRRSingle(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{
for _,rr := range rrSingleValuedtc {

	t.Run(fmt.Sprint(rule,rr), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,	
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rr,
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],
	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}

func TestEnumerationAlertRRMulti(t *testing.T) {

	for _,rule :=range testdata.EnumerationRuleType{
for _,rrmulti := range rrMultiValuedConditionstc {

	t.Run(fmt.Sprint(rule,rrmulti), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrmulti,
		"attribute_based_conditions":attributeConditionstc[0],			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],

	},
}
	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")
	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)
	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}

func TestEnumerationAlertAttribute(t *testing.T) {


	for _,rule :=range testdata.EnumerationRuleType{

for _,attr := range attributeConditionstc {

	t.Run(fmt.Sprint(rule,attr), func(t *testing.T) {


	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/resources/traceable_enumeration",

				Vars: map[string]interface{}{
		"rule_type":       rule,
		"name":            "tf_automation_test ",
		"environments":    []string{"utkarsh_21"},
		"enabled":         testdata.Enabled[0],
		"alert_severity":  testdata.SeverityOptions[0],
		"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
		"threshold_configs": enumthresholdConfigtc[0],
		"request_response_single_valued_conditions":rrSingleValuedtc[0],
		"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
		"attribute_based_conditions":attr,			
		"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
		"platform_url":      os.Getenv("PLATFORM_URL"),
		"data_set_name":    testdata.DataSetName[0],
	},
}



	logger.Log("Starting terraform init and apply")
	terraform.InitAndApply(t, terraformOptions)

	logger.Log("Verifying created resources")

	output := terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
	require.NotEmpty(t,output["id"])

	terraformOptions.Vars["enabled"] = testdata.Enabled[1]
	terraform.Apply(t, terraformOptions)

	output = terraform.OutputMap(t, terraformOptions, "traceable_enumeration")
	require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
	require.NotEmpty(t,output["id"])
	terraform.Destroy(t, terraformOptions)

})

}

}
}






	






	
	


	


































