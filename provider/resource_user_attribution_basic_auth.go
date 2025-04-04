package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUserAttributionBasicAuthRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleBasicAuthCreate,
		Read:   resourceUserAttributionRuleBasicAuthRead,
		Update: resourceUserAttributionRuleBasicAuthUpdate,
		Delete: resourceUserAttributionRuleBasicAuthDelete,

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
			"environment": {
				Type:        schema.TypeString,
				Description: "environment",
				Optional:    true,
			},
			"url_regex": {
				Type:        schema.TypeString,
				Description: "url regex",
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
				Default:     "BASIC_AUTH",
			},
		},
	}
}

func resourceUserAttributionRuleBasicAuthCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	category := d.Get("category").(string)

	scopeType := d.Get("scope_type").(string)

	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {name: "%s", type: %s, scopeType: %s}
			) {
			  results {
				id
				scopeType
				rank
				name
			  }
			  total
			}
		  }`, name, category, scopeType)
	} else if scopeType == "CUSTOM" {
		environment := d.Get("environment").(string)
		url_regex := d.Get("url_regex").(string)
		var scopedQuery string
		if environment != "" && url_regex == "" {
			scopedQuery = fmt.Sprintf(`customScope: {environmentScopes: [{environmentName: "%s"}]}`, environment)
		} else if url_regex != "" && environment == "" {
			scopedQuery = fmt.Sprintf(`customScope: {urlScopes: [{urlMatchRegex: "%s"}]}`, url_regex)
		}
		if scopedQuery == "" {
			return fmt.Errorf("provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {name: "%s", type: %s, scopeType: CUSTOM, %s}
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
		  }`, name, category, scopedQuery)
	} else {
		return fmt.Errorf("expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "createUserAttributionRule", name)
	log.Println(ruleDetails)
	id := ruleDetails["id"].(string)
	d.SetId(id)
	return nil
}

func resourceUserAttributionRuleBasicAuthRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)
	readQuery := "{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName}urlScopes{urlMatchRegex}}}}}"
	responseStr, err := ExecuteQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "userAttributionRules", id)
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	log.Printf("fetching from read %s", ruleDetails)
	name := ruleDetails["name"].(string)
	category := ruleDetails["type"].(string)
	disabled := ruleDetails["disabled"].(bool)
	d.Set("disabled", disabled)
	d.Set("category", category)
	scopeType := ruleDetails["scopeType"]
	d.Set("name", name)
	if scopeType == "SYSTEM_WIDE" {
		d.Set("scope_type", "SYSTEM_WIDE")
	} else {
		envScope := ruleDetails["customScope"].(map[string]interface{})["environmentScopes"]
		urlScope := ruleDetails["customScope"].(map[string]interface{})["urlScopes"]
		if len(envScope.([]interface{})) == 0 {
			d.Set("scope_type", "CUSTOM")
			d.Set("url_regex", urlScope.([]interface{})[0].(map[string]interface{})["urlMatchRegex"])

		} else {
			d.Set("scope_type", "CUSTOM")
			d.Set("environment", envScope.([]interface{})[0].(map[string]interface{})["environmentName"])
		}
	}
	return nil
}

func resourceUserAttributionRuleBasicAuthUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	scopeType := d.Get("scope_type").(string)
	category := d.Get("category").(string)
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
	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			updateUserAttributionRule(
			  rule: {name: "%s", type: %s,id:"%s",rank:%d,disabled: %t, scopeType: %s}
			) {
				id
				scopeType
				rank
				name
			}
		  }`, name, category, id, rank, disabled, scopeType)
	} else if scopeType == "CUSTOM" {
		environment := d.Get("environment").(string)
		url_regex := d.Get("url_regex").(string)
		var scopedQuery string
		if environment != "" && url_regex == "" {
			scopedQuery = fmt.Sprintf(`customScope: {environmentScopes: [{environmentName: "%s"}]}`, environment)
		} else if url_regex != "" && environment == "" {
			scopedQuery = fmt.Sprintf(`customScope: {urlScopes: [{urlMatchRegex: "%s"}]}`, url_regex)
		}
		if scopedQuery == "" {
			return fmt.Errorf("provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			updateUserAttributionRule(
			  rule: {name: "%s", type: %s,id:"%s",rank:%d,disabled: %t scopeType: CUSTOM, %s}
			) {
				id
				scopeType
				rank
				name
				type
			}
		  }`, name, category, id, rank, disabled, scopedQuery)
	} else {
		return fmt.Errorf("expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	rules := response["data"].(map[string]interface{})["updateUserAttributionRule"].(map[string]interface{})
	// log.Printf(ruleDetails)
	updatedId := rules["id"].(string)
	d.SetId(updatedId)
	return nil
}

func resourceUserAttributionRuleBasicAuthDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := ExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
