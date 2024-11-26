package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributionRequestHeaderRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleRequestHeaderCreate,
		Read:   resourceUserAttributionRuleRequestHeaderRead,
		Update: resourceUserAttributionRuleRequestHeaderUpdate,
		Delete: resourceUserAttributionRuleRequestHeaderDelete,

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
			"user_id_location": {
				Type:        schema.TypeString,
				Description: "user id location",
				Required:    true,
			},
			"user_id_regex_capture_group": {
				Type:        schema.TypeString,
				Description: "user id regex capture group",
				Optional:    true,
			},
			"user_role_location": {
				Type:        schema.TypeString,
				Description: "user role location",
				Optional:    true,
			},
			"role_location_regex_capture_group": {
				Type:        schema.TypeString,
				Description: "user role location regex capture group",
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
				Default:     "REQUEST_HEADER",
			},
		},
	}
}

func resourceUserAttributionRuleRequestHeaderCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	scopeType:=d.Get("scope_type").(string)
	environment:=d.Get("environment").(string)
	url_regex:=d.Get("url_regex").(string)
	category:=d.Get("category").(string)
	user_id_location:=d.Get("user_id_location").(string)
	user_id_regex_capture_group:=d.Get("user_id_regex_capture_group").(string)
	user_role_location:=d.Get("user_role_location").(string)
	role_location_regex_capture_group:=d.Get("role_location_regex_capture_group").(string)
	auth_type:=d.Get("auth_type").(string)

	var scopedQuery string
	
	if environment!="" && url_regex=="" {
		scopedQuery=fmt.Sprintf(`customScope: {environmentScopes: [{environmentName: "%s"}]}`,environment)
	} else if url_regex!="" && environment=="" {
		scopedQuery=fmt.Sprintf(`customScope: {urlScopes: [{urlMatchRegex: "%s"}]}`,url_regex)
	}
	var authTypeQuery string
	if auth_type!=""{
		authTypeQuery=fmt.Sprintf(`authentication: { type: "%s" }`,auth_type)
	}else{
		authTypeQuery=""
	}

	parsingTargetUserId:=""
	if user_id_regex_capture_group!=""{
		parsingTargetUserId=fmt.Sprintf(`parsingTarget: {
										type: REGEX_CAPTURE_GROUP
										regexCaptureGroup: "%s"
									}`,user_id_regex_capture_group)
	}
	if user_role_location=="" && role_location_regex_capture_group!=""{
		return fmt.Errorf("role_location_regex_capture_group is not expected here without user_role_location")
	}
	parsingTargetUserRole:=""
	if user_role_location!=""{
		parsingTargetUserRole=fmt.Sprintf(`roleLocation: {
			type: HEADER
			headerName: "%s"
			parsingTarget: {
				type: REGEX_CAPTURE_GROUP
				regexCaptureGroup: "%s"
			}
		}`,user_role_location,role_location_regex_capture_group)
	}

	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {
				name: "%s",
				type: %s,
				scopeType: %s,
				requestHeader: {
					%s
					userIdLocation: {
						type: HEADER
						headerName: "%s"
						%s
					}
					%s
				}
			}
			) {
			  results {
				id
				scopeType
				rank
				name
			  }
			  total
			}
		  }`,name,category,scopeType,authTypeQuery,user_id_location,parsingTargetUserId,parsingTargetUserRole)
	} else if scopeType== "CUSTOM" {
		
		if scopedQuery==""{
			return fmt.Errorf("Provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {name: "%s", 
			  type: %s, 
			  scopeType: CUSTOM, 
			  requestHeader: {
				%s
				userIdLocation: {
					type: HEADER
					headerName: "%s"
					%s
				}
				%s
			  },
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
		  }`,name,category,authTypeQuery,user_id_location,parsingTargetUserId,parsingTargetUserRole,scopedQuery)
	}else{
		return fmt.Errorf("Expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

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
	ruleDetails := getRuleDetailsFromRulesListUsingIdName(response,"createUserAttributionRule",name)
	log.Println(ruleDetails)
	id:=ruleDetails["id"].(string)
	d.SetId(id)
 	return nil
}

func resourceUserAttributionRuleRequestHeaderRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)
	readQuery:="{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName __typename}urlScopes{urlMatchRegex __typename}__typename}requestHeader{authentication{type __typename}userIdLocation{type headerName parsingTarget{regexCaptureGroup type __typename}__typename}roleLocation{type headerName parsingTarget{regexCaptureGroup type __typename}__typename}__typename}__typename}total __typename}}"
	responseStr, err := ExecuteQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}
	ruleDetails:=getRuleDetailsFromRulesListUsingIdName(response,"userAttributionRules" ,id)
	if len(ruleDetails)==0{
		d.SetId("")
		return nil
	}
	log.Printf("fetching from read %s",ruleDetails)
	name:=ruleDetails["name"].(string)
	d.Set("name",name)
	disabled:=ruleDetails["disabled"].(bool)
	d.Set("disabled",disabled)
	category:=ruleDetails["type"].(string)
	d.Set("category",category)
	scopeType:=ruleDetails["scopeType"]
	if scopeType=="SYSTEM_WIDE"{
		d.Set("scope_type", "SYSTEM_WIDE")
		// d.Set("url_regex",nil)
		// d.Set("environment",nil)
	}else{
		envScope := ruleDetails["customScope"].(map[string]interface{})["environmentScopes"]
		urlScope := ruleDetails["customScope"].(map[string]interface{})["urlScopes"]
		if len(envScope.([]interface{}))==0{
			d.Set("scope_type","CUSTOM")
			d.Set("url_regex",urlScope.([]interface{})[0].(map[string]interface{})["urlMatchRegex"])
			// d.Set("environment",nil)
		}else{
			d.Set("scope_type","CUSTOM")
			d.Set("environment",envScope.([]interface{})[0].(map[string]interface{})["environmentName"])
			// d.Set("url_regex",nil)
		}
	}
	requestHeaderDetails:= ruleDetails["requestHeader"]
	if requestHeaderDetails!=nil{
		auth_type:=requestHeaderDetails.(map[string]interface{})["authentication"]
		if auth_type!=nil{
			auth_type=auth_type.(map[string]interface{})["type"]
			d.Set("auth_type",auth_type)
		}else{
			d.Set("auth_type",nil)
		}
		user_id_location:=requestHeaderDetails.(map[string]interface{})["userIdLocation"].(map[string]interface{})["headerName"]
		user_id_regex_capture_group_details,ok:=requestHeaderDetails.(map[string]interface{})["userIdLocation"]
		if ok {
			if user_id_regex_capture_group,ok:=user_id_regex_capture_group_details.(map[string]interface{})["parsingTarget"]; ok {
				d.Set("user_id_regex_capture_group",user_id_regex_capture_group)
			}
		}
	
		if user_role_location_details, ok := requestHeaderDetails.(map[string]interface{})["roleLocation"]; ok && user_role_location_details!=nil{
			if user_role_location,ok:=user_role_location_details.(map[string]interface{})["headerName"]; ok{
				d.Set("user_role_location",user_role_location)
			}
			if role_location_regex_capture_group:=user_role_location_details.(map[string]interface{})["parsingTarget"]; ok {
				d.Set("role_location_regex_capture_group",role_location_regex_capture_group)
			}
		}
		d.Set("user_id_location",user_id_location)
	}
	return nil
}

func resourceUserAttributionRuleRequestHeaderUpdate(d *schema.ResourceData, meta interface{}) error {
	id:=d.Id()
	name := d.Get("name").(string)
	scopeType:=d.Get("scope_type").(string)
	disabled := d.Get("disabled").(bool)
	user_id_location:=d.Get("user_id_location").(string)
	user_id_regex_capture_group:=d.Get("user_id_regex_capture_group").(string)
	user_role_location:=d.Get("user_role_location").(string)
	role_location_regex_capture_group:=d.Get("role_location_regex_capture_group").(string)
	category:=d.Get("category").(string)
	auth_type:=d.Get("auth_type").(string)
	var authTypeQuery string
	if auth_type!=""{
		authTypeQuery=fmt.Sprintf(`authentication: { type: "%s" }`,auth_type)
	}else{
		authTypeQuery=""
	}
	if user_role_location=="" && role_location_regex_capture_group!=""{
		return fmt.Errorf("role_location_regex_capture_group is not expected here without user_role_location")
	}
	parsingTargetUserId:=""
	if user_id_regex_capture_group!=""{
		parsingTargetUserId=fmt.Sprintf(`parsingTarget: {
										type: REGEX_CAPTURE_GROUP
										regexCaptureGroup: "%s"
									}`,user_id_regex_capture_group)
	}
	if user_role_location=="" && role_location_regex_capture_group!=""{
		return fmt.Errorf("role_location_regex_capture_group is not expected here without user_role_location")
	}
	parsingTargetUserRole:=""
	if user_role_location!=""{
		parsingTargetUserRole=fmt.Sprintf(`roleLocation: {
			type: HEADER
			headerName: "%s"
			parsingTarget: {
				type: REGEX_CAPTURE_GROUP
				regexCaptureGroup: "%s"
			}
		}`,user_role_location,role_location_regex_capture_group)
	}

	readQuery:="{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName}urlScopes{urlMatchRegex}}}}}"
	readQueryResStr, err := ExecuteQuery(readQuery, meta)
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
			  rule: {
				name: "%s", 
				type: %s,
				id:"%s",
				disabled: %t,
				rank:%d, 
				scopeType: %s,
				requestHeader: {
					%s
					userIdLocation: {
						type: HEADER
						headerName: "%s"
						%s
					}
					%s
				}
			}
			) {
				id
				scopeType
				rank
				name
			}
		  }`,name,category,id,disabled,rank,scopeType,authTypeQuery,user_id_location,parsingTargetUserId,parsingTargetUserRole)
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
			  rule: {name: "%s", 
			  type: %s,
			  id:"%s",
			  disabled: %t,
			  rank:%d, 
			  scopeType: CUSTOM, 
			  %s,
			  requestHeader: {
				%s
				userIdLocation: {
					type: HEADER
					headerName: "%s"
					%s
				}
				%s
			  }
			}
			) {
				id
				scopeType
				rank
				name
				type
			}
		  }`,name,category,id,disabled,rank,scopedQuery,authTypeQuery,user_id_location,parsingTargetUserId,parsingTargetUserRole)
	}else{
		return fmt.Errorf("Expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

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
	updatedId:=rules["id"].(string)
	d.SetId(updatedId)
	return nil
}

func resourceUserAttributionRuleRequestHeaderDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := ExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}