package data_classification

import (
	"fmt"

	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func GetDatSetQuery(id string, name string, description string, iconType string) string {
	if id == "" {
		return fmt.Sprintf(CREATE_DATA_SET_QUERY, name, description, iconType)
	}
	return fmt.Sprintf(UPDATE_DATA_SET_QUERY, id, name, description, iconType)
}

func GetOverridesCreateQuery(id string, name string, description string, dataSuppressionOverride string, environments []interface{}, spanFilter []interface{}) string {
	spanFilterQuery := ""
	if len(spanFilter) > 0 {
		keyPatterns := spanFilter[0].(map[string]interface{})["key_patterns"].([]interface{})
		keyOperator := keyPatterns[0].(map[string]interface{})["operator"]
		keyValue := keyPatterns[0].(map[string]interface{})["value"]
		valuePatternString := ""
		if valuePattern, ok := spanFilter[0].(map[string]interface{})["value_patterns"].([]interface{}); ok {
			valueOp := valuePattern[0].(map[string]interface{})["operator"]
			value := valuePattern[0].(map[string]interface{})["value"]
			valuePatternString = fmt.Sprintf(VALUE_PATTERN_QUERY, valueOp, value)
		}
		spanFilterQuery = fmt.Sprintf(SPAN_FILTER_OVERRIDES_QUERY, keyOperator, keyValue, valuePatternString)
	}
	envScopedString := ""
	if len(environments) > 0 {
		envScopedString = fmt.Sprintf(ENV_OVERRIDES_QUERY, common.InterfaceToStringSlice(environments))
	}
	if id == "" {
		return fmt.Sprintf(CREATE_OVERRIDES_QUERY, name, description, envScopedString, spanFilterQuery, dataSuppressionOverride)
	}
	return fmt.Sprintf(UPDATE_OVERRIDES_QUERY, id, name, description, envScopedString, spanFilterQuery, dataSuppressionOverride)
}

func GetScopedPatternQuery(scopedPatterns []interface{}) string {
	scopedPatternQuery := ""
	for _, scopedPattern := range scopedPatterns {
		scopedPatternData := scopedPattern.(map[string]interface{})
		environments := scopedPatternData["environments"].([]interface{})
		envScopedQueryString := ""
		if len(environments) > 0 {
			envScopedQueryString = fmt.Sprintf(ENV_SCOPED_QUERY, common.InterfaceToStringSlice(environments))
		}
		scopedPatternName := scopedPatternData["scoped_pattern_name"].(string)
		matchType := scopedPatternData["match_type"].(string)
		urlMatchPatterns := scopedPatternData["url_match_patterns"].([]interface{})
		locations := scopedPatternData["locations"].([]interface{})
		keyPatterns := scopedPatternData["key_patterns"].([]interface{})
		keyOperator := keyPatterns[0].(map[string]interface{})["operator"]
		keyValue := keyPatterns[0].(map[string]interface{})["value"]
		valuePatternString := ""
		if valuePatterns, ok := scopedPatternData["value_patterns"].([]interface{}); ok {
			valueOperator := valuePatterns[0].(map[string]interface{})["operator"]
			valueValue := valuePatterns[0].(map[string]interface{})["value"]
			valuePatternString = fmt.Sprintf(VALUE_PATTERN_QUERY, valueOperator, valueValue)
		}
		patternType := "KEY_VALUE_PATTERN"
		if valuePatternString == "" {
			patternType = "KEY_PATTERN"
		}
		scopedPatternQuery += fmt.Sprintf(SCOPED_PATTERN_QUERY, scopedPatternName, common.InterfaceToStringSlice(urlMatchPatterns), patternType, common.InterfaceToEnumStringSlice(locations), matchType, envScopedQueryString, keyOperator, keyValue, valuePatternString)
	}
	return scopedPatternQuery
}
