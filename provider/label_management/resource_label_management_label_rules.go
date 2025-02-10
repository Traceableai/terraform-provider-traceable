package label_management

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"log"
)

func ResourceLabelApplicationRule() *schema.Resource {
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
	return SuppressListDiff(old, new)
}

func resourceLabelApplicationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	conditionList := d.Get("condition_list").(*schema.Set).List()
	action := d.Get("action").([]interface{})[0].(map[string]interface{})

	conditionListStr, err := BuildConditionList(conditionList)
	if err != nil {
		return err
	}

	actionStr, err := BuildAction(action)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(CREATE_LABEL_RULE_QUERY, name, description, enabled, conditionListStr, actionStr)

	log.Println(query)

	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["createLabelApplicationRule"] != nil {
		id := response["data"].(map[string]interface{})["createLabelApplicationRule"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf("err %s", responseStr)
	}

	return nil
}

func resourceLabelApplicationRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	responseStr, err := common.CallExecuteQuery(READ_LABEL_RULE_QUERY, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
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

				conditionList := ParseConditions(labelData["conditionList"].([]interface{}))
				d.Set("condition_list", conditionList)

				action := ParseAction(labelData["action"].(map[string]interface{}))
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

func resourceLabelApplicationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	conditionList := d.Get("condition_list").(*schema.Set).List()
	action := d.Get("action").([]interface{})[0].(map[string]interface{})

	conditionListStr, err := BuildConditionList(conditionList)
	if err != nil {
		return err
	}

	actionStr, err := BuildAction(action)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(LABEL_RULE_UPDATE_QUERY, id, name, description, enabled, conditionListStr, actionStr)

	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
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
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
	}

	if response["data"] == nil {
		return fmt.Errorf("could not delete Label Application Rule, no data returned")
	}

	d.SetId("")
	return nil
}
