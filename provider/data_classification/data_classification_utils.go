package data_classification

import (
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/provider/common"

)

func ReturnScopedPatternQuery(scopedPatterns []interface{}) (string) {
	scopedPatternQuery:=""
	for _,scopedPattern := range scopedPatterns {
		scopedPatternData := scopedPattern.(map[string]interface{})
		environments := scopedPatternData["environments"].([]interface{})
		envScopedQueryString := ""
		if len(environments)>0{
			envScopedQueryString = fmt.Sprintf(ENV_SCOPED_QUERY,common.InterfaceToStringSlice(environments))
		}
		scopedPatternName := scopedPatternData["scoped_pattern_name"].(string)
		matchType := scopedPatternData["match_type"].(string)
		urlMatchPatterns := scopedPatternData["url_match_patterns"].([]interface{})
		locations := scopedPatternData["locations"].([]interface{})
		keyPatterns := scopedPatternData["key_patterns"].([]interface{})
		keyOperator := keyPatterns[0].(map[string]interface{})["operator"]
		keyValue := keyPatterns[0].(map[string]interface{})["value"]
		valuePatternString := ""
		if valuePatterns,ok := scopedPatternData["value_patterns"].([]interface{}); ok{
			valueOperator := valuePatterns[0].(map[string]interface{})["operator"]
			valueValue := valuePatterns[0].(map[string]interface{})["value"]
			valuePatternString = fmt.Sprintf(VALUE_PATTERN_QUERY,valueOperator,valueValue)
		}
		patternType := "KEY_VALUE_PATTERN"
		if valuePatternString == ""{
			patternType = "KEY_PATTERN"
		}
		scopedPatternQuery += fmt.Sprintf(SCOPED_PATTERN_QUERY,scopedPatternName,common.InterfaceToStringSlice(urlMatchPatterns),patternType,common.InterfaceToEnumStringSlice(locations),matchType,envScopedQueryString,keyOperator,keyValue,valuePatternString)
	}
	return scopedPatternQuery
}