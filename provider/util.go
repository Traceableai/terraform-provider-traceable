package provider

import (
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"strings"
)


func CheckEmptySegments(segments []interface{}) bool {
	for _,seg := range segments{
		if seg == "" {
			return true
		}
	}
	return false
}
func convertToStringSlicetype(input interface{}) []string {
	var output []string
	switch v := input.(type) {
	case string:
		output = append(output, v)
	case []interface{}:
		for _, elem := range v {
			if str, ok := elem.(string); ok && str != "" {
				output = append(output, str)
			}
		}
	}
	return output
}

func convertToInterfaceSlice(input []string) []interface{} {
	var output []interface{}
	for _, elem := range input {
		output = append(output, elem)
	}
	return output
}

func GetRuleDetailsFromRulesListUsingIdName(response map[string]interface{}, arrayJsonKey string, args ...string) map[string]interface{} {
	return common.CallGetRuleDetailsFromRulesListUsingIdName(response, arrayJsonKey, args...)

}

func jsonifyList(list []interface{}) string {
	var strList []string
	for _, item := range list {
		strList = append(strList, fmt.Sprintf(`"%s"`, item))
	}
	return "[" + strings.Join(strList, ", ") + "]"
}
