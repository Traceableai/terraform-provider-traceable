package label_management

import (
	"encoding/json"
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"strings"
)

func BuildConditionList(conditions []interface{}) (string, error) {
	var conditionStrings []string

	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		key := conditionMap["key"].(string)
		operator := conditionMap["operator"].(string)
		valueConditionType := "UNARY_CONDITION"
		stringConditionOperator := ""
		stringConditionValuePart := ""
		stringConditionValueType := ""
		valueConditionPart := ""

		if value, ok := conditionMap["value"].(string); ok && value != "" {
			valueConditionType = "STRING_CONDITION"
			stringConditionOperator = operator
			stringConditionValueType = "VALUE"
			stringConditionValuePart = fmt.Sprintf(`value: "%s"`, value)
		} else if values, ok := conditionMap["values"].([]interface{}); ok && len(values) > 0 {
			valueConditionType = "STRING_CONDITION"
			stringConditionOperator = operator
			stringConditionValueType = "VALUES"
			valueList, _ := json.Marshal(values)
			stringConditionValuePart = fmt.Sprintf(`values: %s`, string(valueList))
		}

		if valueConditionType == "STRING_CONDITION" {
			valueConditionPart = fmt.Sprintf(`
						stringCondition: {
							operator: %s
							stringConditionValueType: %s
							%s
						}`, stringConditionOperator, stringConditionValueType, stringConditionValuePart)
		} else if valueConditionType == "UNARY_CONDITION" {
			valueConditionPart = `unaryCondition: { operator: OPERATOR_EXISTS }`
		}

		conditionString := fmt.Sprintf(`
			{
				leafCondition: {
					keyCondition: { operator: OPERATOR_EQUALS, value: "%s" }
					valueCondition: {
						valueConditionType: %s
						%s
					}
				}
			}`, key, valueConditionType, valueConditionPart)

		conditionStrings = append(conditionStrings, conditionString)
	}

	return fmt.Sprintf("[%s]", strings.Join(conditionStrings, ",")), nil
}

func SuppressListDiff(old, new string) bool {
	oldList := strings.FieldsFunc(old, func(r rune) bool { return r == ',' || r == '[' || r == ']' || r == '{' || r == '}' || r == '"' })
	newList := strings.FieldsFunc(new, func(r rune) bool { return r == ',' || r == '[' || r == ']' || r == '{' || r == '}' || r == '"' })
	if len(oldList) != len(newList) {
		return false
	}
	oldMap := make(map[string]bool)
	for _, v := range oldList {
		oldMap[v] = true
	}
	for _, v := range newList {
		if !oldMap[v] {
			return false
		}
	}
	return true
}

func BuildAction(action map[string]interface{}) (string, error) {
	actionType := action["type"].(string)
	entityTypes := action["entity_types"].([]interface{})
	operation := action["operation"].(string)
	dynamicLabelKey := ""
	staticLabels := ""

	if actionType == "DYNAMIC_LABEL" {
		dynamicLabels := action["dynamic_labels"].([]interface{})[0].(map[string]interface{})
		attribute := dynamicLabels["attribute"]
		if regex, exist := dynamicLabels["regex"].(string); exist && regex != "" {
			dynamicLabelKey = fmt.Sprintf(`dynamicLabel: {
				expression: "${%s}"
				tokenExtractionRules: [{ key: "%s", regexCapture: "%s" }]
				}`, attribute, attribute, regex)
		} else {
			actionType = "DYNAMIC_LABEL_KEY"
			dynamicLabelKey = fmt.Sprintf(`dynamicLabelKey: "%s"`, attribute)
		}
	} else if actionType == "STATIC_LABELS" {
		staticLabelsList := action["static_labels"].([]interface{})
		labels := common.InterfaceToStringSlice(staticLabelsList)
		staticLabels = fmt.Sprintf(`staticLabels: { ids: %s }`, labels)
	}

	entityTypesList, _ := json.Marshal(entityTypes)

	actionStr := fmt.Sprintf(`action: {
		type: %s
		entityTypes: %s
		operation: %s
		%s
		%s
	}`, actionType, string(entityTypesList), operation, dynamicLabelKey, staticLabels)

	return actionStr, nil
}

func ParseConditions(conditions []interface{}) []interface{} {
	var parsedConditions []interface{}

	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})["leafCondition"].(map[string]interface{})
		keyCondition := conditionMap["keyCondition"].(map[string]interface{})
		valueCondition := conditionMap["valueCondition"].(map[string]interface{})

		parsedCondition := map[string]interface{}{
			"key":      keyCondition["value"].(string),
			"operator": keyCondition["operator"].(string),
		}

		if valueConditionType, ok := valueCondition["valueConditionType"].(string); ok {
			switch valueConditionType {
			case "STRING_CONDITION":
				stringCondition := valueCondition["stringCondition"].(map[string]interface{})
				if stringCondition["stringConditionValueType"] == "VALUE" {
					parsedCondition["value"] = stringCondition["value"].(string)
				} else if values, exists := stringCondition["values"].([]interface{}); exists && len(values) > 0 {
					parsedCondition["values"] = convertToStringSlice2(values)
				}
			case "UNARY_CONDITION":
				unaryCondition := valueCondition["unaryCondition"].(map[string]interface{})
				parsedCondition["operator"] = unaryCondition["operator"].(string)
			}
		}

		parsedConditions = append(parsedConditions, parsedCondition)
	}

	return parsedConditions
}

func ParseAction(action map[string]interface{}) []interface{} {
	parsedAction := map[string]interface{}{
		"type":         action["type"].(string),
		"entity_types": convertToStringSlice2(action["entityTypes"].([]interface{})),
		"operation":    action["operation"].(string),
	}

	if action["type"] == "DYNAMIC_LABEL" {
		if dynamicLabelKey, exist := action["dynamicLabelKey"]; exist {
			parsedActionDynamicLabels := map[string]interface{}{
				"attribute": dynamicLabelKey,
				"regex":     "",
			}
			parsedAction["dynamic_labels"] = parsedActionDynamicLabels
		} else if dynamicLabels, exist := action["dynamicLabel"]; exist {
			tokenExtractionRulesMap := dynamicLabels.(map[string]interface{})
			key := tokenExtractionRulesMap["key"]
			regex := tokenExtractionRulesMap["regex"]
			parsedActionDynamicLabels := map[string]interface{}{
				"attribute": key,
				"regex":     regex,
			}
			parsedAction["dynamic_labels"] = parsedActionDynamicLabels
		}
	} else if action["type"] == "STATIC_LABELS" {
		staticLabels := action["staticLabels"].(map[string]interface{})["ids"].([]interface{})
		if len(staticLabels) > 0 {
			parsedAction["static_labels"] = convertToStringSlice2(staticLabels)
		}
	}

	return []interface{}{parsedAction}
}

func convertToStringSlice2(input []interface{}) []string {
	var output []string
	for _, item := range input {
		if str, ok := item.(string); ok {
			output = append(output, str)
		}
	}
	return output
}
