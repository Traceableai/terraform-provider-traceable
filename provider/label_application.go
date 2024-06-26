package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLabelApplicationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLabelApplicationRuleCreate,
		Read:   resourceLabelApplicationRuleRead,
		Update: resourceLabelApplicationRuleUpdate,
		Delete: resourceLabelApplicationRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Label Application Rule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Label Application Rule",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
			"condition_list": {
				Type:        schema.TypeSet,
				Description: "List of conditions for the rule",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Description: "The key for the condition",
							Required:    true,
						},
						"operator": {
							Type:        schema.TypeString,
							Description: "The operator for the condition",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value for the condition (if applicable)",
							Optional:    true,
						},
						"values": {
							Type:        schema.TypeList,
							Description: "The values for the condition (if applicable)",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				DiffSuppressFunc: suppressConditionListDiff,
			},
			"action": {
				Type:        schema.TypeList,
				Description: "Action to apply for the rule",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "The type of action (DYNAMIC_LABEL_KEY or STATIC_LABELS)",
							Required:    true,
						},
						"entity_types": {
							Type:        schema.TypeList,
							Description: "List of entity types",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"operation": {
							Type:        schema.TypeString,
							Description: "The operation to perform",
							Required:    true,
						},
						"dynamic_label_key": {
							Type:        schema.TypeString,
							Description: "The dynamic label key (if applicable)",
							Optional:    true,
						},
						"static_labels": {
							Type:        schema.TypeList,
							Description: "List of static labels (if applicable)",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func suppressConditionListDiff(k, old, new string, d *schema.ResourceData) bool {
	return suppressListDiff(old, new)
}

func suppressListDiff(old, new string) bool {
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

func resourceLabelApplicationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	conditionList := d.Get("condition_list").(*schema.Set).List()
	action := d.Get("action").([]interface{})[0].(map[string]interface{})

	conditionListStr, err := buildConditionList(conditionList)
	if err != nil {
		return err
	}

	actionStr, err := buildAction(action)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`mutation {
		createLabelApplicationRule(
			labelApplicationRuleData: {
				name: "%s"
				description: "%s"
				enabled: %t
				conditionList: %s
				%s
			}
		) {
			id
		}
	}`, name, description, enabled, conditionListStr, actionStr)

	log.Println(query)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while parsing GraphQL response: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["createLabelApplicationRule"] != nil {
		id := response["data"].(map[string]interface{})["createLabelApplicationRule"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf(responseStr)
	}

	return nil
}

func buildConditionList(conditions []interface{}) (string, error) {
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

func buildAction(action map[string]interface{}) (string, error) {
	actionType := action["type"].(string)
	entityTypes := action["entity_types"].([]interface{})
	operation := action["operation"].(string)
	dynamicLabelKey := ""
	staticLabels := ""

	if actionType == "DYNAMIC_LABEL_KEY" {
		dynamicLabelKey = fmt.Sprintf(`dynamicLabelKey: "%s"`, action["dynamic_label_key"].(string))
	} else if actionType == "STATIC_LABELS" {
		staticLabelsList := action["static_labels"].([]interface{})
		labels, _ := json.Marshal(staticLabelsList)
		staticLabels = fmt.Sprintf(`staticLabels: { ids: %s }`, string(labels))
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

func resourceLabelApplicationRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := fmt.Sprintf(`{
		labelApplicationRules {
			results {
				id
				labelApplicationRuleData {
					name
					description
					enabled
					conditionList {
						leafCondition {
							keyCondition {
								operator
								value
							}
							valueCondition {
								valueConditionType
								stringCondition {
									value
									operator
									stringConditionValueType
									values
								}
								unaryCondition {
									operator
								}
							}
						}
					}
					action {
						entityTypes
						staticLabels {
							ids
						}
						operation
						type
						dynamicLabelKey
					}
				}
			}
		}
	}`)

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while parsing GraphQL response: %s", err)
	}

	if response["data"] != nil {
		results := response["data"].(map[string]interface{})["labelApplicationRules"].(map[string]interface{})["results"].([]interface{})
		for _, result := range results {
			resultMap := result.(map[string]interface{})
			if resultMap["id"].(string) == id {
				labelData := resultMap["labelApplicationRuleData"].(map[string]interface{})
				d.Set("name", labelData["name"].(string))
				d.Set("description", labelData["description"].(string))
				d.Set("enabled", labelData["enabled"].(bool))

				conditionList := parseConditions(labelData["conditionList"].([]interface{}))
				d.Set("condition_list", conditionList)

				action := parseAction(labelData["action"].(map[string]interface{}))
				d.Set("action", action)
				break
			}
		}
	} else {
		d.SetId("")
		return nil
	}

	return nil
}

func parseConditions(conditions []interface{}) []interface{} {
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

func parseAction(action map[string]interface{}) []interface{} {
	parsedAction := map[string]interface{}{
		"type":         action["type"].(string),
		"entity_types": convertToStringSlice2(action["entityTypes"].([]interface{})),
		"operation":    action["operation"].(string),
	}

	if action["type"] == "DYNAMIC_LABEL_KEY" {
		parsedAction["dynamic_label_key"] = action["dynamicLabelKey"].(string)
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

func resourceLabelApplicationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	conditionList := d.Get("condition_list").(*schema.Set).List()
	action := d.Get("action").([]interface{})[0].(map[string]interface{})

	conditionListStr, err := buildConditionList(conditionList)
	if err != nil {
		return err
	}

	actionStr, err := buildAction(action)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`mutation {
		updateLabelApplicationRule(
			labelApplicationRule: {
				id: "%s"
				labelApplicationRuleData: {
					name: "%s"
					description: "%s"
					enabled: %t
					conditionList: %s
					%s
				}
			}
		) {
			id
		}
	}`, id, name, description, enabled, conditionListStr, actionStr)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while parsing GraphQL response: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["updateLabelApplicationRule"] != nil {
		updatedId := response["data"].(map[string]interface{})["updateLabelApplicationRule"].(map[string]interface{})["id"].(string)
		d.SetId(updatedId)
	} else {
		return fmt.Errorf("could not update Label Application Rule, no data returned")
	}

	return nil
}

func resourceLabelApplicationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := fmt.Sprintf(`mutation {
		deleteLabelApplicationRule(id: "%s")
	}`, id)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while parsing GraphQL response: %s", err)
	}

	if response["data"] == nil {
		return fmt.Errorf("could not delete Label Application Rule, no data returned")
	}

	d.SetId("")
	return nil
}
