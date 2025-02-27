package rate_limiting

import (
	"fmt"

	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ReturnConditionsStringRateLimit(
	ipReputation string,
	ipAbuseVelocity string,
	labelIdScope []interface{},
	endpointIdScope []interface{},
	requestResponseSingleValuedConditions []interface{},
	requestResponseMultiValuedConditions []interface{},
	attributeBasedConditions []interface{},
	ipLocationType []interface{},
	ipAddress []interface{},
	emailDomain []interface{},
	userAgents []interface{},
	regions []interface{},
	ipOrganisation []interface{},
	ipAsn []interface{},
	ipConnectionType []interface{},
	requestScannerType []interface{},
	userId []interface{},
	dataTypesConditions []interface{},
) (string, error) {
	var (
		finalUserIdQuery                                string
		finalRequestScannerQuery                        string
		finalIpConnectionTypeQuery                      string
		finalIpAsnQuery                                 string
		finalIpOrganisationQuery                        string
		finalRegionsQuery                               string
		finalUserAgentsQuery                            string
		finalIpAddressQuery                             string
		finalEmailDomainQuery                           string
		finalIpLocationQuery                            string
		finalIpAbuseVelocityQuery                       string
		finalIpReputationQuery                          string
		finalAttributedBasedConditionsQuery             string
		finalRequestResponseSingleValuedConditionsQuery string
		finalRequestResponseMultiValuedConditionsQuery  string
		finalScopedQuery                                string
		finalDataTypesConditionsQuery                   string
	)

	if len(userId) > 0 {

		userIdData := userId[0].(map[string]interface{})
		excludeUserId := userIdData["exclude"].(bool)
		if userIdRegexes, ok := userIdData["user_id_regexes"].([]interface{}); ok && len(userIdRegexes) > 0 {
			finalUserIdQuery = fmt.Sprintf(USER_ID_REGEXES_QUERY, common.ReturnQuotedStringList(userIdRegexes), excludeUserId)
		} else if userIds, ok := userIdData["user_ids"].([]interface{}); ok && len(userIds) > 0 {
			finalUserIdQuery = fmt.Sprintf(USER_ID_LIST_QUERY, common.ReturnQuotedStringList(userIds), excludeUserId)
		}
	}

	if len(requestScannerType) > 0 {

		scannerData := requestScannerType[0].(map[string]interface{})
		excludeScannerType := scannerData["exclude"].(bool)
		if scannerTypes, ok := scannerData["scanner_types_list"].([]interface{}); ok && len(scannerTypes) > 0 {
			finalRequestScannerQuery = fmt.Sprintf(REQUEST_SCANNER_TYPE_QUERY, common.ReturnQuotedStringList(scannerTypes), excludeScannerType)
		}
	}

	if len(ipConnectionType) > 0 {

		connectionData := ipConnectionType[0].(map[string]interface{})
		excludeIpConnectionType := connectionData["exclude"].(bool)
		if connectionTypes, ok := connectionData["ip_connection_type_list"].([]interface{}); ok && len(connectionTypes) > 0 {
			finalIpConnectionTypeQuery = fmt.Sprintf(IP_CONNECTION_TYPE_QUERY, common.ReturnQuotedStringList(connectionTypes), excludeIpConnectionType)
		}
	}

	if len(ipAsn) > 0 {

		asnData := ipAsn[0].(map[string]interface{})
		excludeIpAsnType := asnData["exclude"].(bool)
		if asnRegexes, ok := asnData["ip_asn_regexes"].([]interface{}); ok && len(asnRegexes) > 0 {
			finalIpAsnQuery = fmt.Sprintf(IS_ASN_QUERY, common.ReturnQuotedStringList(asnRegexes), excludeIpAsnType)
		}
	}

	if len(ipOrganisation) > 0 {

		orgData := ipOrganisation[0].(map[string]interface{})
		excludeIpOrg := orgData["exclude"].(bool)
		if orgRegexes, ok := orgData["ip_organisation_regexes"].([]interface{}); ok && len(orgRegexes) > 0 {
			finalIpOrganisationQuery = fmt.Sprintf(IP_ORGANISATION_QUERY, common.ReturnQuotedStringList(orgRegexes), excludeIpOrg)
		}
	}

	if len(regions) > 0 {

		regionData := regions[0].(map[string]interface{})
		excludeRegions := regionData["exclude"].(bool)
		if regionIds, ok := regionData["regions_ids"].([]interface{}); ok && len(regionIds) > 0 {
			regionIdentifiers := ""
			for _, region := range regionIds {
				regionIdentifiers += fmt.Sprintf(`{ countryIsoCode: "%s" },`, region.(string))
			}
			regionIdentifiers = regionIdentifiers[:len(regionIdentifiers)-1]
			finalRegionsQuery = fmt.Sprintf(REGION_QUERY, excludeRegions, regionIdentifiers)
		}
	}

	if len(userAgents) > 0 {

		uaData := userAgents[0].(map[string]interface{})
		excludeUserAgents := uaData["exclude"].(bool)
		if uaList, ok := uaData["user_agents_list"].([]interface{}); ok && len(uaList) > 0 {
			finalUserAgentsQuery = fmt.Sprintf(USER_AGENT_QUERY, common.ReturnQuotedStringList(uaList), excludeUserAgents)
		}
	}

	if len(ipAddress) > 0 {

		ipData := ipAddress[0].(map[string]interface{})
		excludeIpAddress := ipData["exclude"].(bool)
		if ipList, ok := ipData["ip_address_list"].([]interface{}); ok && len(ipList) > 0 {
			finalIpAddressQuery = fmt.Sprintf(RAW_INPUT_IP_ADDRESS_QUERY, common.ReturnQuotedStringList(ipList), excludeIpAddress)
		} else if ipType, ok := ipData["ip_address_type"].(string); ok && ipType != "" {
			finalIpAddressQuery = fmt.Sprintf(ALL_EXTERNAL_IP_ADDRESS_QUERY, ipType, excludeIpAddress)
		}
	}

	if len(emailDomain) > 0 {

		emailData := emailDomain[0].(map[string]interface{})
		excludeEmailDomain := emailData["exclude"].(bool)
		if emailRegexes, ok := emailData["email_domain_regexes"].([]interface{}); ok && len(emailRegexes) > 0 {
			finalEmailDomainQuery = fmt.Sprintf(EMAIL_DOMAIN_QUERY, common.ReturnQuotedStringList(emailRegexes), excludeEmailDomain)
		}
	}

	if len(ipLocationType) > 0 {

		locationData := ipLocationType[0].(map[string]interface{})
		excludeIpLocation := locationData["exclude"].(bool)
		if locationTypes, ok := locationData["ip_location_types"].([]interface{}); ok && len(locationTypes) > 0 {
			finalIpLocationQuery = fmt.Sprintf(IP_LOCATION_TYPE_QUERY, common.ReturnQuotedStringList(locationTypes), excludeIpLocation)
		}
	}

	if ipAbuseVelocity != "" {

		finalIpAbuseVelocityQuery = fmt.Sprintf(IP_ABUSE_VELOCITY_QUERY, ipAbuseVelocity)
	}

	if ipReputation != "" {

		finalIpReputationQuery = fmt.Sprintf(IP_REPUTATION_QUERY, ipReputation)
	}

	if len(attributeBasedConditions) > 0 {

		valueTemplatedQuery := `valueCondition: { operator: %s, value: "%s" }`
		for _, condition := range attributeBasedConditions {
			condData := condition.(map[string]interface{})
			keyOperator := condData["key_condition_operator"].(string)
			keyValue := condData["key_condition_value"].(string)
			if valueOperator, ok := condData["value_condition_operator"].(string); ok && valueOperator != "" {
				valueValue := condData["value_condition_value"].(string)
				finalAttributedBasedConditionsQuery += fmt.Sprintf(ATTRIBUTE_BASED_CONDITIONS_QUERY, keyOperator, keyValue, fmt.Sprintf(valueTemplatedQuery, valueOperator, valueValue))
			} else {
				finalAttributedBasedConditionsQuery += fmt.Sprintf(ATTRIBUTE_BASED_CONDITIONS_QUERY, keyOperator, keyValue, "")
			}
		}
	}

	if len(requestResponseSingleValuedConditions) > 0 {

		for _, requestPayloadSingleValuedCondition := range requestResponseSingleValuedConditions {
			requestPayloadSingleValuedConditionData := requestPayloadSingleValuedCondition.(map[string]interface{})
			requestLocation := requestPayloadSingleValuedConditionData["request_location"].(string)
			keyOp := requestPayloadSingleValuedConditionData["operator"].(string)
			keyValue := requestPayloadSingleValuedConditionData["value"].(string)
			finalRequestResponseSingleValuedConditionsQuery += fmt.Sprintf(REQ_RES_CONDITIONS_QUERY, requestLocation, keyOp, keyValue)
		}

	}
	if len(requestResponseMultiValuedConditions) > 0 {

		for _, requestPayloadMultiValuedCondition := range requestResponseMultiValuedConditions {
			requestPayloadMultiValuedConditionData := requestPayloadMultiValuedCondition.(map[string]interface{})
			requestLocation := requestPayloadMultiValuedConditionData["request_location"].(string)
			keyPatterns := requestPayloadMultiValuedConditionData["key_patterns"].([]interface{})
			keyOp := keyPatterns[0].(map[string]interface{})["operator"]
			keyValue := keyPatterns[0].(map[string]interface{})["value"]
			valQuery := ""
			if valuePatterns, ok := requestPayloadMultiValuedConditionData["value_patterns"].([]interface{}); ok {
				valueOp := valuePatterns[0].(map[string]interface{})["operator"]
				value := valuePatterns[0].(map[string]interface{})["value"]
				valQuery = fmt.Sprintf(DATATYPE_VALUE_CONDITIONS, valueOp, value)
			}
			finalRequestResponseMultiValuedConditionsQuery += fmt.Sprintf(MULTI_VALUES_REQ_CONDITIONS, requestLocation, keyOp, keyValue, valQuery)
		}
	}

	if len(dataTypesConditions) > 0 {

		for _, condition := range dataTypesConditions {
			condData := condition.(map[string]interface{})
			dataTypeIds := condData["data_type_ids"].([]interface{})
			dataSetIds := condData["data_sets_ids"].([]interface{})
			dataLocation := condData["data_location"].(string)
			dataLocationQuery := ""
			if dataLocation != "REQUEST_RESPONSE" {
				dataLocationQuery = fmt.Sprintf(DATA_LOCATION_STRING, dataLocation)
			}

			finalDataTypesConditionsQuery += fmt.Sprintf(DATA_TYPES_CONDITIONS_QUERY, common.InterfaceToStringSlice(dataSetIds), common.InterfaceToStringSlice(dataTypeIds), dataLocationQuery)
		}
	}

	if len(endpointIdScope) > 0 {
		finalScopedQuery = fmt.Sprintf(ENDPOINT_SCOPED_QUERY, common.ReturnQuotedStringList(endpointIdScope))
	} else if len(labelIdScope) > 0 {
		finalScopedQuery = fmt.Sprintf(LABEL_ID_SCOPED_QUERY, common.ReturnQuotedStringList(labelIdScope))
	}

	finalConditionsQuery := finalScopedQuery + finalRequestResponseSingleValuedConditionsQuery + finalRequestResponseMultiValuedConditionsQuery + finalAttributedBasedConditionsQuery +
		finalIpReputationQuery + finalIpAbuseVelocityQuery + finalIpLocationQuery +
		finalEmailDomainQuery + finalIpAddressQuery + finalUserAgentsQuery + finalRegionsQuery +
		finalIpOrganisationQuery + finalIpAsnQuery + finalIpConnectionTypeQuery + finalRequestScannerQuery + finalUserIdQuery + finalDataTypesConditionsQuery

	return finalConditionsQuery, nil
}

func returnFinalThresholdConfigQuery(threshold_configs []interface{}) (string, error) {
	finalThresholdConfigQuery := ""
	for _, thresholdConfig := range threshold_configs {
		thresholdConfigData := thresholdConfig.(map[string]interface{})
		apiAggregateType := thresholdConfigData["api_aggregate_type"].(string)
		userAggregateType := thresholdConfigData["user_aggregate_type"].(string)
		rollingWindowCountAllowed := thresholdConfigData["rolling_window_count_allowed"].(int)
		rollingWindowDuration := thresholdConfigData["rolling_window_duration"].(string)
		thresholdConfigType := thresholdConfigData["threshold_config_type"]
		dynamicMeanCalculationDuration := thresholdConfigData["dynamic_mean_calculation_duration"].(string)
		if thresholdConfigType == "ROLLING_WINDOW" {
			finalThresholdConfigQuery += fmt.Sprintf(ROLLING_WINDOW_THRESHOLD_CONFIG_QUERY, apiAggregateType, userAggregateType, rollingWindowCountAllowed, rollingWindowDuration)
		} else if thresholdConfigType == "DYNAMIC" {
			finalThresholdConfigQuery += fmt.Sprintf(DYNAMIC_THRESHOLD_CONFIG_QUERY, apiAggregateType, userAggregateType, rollingWindowCountAllowed, dynamicMeanCalculationDuration, rollingWindowDuration)
		}
	}
	return finalThresholdConfigQuery, nil
}
