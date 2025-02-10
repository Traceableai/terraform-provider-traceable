package enumeration

import (
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
)

func ReturnFinalThresholdConfigQueryEnumeration(threshold_configs []interface{}) (string, error) {
	finalThresholdConfigQuery := ""
	for _, thresholdConfig := range threshold_configs {
		thresholdConfigData := thresholdConfig.(map[string]interface{})
		apiAggregateType := thresholdConfigData["api_aggregate_type"].(string)
		userAggregateType := thresholdConfigData["user_aggregate_type"].(string)
		uniqueValuesAllowed := thresholdConfigData["unique_values_allowed"].(int)
		duration := thresholdConfigData["duration"].(string)
		thresholdConfigType := thresholdConfigData["threshold_config_type"]
		sensitiveParamEvaluationType := thresholdConfigData["sensitive_param_evaluation_type"].(string)
		if thresholdConfigType == "PATH_PARAMS" || thresholdConfigType == "REQUEST_BODY" {
			finalThresholdConfigQuery += fmt.Sprintf(rate_limiting.ENUMERATION_THRESHOLD_CONFIG_QUERY, apiAggregateType, userAggregateType, thresholdConfigType, uniqueValuesAllowed, duration)
		} else if thresholdConfigType == "SENSITIVE_PARAMS" {
			finalThresholdConfigQuery += fmt.Sprintf(rate_limiting.ENUMERATION_THRESHOLD_CONFIG_SENSITIVE_PARAM_QUERY, apiAggregateType, userAggregateType, thresholdConfigType, uniqueValuesAllowed, duration, sensitiveParamEvaluationType)
		}
	}
	return finalThresholdConfigQuery, nil
}
