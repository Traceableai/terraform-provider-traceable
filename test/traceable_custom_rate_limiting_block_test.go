package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)
// var logger=testlog.Logger






// var ipaddresstc=generateIPAddresstc()
// var thresholdConfigtc=generateThresholdConfigstc()
// var rrSingleValuedtc=generateRRSingleValued()
// var rrMultiValuedConditionstc=generateRRMultiValued()
// var attributeConditionstc=generateAttributeConditions()
// var ipLocationstc    =generateIPLocationtc()
// var ipOrganisationtc      =generateIPOrgtc()
// var ipAsnConditiontc                        =generateIPAsntc()
// var ipConnectionType=generateIPConnectiontc()
// var regionstc=generateRegionstc()
// var scannerTypetc=generateRequestScannerTypetc()
// var  emailDomaintc   =generateEmailDomaintc()
// var userAgenttc=generateUserAgentstc()





func TestRateLimitingIpAddress(t *testing.T) {

	  for _,ipaddress :=range ipaddresstc {

			t.Run(fmt.Sprint(ipaddress), func(t *testing.T) {


			terraformOptions := &terraform.Options{
				TerraformDir: "../examples/resources/traceable_rate_limiting_block",
		
						Vars: map[string]interface{}{
        "name":            "tf_automation_rl_block_test ",
				"environments":    []string{"utkarsh_21"},
				"enabled":         testdata.Enabled[0],
				"alert_severity":  testdata.SeverityOptions[0],
				"ip_address":      ipaddress,
				"threshold_configs": thresholdConfigtc[0],
				"request_response_single_valued_conditions":rrSingleValuedtc[0],
				"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
				"attribute_based_conditions":attributeConditionstc[0],
			   "label_id_scope"   : []string{"label1", "label2"},
				"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
				"platform_url":      os.Getenv("PLATFORM_URL"),

			},
		}
		
			logger.Log("Starting terraform init and apply")
			terraform.InitAndApply(t, terraformOptions)
	
			logger.Log("Verifying created resources")
			output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
			require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
			require.NotEmpty(t,output["id"])
			terraformOptions.Vars["enabled"] = testdata.Enabled[1]
      terraform.Apply(t, terraformOptions)
			output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
			require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
      require.NotEmpty(t,output["id"])
			terraform.Destroy(t, terraformOptions)

		})

	}

}


func TestRateLimitingIpAsn(t *testing.T) {

	// fmt.Println(ipAsnConditiontc)


	for _,ipasn :=range  ipAsnConditiontc {

		t.Run(fmt.Sprint(ipasn), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_asn":      ipasn,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],
      "endpoint_id_scope" : []string{"endpoint1", "endpoint2"},						
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}

}

func TestRateLimitingIpConnection(t *testing.T) {

	// fmt.Println(ipAsnConditiontc)


	for _,conn:=range  ipConnectionType {

		t.Run(fmt.Sprint(conn), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_connection_type":      conn,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}

}

func TestRateLimitingRegions(t *testing.T) {



	for _,region:=range  regionstc {

		t.Run(fmt.Sprint(region), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"regions":      region,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}

}

func TestRateLimitingScanner(t *testing.T) {
	for _,scan:=range  scannerTypetc {

		t.Run(fmt.Sprint(scan), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"request_scanner_type":      scan,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}
}


func TestRateLimitingEmailDomain(t *testing.T) {



	for _,tc:=range  emailDomaintc {

		t.Run(fmt.Sprint(tc), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"email_domain":      tc,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}
}


func TestRateLimitingUserAgent(t *testing.T) {



	for _,tc:=range  userAgenttc {

		t.Run(fmt.Sprint(tc), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"user_agents":      tc,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}
}









	


func TestRateLimitingIpOrganisation(t *testing.T) {

	fmt.Println(ipOrganisationtc  )


	for _,iporg :=range ipOrganisationtc  {

		t.Run(fmt.Sprint(iporg), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_organisation":      iporg,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})


}

}

func TestRateLimitingIpLocation(t *testing.T) {

	for _,iplocation := range ipLocationstc{

		t.Run(fmt.Sprint(iplocation), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_location_type":      iplocation,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}



func TestRateLimitingIpReputation(t *testing.T) {

	for _,ipreputation := range testdata.IPReputationOptions{

		t.Run(fmt.Sprint(ipreputation), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_reputation":      ipreputation,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}


func TestRateLimitingIpAbuse(t *testing.T) {

	for _,ipabuse := range testdata.IPAbuseVelocity{

		t.Run(fmt.Sprint(ipabuse), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_abuse_velocity":      ipabuse,
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}



func TestRateLimitingthresholdconfigs(t *testing.T) {

	for _,config := range thresholdConfigtc{

		t.Run(fmt.Sprint(config), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
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

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}

func TestRateLimitingRRSingle(t *testing.T) {

	for _,rr := range rrSingleValuedtc {

		t.Run(fmt.Sprint(rr), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rr,
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}

func TestRateLimitingRRMulti(t *testing.T) {

	for _,rrmulti := range rrMultiValuedConditionstc {

		t.Run(fmt.Sprint(rrmulti), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrmulti,
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}

func TestRateLimitingAttribute(t *testing.T) {

	for _,attr := range attributeConditionstc {

		t.Run(fmt.Sprint(attr), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_rate_limiting_block",
	
					Vars: map[string]interface{}{
			"name":            "tf_automation_rl_block_test ",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_abuse_velocity":      testdata.IPAbuseVelocity[0],
			"threshold_configs": thresholdConfigtc[0],
			"request_response_single_valued_conditions":rrSingleValuedtc[0],
			"request_response_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"attribute_based_conditions":attr,			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),

		},
	}
	


		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_rate_limiting_block")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})

}

}






		
	




			
		
		


		


	
	

	

	
























