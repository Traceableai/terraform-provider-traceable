package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributionCustomJsonRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleCustomJsonCreate,
		Read:   resourceUserAttributionRuleCustomJsonRead,
		Update: resourceUserAttributionRuleCustomJsonUpdate,
		Delete: resourceUserAttributionRuleCustomJsonDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the user attribution rule",
				Required:    true,
			},
			"scope_type": {
				Type:        schema.TypeString,
				Description: "system wide, environment, url regex",
				Required:    true,
			},
			"url_regex": {
				Type:        schema.TypeString,
				Description: "url regex",
				Optional:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Description: "environement of rule",
				Optional:    true,
			},
			"auth_type_json": {
				Type:        schema.TypeString,
				Description: "auth type json",
				Required:    true,
			},
			"user_id_json": {
				Type:        schema.TypeString,
				Description: "user id json",
				Optional:    true,
			},
			"user_role_json": {
				Type:        schema.TypeString,
				Description: "user role json",
				Optional:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Optional:    true,
				Default:     false,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "Type of user attribution rule",
				Optional:    true,
				Default:     "CUSTOM_JSON",
			},
		},
	}
}

func resourceUserAttributionRuleCustomJsonCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	category := d.Get("category").(string)
	scope_type := d.Get("scope_type").(string)
	environment := d.Get("environment").(string)
	url_regex := d.Get("url_regex").(string)
	user_role_json := d.Get("user_role_json").(string)
	auth_type_json := d.Get("auth_type_json").(string)
	user_id_json := d.Get("user_id_json").(string)

	if scope_type != "SYSTEM_WIDE" && scope_type != "CUSTOM" {
		return fmt.Errorf("scope_type supported string is SYSTEM_WIDE or CUSTOM")
	}

	customScopeString := ""
	if scope_type == "CUSTOM" {
		if environment != "" && url_regex == "" {
			customScopeString = fmt.Sprintf(`customScope: { environmentScopes: [{ environmentName: "%s" }] }`, environment)
		} else if environment == "" && url_regex != "" {
			customScopeString = fmt.Sprintf(`customScope: { urlScopes: [{ urlMatchRegex: "%s" }] }`, url_regex)
		} else {
			return fmt.Errorf("Required environment or url_regex")
		}
	}

	customJsonString := ""
	if user_id_json != "" && user_role_json != "" && auth_type_json != "" {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s, userIdJson: %s, userRoleJson: %s }`, auth_type_json, user_id_json, user_role_json)
	} else if user_id_json == "" && user_role_json != "" && auth_type_json != "" {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s, userRoleJson: %s, }`, auth_type_json, user_role_json)
	} else if user_id_json != "" && user_role_json == "" && auth_type_json != "" {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s, userIdJson: %s, }`, auth_type_json, user_id_json)
	} else {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s }`, auth_type_json)
	}

	var query string
	query = fmt.Sprintf(`mutation {
		createUserAttributionRule(
		  input: {
		  name: "%s", 
		  type: %s, 
		  scopeType: %s, 
		  %s
		  %s
		}
		) {
		  results {
			id
			scopeType
			rank
			name
			type
		  }
		  total
		}
	  }`, name, category, scope_type, customJsonString, customScopeString)

	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "createUserAttributionRule", name)
	log.Println(ruleDetails)
	id := ruleDetails["id"].(string)
	d.SetId(id)
	return nil
}

func resourceUserAttributionRuleCustomJsonRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)
	readQuery := "{ userAttributionRules { results { id scopeType rank name type disabled customScope { environmentScopes { environmentName __typename } urlScopes { urlMatchRegex __typename } __typename } customJson { authTypeJson userIdJson userRoleJson __typename } } } }"
	responseStr, err := ExecuteQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}
	log.Printf("Response from read %s", responseStr)
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "userAttributionRules", id)
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	log.Printf("fetching from read %s", ruleDetails)
	name := ruleDetails["name"].(string)
	scopeType := ruleDetails["scopeType"].(string)
	category := ruleDetails["type"].(string)
	d.Set("category", category)
	d.Set("name", name)
	disabled := ruleDetails["disabled"].(bool)
	d.Set("disabled", disabled)
	if scopeType == "SYSTEM_WIDE" {
		d.Set("scope_type", "SYSTEM_WIDE")
		// d.Set("url_regex",nil)
		// d.Set("environment",nil)
	} else {
		envScope := ruleDetails["customScope"].(map[string]interface{})["environmentScopes"]
		urlScope := ruleDetails["customScope"].(map[string]interface{})["urlScopes"]
		if len(envScope.([]interface{})) == 0 {
			d.Set("scope_type", "CUSTOM")
			d.Set("url_regex", urlScope.([]interface{})[0].(map[string]interface{})["urlMatchRegex"])
			// d.Set("environment",nil)
		} else {
			d.Set("scope_type", "CUSTOM")
			d.Set("environment", envScope.([]interface{})[0].(map[string]interface{})["environmentName"])
			// d.Set("url_regex",nil)
		}
	}
	customJsonDetails := ruleDetails["customJson"]
	if customJsonDetails != nil {

		authTypeJson, _ := json.Marshal(customJsonDetails.(map[string]interface{})["authTypeJson"])
		userIdJson, _ := json.Marshal(customJsonDetails.(map[string]interface{})["userIdJson"])
		userRoleJson, _ := json.Marshal(customJsonDetails.(map[string]interface{})["userRoleJson"])

		d.Set("auth_type_json", authTypeJson)
		d.Set("user_id_json", userIdJson)
		d.Set("user_role_json", userRoleJson)

		return nil
	}
	return nil
}

func resourceUserAttributionRuleCustomJsonUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	readQuery := "{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName}urlScopes{urlMatchRegex}}}}}"
	readQueryResStr, err := ExecuteQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var readResponse map[string]interface{}
	if err := json.Unmarshal([]byte(readQueryResStr), &readResponse); err != nil {
		return err
	}
	readRuleDetails := GetRuleDetailsFromRulesListUsingIdName(readResponse, "userAttributionRules", id)
	if len(readRuleDetails) == 0 {
		return nil
	}
	rank := int(readRuleDetails["rank"].(float64))
	category := d.Get("category").(string)
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	scope_type := d.Get("scope_type").(string)
	environment := d.Get("environment").(string)
	url_regex := d.Get("url_regex").(string)
	user_role_json := d.Get("user_role_json").(string)
	auth_type_json := d.Get("auth_type_json").(string)
	user_id_json := d.Get("user_id_json").(string)

	if scope_type != "SYSTEM_WIDE" && scope_type != "CUSTOM" {
		return fmt.Errorf("scope_type supported string is SYSTEM_WIDE or CUSTOM")
	}

	customScopeString := ""
	if scope_type == "CUSTOM" {
		if environment != "" && url_regex == "" {
			customScopeString = fmt.Sprintf(`customScope: { environmentScopes: [{ environmentName: "%s" }] }`, environment)
		} else if environment == "" && url_regex != "" {
			customScopeString = fmt.Sprintf(`customScope: { urlScopes: [{ urlMatchRegex: "%s" }] }`, url_regex)
		} else {
			return fmt.Errorf("Required environment or url_regex")
		}
	}

	customJsonString := ""
	if user_id_json != "" && user_role_json != "" && auth_type_json != "" {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s, userIdJson: %s, userRoleJson: %s }`, auth_type_json, user_id_json, user_role_json)
	} else if user_id_json == "" && user_role_json != "" && auth_type_json != "" {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s, userRoleJson: %s, }`, auth_type_json, user_role_json)
	} else if user_id_json != "" && user_role_json == "" && auth_type_json != "" {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s, userIdJson: %s, }`, auth_type_json, user_id_json)
	} else {
		customJsonString = fmt.Sprintf(`customJson: { authTypeJson: %s }`, auth_type_json)
	}

	var query string
	query = fmt.Sprintf(`mutation {
		updateUserAttributionRule(
		  rule: {
			id:"%s",
			rank:%d
			name: "%s", 
			disabled: %t,
			type: %s, 
			scopeType: %s, 
			%s
			%s
		}
		) {
			id
			scopeType
			rank
			name
			type
		}
	  }`, id, rank, name, disabled, category, scope_type, customJsonString, customScopeString)

	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	rules := response["data"].(map[string]interface{})["updateUserAttributionRule"].(map[string]interface{})
	// log.Printf(ruleDetails)
	updatedId := rules["id"].(string)
	d.SetId(updatedId)
	return nil
}

func resourceUserAttributionRuleCustomJsonDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := ExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
