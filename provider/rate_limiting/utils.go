package rate_limiting

import (
    "fmt"
    "github.com/traceableai/terraform-provider-traceable/provider/common"
)

func returnConditionsStringRateLimit(
    ip_reputation string,
    ip_abuse_velocity string,
    label_id_scope []interface{},
    endpoint_id_scope []interface{},
    req_res_conditions []interface{},
    attribute_based_conditions []interface{},
    ip_location_type []interface{},
    ip_address []interface{},
    email_domain []interface{},
    user_agents []interface{},
    regions []interface{},
    ip_organisation []interface{},
    ip_asn []interface{},
    ip_connection_type []interface{},
    request_scanner_type []interface{},
    user_id []interface{},
 ) (string,error){
    finalUserIdQuery:=""
        if len(user_id)>0{
            excludeUserId := user_id[0].(map[string]interface{})["exclude"].(bool)
            if userIdRegexes,ok := user_id[0].(map[string]interface{})["user_id_regexes"].([]interface{}) ; ok {
                fmt.Printf("this is len useridregex %d",len(userIdRegexes))
                if len(userIdRegexes)>0{
                    finalUserIdQuery=fmt.Sprintf(USER_ID_REGEXES_QUERY,common.ReturnQuotedStringList(userIdRegexes),excludeUserId)
                }
            }
            if userIds,ok := user_id[0].(map[string]interface{})["user_ids"].([]interface{}); ok {
                fmt.Printf("this is len userid %d",len(userIds))
                if len(userIds)>0{
                    finalUserIdQuery=fmt.Sprintf(USER_ID_LIST_QUERY,common.ReturnQuotedStringList(userIds),excludeUserId)
                }
            }
        }

        finalRequestScannerQuery := ""
        if len(request_scanner_type)>0{
            excludeScannerType := request_scanner_type[0].(map[string]interface{})["exclude"].(bool)
            if scannerTypes,ok := request_scanner_type[0].(map[string]interface{})["scanner_types_list"].([]interface{}); ok {
                finalRequestScannerQuery=fmt.Sprintf(REQUEST_SCANNER_TYPE_QUERY,common.ReturnQuotedStringList(scannerTypes),excludeScannerType)
            }
        }

        finalIpConnectionTypeQuery := ""
        if len(ip_connection_type)>0{
            excludeIpConnectionType := ip_connection_type[0].(map[string]interface{})["exclude"].(bool)
            if ipConnectionTypeList,ok := ip_connection_type[0].(map[string]interface{})["ip_connection_type_list"].([]interface{}); ok {
                finalRequestScannerQuery=fmt.Sprintf(IP_CONNECTION_TYPE_QUERY,common.ReturnQuotedStringList(ipConnectionTypeList),excludeIpConnectionType)
            }
        }
        finalIpAsnQuery := ""
        if len(ip_asn)>0{
            excludeIpAsnType := ip_asn[0].(map[string]interface{})["exclude"].(bool)
            if ipAsnRegexesList,ok := ip_asn[0].(map[string]interface{})["ip_asn_regexes"].([]interface{}); ok {
                finalIpAsnQuery=fmt.Sprintf(IS_ASN_QUERY,common.ReturnQuotedStringList(ipAsnRegexesList),excludeIpAsnType)
            }
        }
        finalIpOrganisationQuery := ""
        if len(ip_organisation)>0{
            excludeIpOrg := ip_organisation[0].(map[string]interface{})["exclude"].(bool)
            if ipOrganisationRegexesList,ok := ip_organisation[0].(map[string]interface{})["ip_organisation_regexes"].([]interface{}); ok {
                finalIpOrganisationQuery=fmt.Sprintf(IP_ORGANISATION_QUERY,common.ReturnQuotedStringList(ipOrganisationRegexesList),excludeIpOrg)
            }
        }
        finalRegionsQuery := ""
        if len(regions)>0{
            excludeRegions := regions[0].(map[string]interface{})["exclude"].(bool)
            if regionIds,ok := regions[0].(map[string]interface{})["regions_ids"].([]interface{}); ok {
                regionIdentifiers := ""
                for _,region := range regionIds{
                    regionIdentifiers+=fmt.Sprintf(`{ countryIsoCode: "%s" }`,region.(string))
                    regionIdentifiers+=","
                }
                regionIdentifiers = regionIdentifiers[:len(regionIdentifiers)-1]
                finalRegionsQuery=fmt.Sprintf(REGION_QUERY,excludeRegions,regionIdentifiers)
            }
        }
        finalUserAgentsQuery := ""
        if len(user_agents)>0{
            excludeUserAgents := user_agents[0].(map[string]interface{})["exclude"].(bool)
            if userAgentsList,ok := user_agents[0].(map[string]interface{})["user_agents_list"].([]interface{}); ok {
                if len(userAgentsList)>0{
                    finalUserAgentsQuery=fmt.Sprintf(USER_AGENT_QUERY,common.ReturnQuotedStringList(userAgentsList),excludeUserAgents)
                }
            }
        }
        finalIpAddressQuery := ""
        if len(ip_address)>0{
            excludeIpAddress := ip_address[0].(map[string]interface{})["exclude"].(bool)
            if IpAddressList,ok := ip_address[0].(map[string]interface{})["ip_address_list"].([]interface{}); ok {
                if len(IpAddressList)>0{
                    finalIpAddressQuery=fmt.Sprintf(RAW_INPUT_IP_ADDRESS_QUERY,common.ReturnQuotedStringList(IpAddressList),excludeIpAddress)
                }

            }
            if ipAddressConditionType,ok := ip_address[0].(map[string]interface{})["ip_address_type"].(string); ok {
                if ipAddressConditionType!=""{
                    finalIpAddressQuery=fmt.Sprintf(ALL_EXTERNAL_IP_ADDRESS_QUERY,ipAddressConditionType,excludeIpAddress)
                }
            }
        }
        finalEmailDomainQuery := ""
        if len(email_domain)>0{
            excludeEmailDomain := email_domain[0].(map[string]interface{})["exclude"].(bool)
            if emailDomainRegexList,ok := email_domain[0].(map[string]interface{})["email_domain_regexes"].([]interface{}); ok {
                finalEmailDomainQuery=fmt.Sprintf(EMAIL_DOMAIN_QUERY,common.ReturnQuotedStringList(emailDomainRegexList),excludeEmailDomain)
            }
        }
        finalIpLocationQuery := ""
        if len(ip_location_type)>0{
            excludeIpLocation := ip_location_type[0].(map[string]interface{})["exclude"].(bool)
            if ipLocationTypeList,ok := ip_location_type[0].(map[string]interface{})["ip_location_types"].([]interface{}); ok {
                finalIpLocationQuery=fmt.Sprintf(IP_LOCATION_TYPE_QUERY,common.ReturnQuotedStringList(ipLocationTypeList),excludeIpLocation)
            }
        }
        finalIpAbuseVelocityQuery := ""
        if ip_abuse_velocity!=""{
            finalIpAbuseVelocityQuery=fmt.Sprintf(IP_ABUSE_VELOCITY_QUERY,ip_abuse_velocity)
        }
        finalIpReputationQuery := ""
        if ip_reputation!=""{
            finalIpReputationQuery=fmt.Sprintf(IP_REPUTATION_QUERY,ip_reputation)
        }
        finalAttributedBasedConditionsQuery := ""
        if len(attribute_based_conditions) > 0{
            valueTemplatedQuery := `valueCondition: { operator: %s, value: "%s" }`
            for _,attBasedCondition := range attribute_based_conditions {
                keyConditionOperator := attBasedCondition.(map[string]interface{})["key_condition_operator"]
                keyConditionValue := attBasedCondition.(map[string]interface{})["key_condition_value"]
                valueConditionOperator := attBasedCondition.(map[string]interface{})["value_condition_operator"]
                valueConditionValue := attBasedCondition.(map[string]interface{})["value_condition_value"]
                // if (valueConditionOperator!="" && valueConditionValue=="") || (valueConditionValue!="" && valueConditionOperator==""){
                //     return "",fmt.Errorf("required both values value_condition_value and value_condition_operator")
                // }
                if valueConditionOperator!="" && valueConditionValue!="" {
                    finalAttributedBasedConditionsQuery+=fmt.Sprintf(ATTRIBUTE_BASED_CONDITIONS_QUERY,keyConditionOperator,keyConditionValue,fmt.Sprintf(valueTemplatedQuery,valueConditionOperator,valueConditionValue))
                }else{
                    finalAttributedBasedConditionsQuery += fmt.Sprintf(ATTRIBUTE_BASED_CONDITIONS_QUERY,keyConditionOperator,keyConditionValue,"")
                }
            }
        }

        finalReqResConditionsQuery := ""
        for _,reqResCondition := range req_res_conditions {
            reqResConditionData := reqResCondition.(map[string]interface{})
            metadata_type := reqResConditionData["metadata_type"].(string)
            req_res_operator := reqResConditionData["req_res_operator"].(string)
            req_res_value := reqResConditionData["req_res_value"].(string)
            finalReqResConditionsQuery += fmt.Sprintf(REQ_RES_CONDITIONS_QUERY,metadata_type,req_res_operator,req_res_value)

        }
        finalScopedQuery := ""
        // if len(endpoint_id_scope)>0 && len(label_id_scope)>0{
        //     return "",fmt.Errorf("required one of endpoint_id_scope or label_id_scope")
        // }
        if len(endpoint_id_scope)>0{
            finalScopedQuery=fmt.Sprintf(ENDPOINT_SCOPED_QUERY,common.ReturnQuotedStringList(endpoint_id_scope))
        }else if len(label_id_scope) > 0{
            finalScopedQuery=fmt.Sprintf(LABEL_ID_SCOPED_QUERY,common.ReturnQuotedStringList(label_id_scope))
        }
    finalConditionsQuery := finalScopedQuery+finalReqResConditionsQuery+finalAttributedBasedConditionsQuery+finalIpReputationQuery+finalIpAbuseVelocityQuery+finalIpLocationQuery+finalEmailDomainQuery+finalIpAddressQuery+finalUserAgentsQuery+finalRegionsQuery+finalIpOrganisationQuery+finalIpAsnQuery+finalIpConnectionTypeQuery+finalRequestScannerQuery+finalUserIdQuery
    return finalConditionsQuery,nil
 }

func returnFinalThresholdConfigQuery(threshold_configs []interface{}) (string,error){
    finalThresholdConfigQuery := ""
    for _, thresholdConfig := range threshold_configs {
		thresholdConfigData := thresholdConfig.(map[string]interface{})
		apiAggregateType := thresholdConfigData["api_aggregate_type"].(string)
		rollingWindowCountAllowed := thresholdConfigData["rolling_window_count_allowed"].(int)
		rollingWindowDuration := thresholdConfigData["rolling_window_duration"].(string)
		thresholdConfigType := thresholdConfigData["threshold_config_type"]
		dynamicMeanCalculationDuration := thresholdConfigData["dynamic_mean_calculation_duration"].(string)
		if thresholdConfigType == "ROLLING_WINDOW" {
			finalThresholdConfigQuery += fmt.Sprintf(ROLLING_WINDOW_THRESHOLD_CONFIG_QUERY, apiAggregateType, rollingWindowCountAllowed, rollingWindowDuration)
		} else if thresholdConfigType == "DYNAMIC" {
			// if dynamicMeanCalculationDuration == "" {
			// 	return "",fmt.Errorf("required dynamic_mean_calculation_duration for dynamic threshold_config_type")
			// }
			finalThresholdConfigQuery += fmt.Sprintf(DYNAMIC_THRESHOLD_CONFIG_QUERY, apiAggregateType, rollingWindowCountAllowed, dynamicMeanCalculationDuration, rollingWindowDuration)
		}
	}
    return finalThresholdConfigQuery,nil
}
