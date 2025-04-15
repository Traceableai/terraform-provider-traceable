package test

import (
	"github.com/traceableai/terraform-provider-traceable/test/testdata"
	"github.com/traceableai/terraform-provider-traceable/test/testlog"
	"fmt"

)

func generateIPAddresstc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 

	addressesList := [][]string{{"1.1.1.1", "1.1.1.1/32"},{"1.1.1.1"}}
	excludeOptions := []bool{true, false}

	for _, ipList := range addressesList {
		var temp []map[string]interface{} 
		for _, exclude := range excludeOptions {
			temp = []map[string]interface{}{
				{
				"ip_address_list": ipList,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}

func generateThresholdConfigstc() [][]map[string]interface{} {
	var results [][]map[string]interface{}
	var durations =map[string][]string{
		"DYNAMIC":[]string{"PT86400S",},
		"ROLLING_WINDOW":[]string{"",},
	}

	for _, apitype := range testdata.ApiAggregateType {
		for _ ,usertype :=range testdata.UserAggregateType{
			for _,config :=range testdata.ThresholdConfigType{
				for _,duration :=range durations[config] {

		temp := []map[string]interface{}{
			{
			"api_aggregate_type": apitype,
			"user_aggregate_type":usertype,
			"rolling_window_count_allowed":10,
			"rolling_window_duration":"PT60s",
			"threshold_config_type":config,
			"dynamic_mean_calculation_duration": duration,

		},
	}
	results = append(results, temp)
	}
}
	}
}
return results
}


func generateRRSingleValued() [][]map[string]interface{} {
	var results [][]map[string]interface{} 



	for _, loc := range testdata.RRSingle {
		for _,operator :=range testdata.MatchOperatorskey[loc]{

			temp := []map[string]interface{}{
				{
				"request_location": loc,
				"operator":operator,
				"value":"200",
			},
		}
		results = append(results, temp) 
		}
	}

	return results
}

func generateRRMultiValued() [][]map[string]interface{} {
	var results [][]map[string]interface{} 



	for _, loc := range testdata.RRMulti {
		for _,operator :=range testdata.MatchOperatorskey[loc]{

			temp := []map[string]interface{}{
				{
				"request_location": loc,
				"key_patterns":[]map[string]interface{}{{
					"operator":operator,
					"value":"200",
				}},
				"value_patterns":[]map[string]interface{}{{
					"operator":operator,
					"value":"200",
				}},
			},
		}
		results = append(results, temp) 
		}
	}

	return results
}

func generateAttributeConditions()[][]map[string]interface{}{
	var results [][]map[string]interface{} 

	for _,operator :=range testdata.MatchOperatorOptions{

		temp := []map[string]interface{}{
			
				{
					"key_condition_operator": operator,
					"key_condition_value":"200",
					"value_condition_operator":operator,
					"value_condition_value":"200",
				},
			
	}
	results = append(results, temp) 
	}

return results

}


func generateIPLocationtc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 

	locationList := [][]string{{"BOT"},{"TOR_EXIT_NODE"},{"SCANNER"},{"PUBLIC_PROXY"},{"HOSTING_PROVIDER"},{"ANONYMOUS_VPN"},{"BOT","SCANNER"}}


	for _, loc := range locationList {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"ip_location_types": loc,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}

func generateUserAgentstc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	userAgentList := [][]string{{"Mozilla/5.0","curl/7.68.0"}}
	for _, agents := range userAgentList  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"user_agents_list": agents,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}
func generateRegionstc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	regionsList := [][]string{{"AX", "DZ"},{"AX"}}
	for _, regions := range regionsList  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"regions_ids": regions,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}
func generateIPOrgtc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	orgsList := [][]string{{"ExampleOrg"},{"ExampleOrg1","ExampleOrg2"}}
	for _, orgs := range orgsList  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"ip_organisation_regexes": orgs,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}
func generateIPAsntc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	asnList := [][]string{{"12345"},{"12345","13678"}}
	for _, asn := range asnList  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"ip_asn_regexes": asn,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}
func generateIPConnectiontc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	connections := [][]string{{"RESIDENTIAL","MOBILE","CORPORATE"},{"DATA_CENTER","EDUCATION"}}
	for _, connection := range connections  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"ip_connection_type_list": connection,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}
func generateRequestScannerTypetc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	scanners := [][]string{{"Rapid7 InsightAppSec","Qualys"},{"Qualys"}}
	for _, scans := range scanners  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"scanner_types_list": scans,
				"exclude":         exclude,
			},
		}
		}
		results = append(results, temp) 
	}

	return results
}
func generateUserIdtc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	userRegexes := [][]string{{"user.*"},{"user.*","enduser.*"}}
	userIds:=[][]string{{"user123", "user456"},{"user123"}}
	for _,regexes :=range userRegexes{
	for _, ids := range userIds  {
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"user_id_regexes": regexes,
				"user_ids": ids,
				"exclude": exclude,
			},
		}
		}
		results = append(results, temp) 
	}
}

	return results
}
func generateEmailDomaintc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 
	emailRegexes := [][]string{{".*@example.com"},{".*@example.com",".*@example1.com"}}
	for _,regexes :=range emailRegexes{
		var temp []map[string]interface{} 
		for _, exclude := range testdata.ExcludeOptions {
			temp = []map[string]interface{}{
				{
				"email_domain_regexes": regexes,
				"exclude": exclude,
			},
		}
		}
		results = append(results, temp) 
	
}

	return results
}




func generateEnumerationthresholdConfigtc() [][]map[string]interface{} {
	var results [][]map[string]interface{}
	var durations =[]string{"PT60S"}
	var sesnsitiveParams =map[string][]string{
		"PATH_PARAMS":{""},
		"REQUEST_BODY":{""},
		"SENSITIVE_PARAMS":{"SELECTED_DATA_TYPES","ALL"},

	}

	for _, apitype := range testdata.ApiAggregateType {
		for _ ,usertype :=range testdata.UserAggregateType{
			for _,config :=range testdata.EnumerationThresholdConfigtype{
				for _,duration :=range durations {
					for _,params :=range sesnsitiveParams[config]{

		temp := []map[string]interface{}{
			{
			"api_aggregate_type": apitype,
			"user_aggregate_type":usertype,
			"unique_values_allowed":10,
			"duration":duration,
			"threshold_config_type":config,
			"sensitive_param_evaluation_type":params ,

		},
	}
	results = append(results, temp)
	}
}
	}
}
	}
return results

}

func generateEnumerationDataLoctc() [][]map[string]interface{} {
	var results [][]map[string]interface{}

	for _,loc := range testdata.DataLocation{
	 	temp := []map[string]interface{}{
			{
			"data_location":loc ,		
		},
	}
	results = append(results, temp)
	}
	return results
}


func generateDlpDataTypeConditionstc()  [][]map[string]interface{} {
	var results [][]map[string]interface{
		
	}

	for _,enabled  :=range []bool{true}{
		for attr := range testdata.DlpCustomLocationAttr{
        for  _,op := range testdata.DlpCustomLocationAttr[attr]{
					temp := []map[string]interface{}{
						{
           
           "custom_location_data_type_key_value_pair_matching":enabled,
					 "custom_location_attribute":attr,
					 "custom_location_attribute_key_operator":op,
					 "custom_location_attribute_value":"sample",

						},
				}
				results = append(results, temp)
		}

	}




}
return results
}


func generateDlpRRSingleValuedtc() [][]map[string]interface{} {
	var results [][]map[string]interface{} 



	for _, loc := range testdata.DlpRRSingle {
		for _,operator :=range testdata.MatchOperatorskey[loc]{

			temp := []map[string]interface{}{
				{
				"request_location": loc,
				"operator":operator,
				"value":"200",
			},
		}
		results = append(results, temp) 
		}
	}

	return results
}



func generateCustomRRSingleValuedtc() [][]map[string]interface{}{
	var results [][]map[string]interface{} 
  categories := []string{"REQUEST","RESPONSE"}

	for _,cat := range categories{
	for _,key := range testdata.CustomKeys[cat]{
		for _,operator :=range testdata.MatchOperatorskey[key]{
        
			temp := []map[string]interface{}{
				{
				"match_category":cat ,
				"match_key":key,
				"match_operator":operator,
				"match_value":"200",
			},


		}
		results = append(results, temp) 

	}
}


}
return results
}

func generateCustomRRMultiValuedtc() [][]map[string]interface{}{
	var results [][]map[string]interface{} 
  categories := []string{"REQUEST","RESPONSE"}

	for _,cat := range categories{
	for _,key := range testdata.CustomMultiKeys[cat]{
		
        
			temp := []map[string]interface{}{
				{
				"match_category":cat ,
				"key_value_tag":key,
				"key_match_operator":testdata.MatchOperatorskey["Key"][0],
				"match_key":"hell0",
				"value_match_operator":testdata.MatchOperatorOptions[0],
				"match_value":"200",
			},


		}

		results = append(results, temp) 

	}
}

return results
}
func generateCustomAttributetc() [][]map[string]interface{}{
	var results [][]map[string]interface{} 

	for _,opkey := range testdata.MatchOperatorskey["Key"]{
	for _,opval := range testdata.MatchOperatorskey["Value"]{
		
        fmt.Println("hello")
			temp := []map[string]interface{}{
				{
				"key_condition_operator":opkey ,
				"key_condition_value":"key",
				"value_condition_operator":opval,
				"value_condition_value":"200",
			
			},


		}

		results = append(results, temp) 

	}
}

return results
}

func generateRequestHeadertc() [][]map[string]interface{}{
	var results [][]map[string]interface{} 
	headerKey:= []string{"key1"}
	headerValue:=[]string{"key2"}

	for _,key := range headerKey{
	for _,val := range headerValue{
		
        fmt.Println("hello")
			temp := []map[string]interface{}{
				{
				"header_key":key ,
				"header_value":val,
			},


		}

		results = append(results, temp) 

	}
}

return results
}























var ipaddresstc=generateIPAddresstc()
var thresholdConfigtc=generateThresholdConfigstc()
var rrSingleValuedtc=generateRRSingleValued()
var rrMultiValuedConditionstc=generateRRMultiValued()
var attributeConditionstc=generateAttributeConditions()
var ipLocationstc    =generateIPLocationtc()
var ipOrganisationtc      =generateIPOrgtc()
var ipAsnConditiontc                        =generateIPAsntc()
var ipConnectionType=generateIPConnectiontc()
var regionstc=generateRegionstc()
var scannerTypetc=generateRequestScannerTypetc()
var  emailDomaintc   =generateEmailDomaintc()
var userAgenttc=generateUserAgentstc()
var enumthresholdConfigtc=generateEnumerationthresholdConfigtc()
var enumdataLoctc=generateEnumerationDataLoctc()
var dlpdataTypeConditionstc=generateDlpDataTypeConditionstc()
var dlpRRSingletc=generateDlpRRSingleValuedtc()
var customRRSingletc=generateCustomRRSingleValuedtc()
var customRRMultitc=generateCustomRRMultiValuedtc()
var customAttributetc=generateCustomAttributetc()
var injectRequestHeadertc=generateRequestHeadertc()


var logger=testlog.Logger



