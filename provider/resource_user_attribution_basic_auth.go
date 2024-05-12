package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributionBasicAuthRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleBasicAuthCreate,
		Read:   resourceUserAttributionRuleBasicAuthRead,
		Update: resourceUserAttributionRuleBasicAuthUpdate,
		Delete: resourceUserAttributionRuleBasicAuthDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "name of the user attribution rule",
				Required:    true,
			},
			"scope_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "system wide, environment, url regex",
				Required:    true,
			},
			"environment": &schema.Schema{
				Type:        schema.TypeString,
				Description: "environment",
				Optional:    true,
			},
			"url_regex": &schema.Schema{
				Type:        schema.TypeString,
				Description: "url regex",
				Optional:    true,
			},
		},
	}
}

func resourceUserAttributionRuleBasicAuthCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	scopeType:=d.Get("scope_type").(string)
	
	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {name: "%s", type: BASIC_AUTH, scopeType: %s}
			) {
			  results {
				id
				scopeType
				rank
				name
			  }
			  total
			}
		  }`,name,scopeType)
	} else if scopeType== "CUSTOM" {
		environment:=d.Get("environment").(string)
		url_regex:=d.Get("url_regex").(string)
		var scopedQuery string
		if environment!="" && url_regex=="" {
			scopedQuery=fmt.Sprintf(`customScope: {environmentScopes: [{environmentName: "%s"}]}`,environment)
		} else if url_regex!="" && environment=="" {
			scopedQuery=fmt.Sprintf(`customScope: {urlScopes: [{urlMatchRegex: "%s"}]}`,url_regex)
		}
		if scopedQuery==""{
			return fmt.Errorf("Provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {name: "%s", type: BASIC_AUTH, scopeType: CUSTOM, %s}
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
		  }`,name,scopedQuery)
	}else{
		return fmt.Errorf("Expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		fmt.Errorf("Error:", err)
	}
	log.Println("This is the graphql query %s", query)
	log.Println("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Errorf("Error:", err)
	}
	ruleDetails := getRuleDetailsFromRulesListUsingIdName(response,"createUserAttributionRule",name)
	log.Println(ruleDetails)
	id:=ruleDetails["id"].(string)
	d.SetId(id)
 	return nil
}

func resourceUserAttributionRuleBasicAuthRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	readQuery:="{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName}urlScopes{urlMatchRegex}}}}}"
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}
	ruleDetails:=getRuleDetailsFromRulesListUsingIdName(response,"userAttributionRules" ,id)
	if len(ruleDetails)==0{
		return nil
	}
	log.Println("fetching from read %s",ruleDetails)
	name:=ruleDetails["name"].(string)
	scopeType:=ruleDetails["scopeType"]
	d.Set("name",name)
	if scopeType=="SYSTEM_WIDE"{
		d.Set("scope_type", scopeType)
	}else{
		envScope := ruleDetails["customScope"].(map[string]interface{})["environmentScopes"]
		urlScope := ruleDetails["customScope"].(map[string]interface{})["urlScopes"]
		if len(envScope.([]interface{}))==0{
			d.Set("scopeType",scopeType)
			d.Set("url_regex",urlScope.([]interface{})[0].(map[string]interface{})["urlScopes"])
			
		}else{
			d.Set("scopeType",scopeType)
			d.Set("environment",envScope.([]interface{})[0].(map[string]interface{})["environmentName"])
		}
	}
	return nil
}

func resourceUserAttributionRuleBasicAuthUpdate(d *schema.ResourceData, meta interface{}) error {
	id:=d.Id()
	name := d.Get("name").(string)
	scopeType:=d.Get("scope_type").(string)
	readQuery:="{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName}urlScopes{urlMatchRegex}}}}}"
	readQueryResStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var readResponse map[string]interface{}
	if err := json.Unmarshal([]byte(readQueryResStr), &readResponse); err != nil {
		return err
	}
	readRuleDetails:=getRuleDetailsFromRulesListUsingIdName(readResponse,"userAttributionRules" ,id)
	if len(readRuleDetails)==0{
		return nil
	}
	rank:=int(readRuleDetails["rank"].(float64))
	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			updateUserAttributionRule(
			  rule: {name: "%s", type: BASIC_AUTH,id:"%s",rank:%d, scopeType: %s}
			) {
				id
				scopeType
				rank
				name
			}
		  }`,name,id,rank,scopeType)
	} else if scopeType== "CUSTOM" {
		environment:=d.Get("environment").(string)
		url_regex:=d.Get("url_regex").(string)
		var scopedQuery string
		if environment!="" && url_regex=="" {
			scopedQuery=fmt.Sprintf(`customScope: {environmentScopes: [{environmentName: "%s"}]}`,environment)
		} else if url_regex!="" && environment=="" {
			scopedQuery=fmt.Sprintf(`customScope: {urlScopes: [{urlMatchRegex: "%s"}]}`,url_regex)
		}
		if scopedQuery==""{
			return fmt.Errorf("Provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			updateUserAttributionRule(
			  rule: {name: "%s", type: BASIC_AUTH,id:"%s",rank:%d, scopeType: CUSTOM, %s}
			) {
				id
				scopeType
				rank
				name
				type
			}
		  }`,name,id,rank,scopedQuery)
	}else{
		return fmt.Errorf("Expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		fmt.Errorf("Error:", err)
	}
	log.Println("This is the graphql query %s", query)
	log.Println("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Errorf("Error:", err)
	}
	rules := response["data"].(map[string]interface{})["updateUserAttributionRule"].(map[string]interface{})
	// log.Println(ruleDetails)
	updatedId:=rules["id"].(string)
	d.SetId(updatedId)
	return nil
}

func resourceUserAttributionRuleBasicAuthDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := executeQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
