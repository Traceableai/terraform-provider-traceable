package provider

import (
	"fmt"
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

// function to convert a list of strings to a GraphQL-compatible string list
func jsonifyList(list []interface{}) string {
	var strList []string
	for _, item := range list {
		strList = append(strList, fmt.Sprintf(`"%s"`, item))
	}
	return "[" + strings.Join(strList, ", ") + "]"
}

func parseConditions(conditions []interface{}) string {
	var conditionsStrs []string
	for _, c := range conditions {
		cond := c.(map[string]interface{})
		keyCondition := cond["key_condition"].([]interface{})[0].(map[string]interface{})
		valueCondition := cond["value_condition"].([]interface{})[0].(map[string]interface{})

		keyCondStr := fmt.Sprintf("{operator: %s, value: \"%s\"}",
			keyCondition["operator"].(string),
			keyCondition["value"].(string))

		valCondStr := jsonifyValueCondition(valueCondition)

		conditionsStrs = append(conditionsStrs, fmt.Sprintf("{leafCondition: {keyCondition: %s, valueCondition: %s}}",
			keyCondStr, valCondStr))
	}
	return strings.Join(conditionsStrs, ", ")
}

func jsonifyValueCondition(cond map[string]interface{}) string {
	valueCondType := cond["value_condition_type"].(string)
	switch valueCondType {
	case "STRING_CONDITION":
		stringCond := cond["string_condition"].(map[string]interface{})
		values := "null"
		if v, ok := stringCond["values"]; ok {
			valueList := v.([]interface{})
			valueStrs := make([]string, len(valueList))
			for i, val := range valueList {
				valueStrs[i] = fmt.Sprintf("\"%s\"", val.(string))
			}
			values = "[" + strings.Join(valueStrs, ", ") + "]"
		}
		return fmt.Sprintf("{valueConditionType: %s, stringCondition: {value: \"%s\", operator: %s, stringConditionValueType: %s, values: %s}}",
			valueCondType,
			stringCond["value"].(string),
			stringCond["operator"].(string),
			stringCond["string_condition_value_type"].(string),
			values)
	case "UNARY_CONDITION":
		unaryCond := cond["unary_condition"].(map[string]interface{})
		return fmt.Sprintf("{valueConditionType: %s, unaryCondition: {operator: %s}}",
			valueCondType,
			unaryCond["operator"].(string))
	default:
		return "{}"
	}
}

func jsonifyAction(action map[string]interface{}) string {
	entityTypes := action["entity_types"].([]interface{})
	entityTypeStrs := make([]string, len(entityTypes))
	for i, et := range entityTypes {
		entityTypeStrs[i] = fmt.Sprintf("\"%s\"", et.(string))
	}

	actionType := action["type"].(string)
	operation := action["operation"].(string)
	entityTypeList := "[" + strings.Join(entityTypeStrs, ", ") + "]"

	actionStr := fmt.Sprintf("{type: \"%s\", entityTypes: %s, operation: \"%s\"", actionType, entityTypeList, operation)

	if actionType == "DYNAMIC_LABEL_KEY" {
		dynamicLabelKey := action["dynamic_label_key"].(string)
		actionStr += fmt.Sprintf(", dynamicLabelKey: \"%s\"}", dynamicLabelKey)
	} else if actionType == "STATIC_LABELS" {
		ids := action["static_labels"].([]interface{})
		idStrs := make([]string, len(ids))
		for i, id := range ids {
			idStrs[i] = fmt.Sprintf("\"%s\"", id.(string))
		}
		idsList := "[" + strings.Join(idStrs, ", ") + "]"
		actionStr += fmt.Sprintf(", staticLabels: {ids: %s}}", idsList)
	} else {
		actionStr += "}"
	}

	return actionStr
}

func jsonifyConditionList(conditions []interface{}) string {
	var conditionStrings []string

	for _, c := range conditions {
		condition := c.(map[string]interface{})
		keyCondition := condition["key_condition"].(map[string]interface{})
		valueCondition := condition["value_condition"].(map[string]interface{})

		// Serialize the key condition
		keyConditionStr := fmt.Sprintf("{operator: %s, value: \"%s\"}",
			keyCondition["operator"].(string), keyCondition["value"].(string))

		// Serialize the value condition
		valueConditionStr := jsonifyValueCondition(valueCondition)

		// Combine key and value conditions into a leafCondition
		leafConditionStr := fmt.Sprintf("{leafCondition: {keyCondition: %s, valueCondition: %s}}",
			keyConditionStr, valueConditionStr)

		conditionStrings = append(conditionStrings, leafConditionStr)
	}

	// Join all conditions into a single string
	return strings.Join(conditionStrings, ", ")
}
