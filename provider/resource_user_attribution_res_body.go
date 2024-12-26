package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributionResponseBodyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleResponseBodyCreate,
		Read:   resourceUserAttributionRuleResponseBodyRead,
		Update: resourceUserAttributionRuleResponseBodyUpdate,
		Delete: resourceUserAttributionRuleResponseBodyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the user attribution rule",
				Required:    true,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Description: "auth type of the user attribution rule",
				Optional:    true,
			},
			"url_regex": {
				Type:        schema.TypeString,
				Description: "url regex",
				Required:    true,
			},
			"user_id_location_json_path": {
				Type:        schema.TypeString,
				Description: "user id location json path",
				Required:    true,
			},
			"user_role_location_json_path": {
				Type:        schema.TypeString,
				Description: "user role location json path",
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
				Default:     "RESPONSE_BODY",
			},
		},
	}
}

func resourceUserAttributionRuleResponseBodyCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	url_regex := d.Get("url_regex").(string)
	category := d.Get("category").(string)
	user_role_location_json_path := d.Get("user_role_location_json_path").(string)
	user_id_location_json_path := d.Get("user_id_location_json_path").(string)
	auth_type := d.Get("auth_type").(string)

	var authTypeQuery string
	if auth_type != "" {
		authTypeQuery = fmt.Sprintf(`authentication: { type: "%s" }`, auth_type)
	} else {
		authTypeQuery = ""
	}

	roleLocationString := ""
	if user_role_location_json_path != "" {
		roleLocationString = fmt.Sprintf(`roleLocation: { type: JSON_PATH, jsonPath: "%s" }`, user_role_location_json_path)
	}
	var query string
	query = fmt.Sprintf(`mutation {
		createUserAttributionRule(
		  input: {
		  name: "%s", 
		  type: %s, 
		  scopeType: CUSTOM, 
		  responseBody: {
			%s
			userIdLocation: { type: JSON_PATH, jsonPath: "%s" }
			%s
		  },
		  customScope: { urlScopes: [{ urlMatchRegex: "%s" }] }
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
	  }`, name, category, authTypeQuery, user_id_location_json_path, roleLocationString, url_regex)

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

func resourceUserAttributionRuleResponseBodyRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)
	readQuery := "{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName __typename} urlScopes{urlMatchRegex __typename} __typename} responseBody{authentication{type __typename} userIdLocation{type jsonPath __typename} roleLocation{type jsonPath __typename} condition{type urlMatchRegex __typename} __typename}}}}"
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
	d.Set("name", name)
	category := ruleDetails["type"].(string)
	d.Set("category", category)
	disabled := ruleDetails["disabled"].(bool)
	d.Set("disabled", disabled)
	urlScope := ruleDetails["customScope"].(map[string]interface{})["urlScopes"]
	d.Set("url_regex", urlScope.([]interface{})[0].(map[string]interface{})["urlMatchRegex"])
	responseBodyDetails := ruleDetails["responseBody"]
	if responseBodyDetails != nil {

		auth_type := responseBodyDetails.(map[string]interface{})["authentication"]
		if auth_type != nil {
			auth_type = auth_type.(map[string]interface{})["type"]
			d.Set("auth_type", auth_type)
		} else {
			d.Set("auth_type", nil)
		}

		userIdLocationDetails := responseBodyDetails.(map[string]interface{})["userIdLocation"]
		if userIdLocationDetails != nil {
			user_id_location_json_path := userIdLocationDetails.(map[string]interface{})["jsonPath"]
			d.Set("user_id_location_json_path", user_id_location_json_path)
		}

		userRoleLocationDetails := responseBodyDetails.(map[string]interface{})["roleLocation"]
		if userRoleLocationDetails != nil {
			user_role_location_json_path := userRoleLocationDetails.(map[string]interface{})["jsonPath"]
			d.Set("user_role_location_json_path", user_role_location_json_path)
		}
	}

	return nil
}

func resourceUserAttributionRuleResponseBodyUpdate(d *schema.ResourceData, meta interface{}) error {
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

	name := d.Get("name").(string)
	url_regex := d.Get("url_regex").(string)
	user_role_location_json_path := d.Get("user_role_location_json_path").(string)
	user_id_location_json_path := d.Get("user_id_location_json_path").(string)
	auth_type := d.Get("auth_type").(string)
	disabled := d.Get("disabled").(bool)

	authTypeQuery := ""
	if auth_type != "" {
		authTypeQuery = fmt.Sprintf(`authentication: { type: "%s" }`, auth_type)
	}

	roleLocationString := ""
	if user_role_location_json_path != "" {
		roleLocationString = fmt.Sprintf(`roleLocation: { type: JSON_PATH, jsonPath: "%s" }`, user_role_location_json_path)
	}
	var query string
	query = fmt.Sprintf(`mutation {
		updateUserAttributionRule(
		  rule: {
		  id: "%s",
		  rank: %d,
		  disabled: %t,
		  name: "%s", 
		  type: RESPONSE_BODY, 
		  scopeType: CUSTOM, 
		  responseBody: {
			%s
			userIdLocation: { type: JSON_PATH, jsonPath: "%s" }
			%s
		  },
		  customScope: { urlScopes: [{ urlMatchRegex: "%s" }] }
		}
		) {
			id
			scopeType
			rank
			name
			type
		}
	  }`, id, rank, disabled, name, authTypeQuery, user_id_location_json_path, roleLocationString, url_regex)

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

func resourceUserAttributionRuleResponseBodyDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := ExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
