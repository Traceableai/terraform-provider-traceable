package dlp

import (
	"fmt"

	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
)

func ReturnFinalThresholdConfigQueryDlp(thresholdConfigs []interface{}) (string, error) {
	thresholdConfigsData := thresholdConfigs[0].(map[string]interface{})
	rollingWindowThresholdConfigs := thresholdConfigsData["rolling_window_threshold_config"].([]interface{})
	dynamicThresholdConfigs := thresholdConfigsData["dynamic_threshold_config"].([]interface{})
	valueBasedThresholdConfigs := thresholdConfigsData["value_based_threshold_config"].([]interface{})
	var finalRollingWindowThresholdConfigString, finalDynamicThresholdConfigString, finalValueBasedThresholdConfigString string
	var finalThresholdConfigQuery string
	for _, rollingWindowThresholdConfig := range rollingWindowThresholdConfigs {
		rollingWindowThresholdConfigData := rollingWindowThresholdConfig.(map[string]interface{})
		userAggregateType := rollingWindowThresholdConfigData["user_aggregate_type"]
		apiAggregateType := rollingWindowThresholdConfigData["api_aggregate_type"]
		countAllowed := rollingWindowThresholdConfigData["count_allowed"]
		duration := rollingWindowThresholdConfigData["duration"]
		finalRollingWindowThresholdConfigString += fmt.Sprintf(rate_limiting.ROLLING_WINDOW_THRESHOLD_CONFIG_QUERY, apiAggregateType, userAggregateType, countAllowed, duration)
	}
	for _, dynamicThresholdConfig := range dynamicThresholdConfigs {
		dynamicThresholdConfigData := dynamicThresholdConfig.(map[string]interface{})
		userAggregateType := dynamicThresholdConfigData["user_aggregate_type"]
		apiAggregateType := dynamicThresholdConfigData["api_aggregate_type"]
		percentageExceedingMeanAllowed := dynamicThresholdConfigData["percentage_exceeding_mean_allowed"]
		meanCalculationDuration := dynamicThresholdConfigData["mean_calculation_duration"]
		duration := dynamicThresholdConfigData["duration"]
		finalDynamicThresholdConfigString += fmt.Sprintf(rate_limiting.DYNAMIC_THRESHOLD_CONFIG_QUERY, apiAggregateType, userAggregateType, percentageExceedingMeanAllowed, meanCalculationDuration, duration)
	}
	for _, valueBasedThresholdConfig := range valueBasedThresholdConfigs {
		valueBasedThresholdConfigData := valueBasedThresholdConfig.(map[string]interface{})
		userAggregateType := valueBasedThresholdConfigData["user_aggregate_type"]
		apiAggregateType := valueBasedThresholdConfigData["api_aggregate_type"]
		uniqueValuesAllowed := valueBasedThresholdConfigData["unique_values_allowed"]
		duration := valueBasedThresholdConfigData["duration"]
		sensitiveParamsEvaluationType := valueBasedThresholdConfigData["sensitive_params_evaluation_type"]
		finalValueBasedThresholdConfigString += fmt.Sprintf(DLP_VALUE_BASED_THRESHOLD_CONFIG_QUERY, apiAggregateType, userAggregateType, uniqueValuesAllowed, duration, sensitiveParamsEvaluationType)
	}
	finalThresholdConfigQuery = finalDynamicThresholdConfigString + finalRollingWindowThresholdConfigString + finalValueBasedThresholdConfigString
	if finalThresholdConfigQuery == "" {
		return "", fmt.Errorf("required at least one threshold config")
	}
	return finalThresholdConfigQuery, nil
}

func GetConditionsStringDlp(targetScope []interface{}, ipAddress []interface{}, regions []interface{}, ipLocationTypes []interface{}, dataTypesConditions []interface{}, requestPayloadSingleValuedConditions []interface{}, requestPayloadMultiValuedConditions []interface{}) string {
	finalConditionsQuery := ""
	if len(ipAddress) > 0 {

		finalConditionsQuery += fmt.Sprintf(rate_limiting.RAW_INPUT_IP_ADDRESS_QUERY, common.InterfaceToStringSlice(ipAddress), false)
	}
	if len(regions) > 0 {
		regionIdentifiers := ""
		for _, region := range regions {
			regionIdentifiers += fmt.Sprintf(`{ countryIsoCode: "%s" },`, region.(string))
		}
		regionIdentifiers = regionIdentifiers[:len(regionIdentifiers)-1]
		finalConditionsQuery += fmt.Sprintf(rate_limiting.REGION_QUERY, false, regionIdentifiers)
	}
	if len(ipLocationTypes) > 0 {

		finalConditionsQuery += fmt.Sprintf(rate_limiting.IP_LOCATION_TYPE_QUERY, common.InterfaceToEnumStringSlice(ipLocationTypes), false)
	}
	if len(dataTypesConditions) > 0 {

		dataTypesConditionsData := dataTypesConditions[0].(map[string]interface{})
		dataTypes := dataTypesConditionsData["data_types"].([]interface{})
		dataTypeIds := dataTypes[0].(map[string]interface{})["data_type_ids"].([]interface{})
		dataSetIds := dataTypes[0].(map[string]interface{})["data_sets_ids"].([]interface{})
		customLocationDataTypeKeyValuePairMatching := dataTypesConditionsData["custom_location_data_type_key_value_pair_matching"].(bool)
		customLocationAttribute := dataTypesConditionsData["custom_location_attribute"].(string)
		customLocationAttributeKeyOp := dataTypesConditionsData["custom_location_attribute_key_operator"].(string)
		customLocationAttributeKeyVal := dataTypesConditionsData["custom_location_attribute_value"].(string)
		if customLocationDataTypeKeyValuePairMatching {
			keyValQuery := ""
			if customLocationAttribute != "REQUEST_BODY" {
				keyValQuery = fmt.Sprintf(DATATYPE_KEY_CONDITIONS, customLocationAttributeKeyOp, customLocationAttributeKeyVal)
			}
			finalConditionsQuery += fmt.Sprintf(DLP_REQUEST_BASED_DATATYPE_CONDITIONS, common.InterfaceToStringSlice(dataSetIds), common.InterfaceToStringSlice(dataTypeIds), customLocationAttribute, keyValQuery)
		} else {
			dataLocationString := fmt.Sprintf(rate_limiting.DATA_LOCATION_STRING, "REQUEST")
			finalConditionsQuery += fmt.Sprintf(rate_limiting.DATA_TYPES_CONDITIONS_QUERY, common.InterfaceToStringSlice(dataSetIds), common.InterfaceToStringSlice(dataTypeIds), dataLocationString)
		}
	}
	if len(requestPayloadSingleValuedConditions) > 0 {

		for _, requestPayloadSingleValuedCondition := range requestPayloadSingleValuedConditions {
			requestPayloadSingleValuedConditionData := requestPayloadSingleValuedCondition.(map[string]interface{})
			requestLocation := requestPayloadSingleValuedConditionData["request_location"].(string)
			keyOp := requestPayloadSingleValuedConditionData["operator"].(string)
			keyValue := requestPayloadSingleValuedConditionData["value"].(string)
			finalConditionsQuery += fmt.Sprintf(rate_limiting.REQ_RES_CONDITIONS_QUERY, requestLocation, keyOp, keyValue)
		}
	}
	if len(requestPayloadMultiValuedConditions) > 0 {

		for _, requestPayloadMultiValuedCondition := range requestPayloadMultiValuedConditions {
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
			finalConditionsQuery += fmt.Sprintf(MULTI_VALUES_REQ_CONDITIONS, requestLocation, keyOp, keyValue, valQuery)
		}
	}
	if len(targetScope) > 0 {

		targetScopeVal := targetScope[0].(map[string]interface{})
		serviceIds := targetScopeVal["service_ids"].([]interface{})
		urlRegex := targetScopeVal["url_regex"].([]interface{})
		if len(serviceIds) > 0 {
			finalConditionsQuery += fmt.Sprintf(SERVICE_SCOPE, common.InterfaceToStringSlice(serviceIds))
		}
		finalConditionsQuery += fmt.Sprintf(URL_REGEX_SCOPE, common.InterfaceToStringSlice(urlRegex))
	}

	return finalConditionsQuery
}

func GetTransactionActionConfigs(ruleType string, severity string, expiry string) string {
	finalTransactionConfigQuery := ""
	if ruleType == "ALLOW" && expiry == "" {
		finalTransactionConfigQuery = fmt.Sprintf(ALLOW_INDEFINITE_TRANSACTION_CONFIGS, ruleType)
	} else if ruleType == "ALLOW" && expiry != "" {
		finalTransactionConfigQuery = fmt.Sprintf(ALLOW_WITH_EXPIRY, ruleType, expiry)
	} else if ruleType == "BLOCK" && expiry == "" {
		finalTransactionConfigQuery = fmt.Sprintf(BLOCK_INDEFINITE_TRANSACTION_CONFIGS, ruleType, severity)
	} else if ruleType == "BLOCK" && expiry != "" {
		finalTransactionConfigQuery = fmt.Sprintf(BLOCK_WITH_EXPIRY, ruleType, severity, expiry)
	} else if ruleType == "ALERT" {
		finalTransactionConfigQuery = fmt.Sprintf(ALERT_TRANSACTION_CONFIGS, ruleType, severity)
	}
	return finalTransactionConfigQuery
}
