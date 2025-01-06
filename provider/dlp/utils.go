package dlp

import (
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
)

func ReturnFinalThresholdConfigQueryDlp(thresholdConfigs []interface{}) (string,error) {
	thresholdConfigsData := thresholdConfigs[0].(map[string]interface{})
	rollingWindowThresholdConfigs := thresholdConfigsData["rolling_window_threshold_config"].([]interface{})
	dynamicThresholdConfigs := thresholdConfigsData["dynamic_threshold_config"].([]interface{})
	valueBasedThresholdConfigs := thresholdConfigsData["value_based_threshold_config"].([]interface{})
	var finalRollingWindowThresholdConfigString,finalDynamicThresholdConfigString,finalValueBasedThresholdConfigString string
	var finalThresholdConfigQuery string
	for _,rollingWindowThresholdConfig := range rollingWindowThresholdConfigs {
		rollingWindowThresholdConfigData := rollingWindowThresholdConfig.(map[string]interface{})
		userAggregateType := rollingWindowThresholdConfigData["user_aggregate_type"]
		apiAggregateType := rollingWindowThresholdConfigData["api_aggregate_type"]
		countAllowed := rollingWindowThresholdConfigData["count_allowed"]
		duration := rollingWindowThresholdConfigData["duration"]
		finalRollingWindowThresholdConfigString += fmt.Sprintf(rate_limiting.ROLLING_WINDOW_THRESHOLD_CONFIG_QUERY,apiAggregateType,userAggregateType,countAllowed,duration)
	}
	for _,dynamicThresholdConfig := range dynamicThresholdConfigs {
		dynamicThresholdConfigData := dynamicThresholdConfig.(map[string]interface{})
		userAggregateType := dynamicThresholdConfigData["user_aggregate_type"]
		apiAggregateType := dynamicThresholdConfigData["api_aggregate_type"]
		percentageExceedingMeanAllowed := dynamicThresholdConfigData["percentage_exceeding_mean_allowed"]
		meanCalculationDuration := dynamicThresholdConfigData["mean_calculation_duration"]
		duration := dynamicThresholdConfigData["duration"]
		finalDynamicThresholdConfigString += fmt.Sprintf(rate_limiting.DYNAMIC_THRESHOLD_CONFIG_QUERY,apiAggregateType,userAggregateType,percentageExceedingMeanAllowed,meanCalculationDuration,duration)
	}
	for _,valueBasedThresholdConfig := range valueBasedThresholdConfigs {
		valueBasedThresholdConfigData := valueBasedThresholdConfig.(map[string]interface{})
		userAggregateType := valueBasedThresholdConfigData["user_aggregate_type"]
		apiAggregateType := valueBasedThresholdConfigData["api_aggregate_type"]
		uniqueValuesAllowed := valueBasedThresholdConfigData["unique_values_allowed"]
		duration := valueBasedThresholdConfigData["duration"]
		sensitiveParamsEvaluationType := valueBasedThresholdConfigData["sensitive_params_evaluation_type"]
		finalValueBasedThresholdConfigString += fmt.Sprintf(DLP_VALUE_BASED_THRESHOLD_CONFIG_QUERY,apiAggregateType,userAggregateType,uniqueValuesAllowed,duration,sensitiveParamsEvaluationType)
	}
	finalThresholdConfigQuery = finalDynamicThresholdConfigString+finalRollingWindowThresholdConfigString+finalValueBasedThresholdConfigString
	if finalThresholdConfigQuery == ""{
		return "",fmt.Errorf("required at least one threshold config")
	}
	return finalThresholdConfigQuery,nil
}