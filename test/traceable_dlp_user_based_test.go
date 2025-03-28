package test

import (
	"fmt"
	"os"
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
)



func TestDlpUseBasedIpAddress(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,ipaddress :=range ipaddresstc {

		t.Run(fmt.Sprint(rule,ipaddress), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_address":      ipaddress,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}


func TestDlpUseBasedIpLocation(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,iplocation :=range ipLocationstc {

		t.Run(fmt.Sprint(rule,iplocation), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_location_type":     iplocation,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}



func TestDlpUseBasedIpReputation(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,ipReputation :=range testdata.IPReputationOptions {

		t.Run(fmt.Sprint(rule,ipReputation), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_reputation":     ipReputation,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}


func TestDlpUseBasedIpAbuseVelocity(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,ipAbuseVelocity :=range testdata.IPAbuseVelocity {

		t.Run(fmt.Sprint(rule,ipAbuseVelocity), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_abuse_velocity":     ipAbuseVelocity,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}



func TestDlpUseBasedEmailDomain(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,emailDomain :=range emailDomaintc {

		t.Run(fmt.Sprint(rule,emailDomain), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"email_domain":     emailDomain,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}



func TestDlpUseBasedRegion(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,region :=range regionstc {

		t.Run(fmt.Sprint(rule,region), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"regions":     region,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}

func TestDlpUseBasedIpOrganisation(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,ipOrganisation :=range ipOrganisationtc {

		t.Run(fmt.Sprint(rule,ipOrganisation), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_organisation":     ipOrganisation,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}

func TestDlpUseBasedIpAsn(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,ipAsn :=range ipAsnConditiontc {

		t.Run(fmt.Sprint(rule,ipAsn), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_asn":     ipAsn,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}



func TestDlpUseBasedIpConnection(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,conn :=range ipConnectionType {

		t.Run(fmt.Sprint(rule,conn), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"ip_connection_type":     conn,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}


func TestDlpUseBasedReqScanner(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,scan :=range scannerTypetc {

		t.Run(fmt.Sprint(rule,scan), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"request_scanner_type":     scan,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}

func TestDlpUseBasedUserId(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,userId :=range userIdtc {

		t.Run(fmt.Sprint(rule,userId), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"user_id":     userId,
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}

func TestDlpUseBasedRRSingle(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,rr :=range rrSingleValuedtc {

		t.Run(fmt.Sprint(rule,rr), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"user_id":     userIdtc[0],
			"request_payload_single_valued_conditions":rr,
			"request_payload_multi_valued_conditions":rrMultiValuedConditionstc[0],
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}

func TestDlpUseBasedRRMulti(t *testing.T) {


  for _,rule :=range testdata.EnumerationRuleType{
	for _,rr :=range rrMultiValuedConditionstc {

		t.Run(fmt.Sprint(rule,rr), func(t *testing.T) {


		terraformOptions := &terraform.Options{
			TerraformDir: "../examples/resources/traceable_dlp_user_based",
	
					Vars: map[string]interface{}{
		  "rule_type":       rule,
			"name":            "tf_automation_dlp_user_based_test",
			"environments":    []string{"utkarsh_21"},
			"enabled":         testdata.Enabled[0],
			"alert_severity":  testdata.SeverityOptions[0],
			"user_id":     userIdtc[0],
			"request_payload_single_valued_conditions":rrSingleValuedtc[0],
			"request_payload_multi_valued_conditions":rr,
			"dynamic_threshold_config": dlpDynamicThresholdConfigtc[0],
			// "value_based_threshold_config":dlpValueThresholdConfigtc[0],
			"rolling_window_threshold_config":dlpRollingThresholdConfigtc[0],
			"endpoint_id_scope" : []string{"endpoint1", "endpoint2"},			
			"attribute_based_conditions":attributeConditionstc[0],			
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),		
			"platform_url":      os.Getenv("PLATFORM_URL"),
			"data_set_name":    testdata.DataSetName[0],
			"data_location":    datalocationtc[0],


		},
	}
	
		logger.Log("Starting terraform init and apply")
		terraform.InitAndApply(t, terraformOptions)

		logger.Log("Verifying created resources")
		output := terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t, fmt.Sprintf("%t",testdata.Enabled[0]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraformOptions.Vars["enabled"] = testdata.Enabled[1]
		terraform.Apply(t, terraformOptions)
		output = terraform.OutputMap(t, terraformOptions, "traceable_dlp_user_based")
		require.Equal(t,fmt.Sprintf("%t",testdata.Enabled[1]),output["enabled"])
		require.NotEmpty(t,output["id"])
		terraform.Destroy(t, terraformOptions)

	})
	

}
	}

}









