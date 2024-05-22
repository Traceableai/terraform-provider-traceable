package provider

import (
	"fmt"
	"log"
	"strings"
)

func listToString(stringArray []string) string {
	var formattedStrings []string
	for _, s := range stringArray {
		formattedStrings = append(formattedStrings, fmt.Sprintf(`"%s"`, s))
	}
	return strings.Join(formattedStrings, ", ")
}

func toStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		stringSlice[i] = v.(string)
	}
	return stringSlice
}

func convertToStringSlice(data []interface{}) []interface{} {
	var result []interface{}
	for _, v := range data {
		result = append(result, v.(interface{}))
	}
	return result
}

func getRuleDetailsFromRulesListUsingIdName(response map[string]interface{}, arrayJsonKey string, args ...string) map[string]interface{} {
	var res map[string]interface{}
	rules := response["data"].(map[string]interface{})[arrayJsonKey].(map[string]interface{})
	results := rules["results"].([]interface{})
	id_name := args[0]
	if len(args)==1{
		args = append(args, "id")
		args = append(args, "name")
	}
	// log.Println(id_name)
	// log.Println(results)
	for _, rule := range results {
		ruleData := rule.(map[string]interface{})
		log.Println(ruleData)
		rule_id := ruleData[args[1]].(string)
		var rule_name string
		var ok bool
		if rule_name, ok = ruleData[args[2]].(string); ok {
			// fmt.Println("Rule Name:", rule_name)
		} else {
			rule_name = ""
		}
		if rule_id == id_name || rule_name == id_name {
			// log.Println("Inside if block %s",rule)
			return rule.(map[string]interface{})
		}
	}
	return res
}

// function to convert a list of strings to a GraphQL-compatible string list
func jsonifyList(list []interface{}) string {
	var strList []string
	for _, item := range list {
		strList = append(strList, fmt.Sprintf(`"%s"`, item))
	}
	return "[" + strings.Join(strList, ", ") + "]"
}
