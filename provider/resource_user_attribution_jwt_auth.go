package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributionJwtAuthRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleJwtAuthCreate,
		Read:   resourceUserAttributionRuleJwtAuthRead,
		Update: resourceUserAttributionRuleJwtAuthUpdate,
		Delete: resourceUserAttributionRuleJwtAuthDelete,

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
			"jwt_location": {
				Type:        schema.TypeString,
				Description: "header or cookie",
				Required:    true,
			},
			"jwt_key": {
				Type:        schema.TypeString,
				Description: "header name for jwt in header or cookie",
				Required:    true,
			},
			"token_capture_group": {
				Type:        schema.TypeString,
				Description: "token capture group",
				Optional:    true,
			},
			"user_id_location_json_path": {
				Type:        schema.TypeString,
				Description: "user id location json path",
				Optional:    true,
			},
			"user_id_claim": {
				Type:        schema.TypeString,
				Description: "user id claim",
				Required:    true,
			},
			"user_role_location_json_path": {
				Type:        schema.TypeString,
				Description: "user role location json path",
				Optional:    true,
			},
			"user_role_claim": {
				Type:        schema.TypeString,
				Description: "user role claim",
				Optional:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceUserAttributionRuleJwtAuthCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	scopeType:=d.Get("scope_type").(string)
	environment:=d.Get("environment").(string)
	url_regex:=d.Get("url_regex").(string)
	token_capture_group:=d.Get("token_capture_group").(string)
	jwt_location:=d.Get("jwt_location").(string)
	jwt_key:=d.Get("jwt_key").(string)
	user_id_claim:=d.Get("user_id_claim").(string)
	user_role_claim:=d.Get("user_role_claim").(string)
	user_role_location_json_path:=d.Get("user_role_location_json_path").(string)
	user_id_location_json_path:=d.Get("user_id_location_json_path").(string)
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

	tokenCaptureGroup:=""
	if token_capture_group!=""{
		tokenCaptureGroup=fmt.Sprintf(`parsingTarget: {
										type: REGEX_CAPTURE_GROUP
										regexCaptureGroup: "%s"
									}`,token_capture_group)
	}
	userIdLocationString:=""
	if user_id_location_json_path!=""{
		userIdLocationString=fmt.Sprintf(`userIdLocation: { type: JSON_PATH, jsonPath: "%s" }`,user_id_location_json_path)
	}
	roleLocationString:=""
	if user_role_location_json_path!=""{
		roleLocationString=fmt.Sprintf(`roleLocation: { type: JSON_PATH, jsonPath: "%s" }`,user_role_location_json_path)
	}
	location_key:="headerName"
	if jwt_location=="COOKIE"{
		location_key="cookieName"
	}
	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {
				name: "%s",
				type: JWT,
				scopeType: %s,
				jwt: {
					%s
					location: {
						type: %s
						%s: "%s"
						%s
					}
					%s
					%s
					roleClaim: "%s"
					userIdClaim: "%s"
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
		  }`,name,scopeType,authTypeQuery,jwt_location,location_key,jwt_key,tokenCaptureGroup,userIdLocationString,roleLocationString,user_role_claim,user_id_claim)
	} else if scopeType== "CUSTOM" {
		
		if scopedQuery==""{
			return fmt.Errorf("Provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			createUserAttributionRule(
			  input: {
			  name: "%s", 
			  type: JWT, 
			  scopeType: CUSTOM, 
			  jwt: {
				%s
				location: {
					type: %s
					%s: "%s"
					%s
				}
				%s
				%s
				roleClaim: "%s"
				userIdClaim: "%s"
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
		  }`,name,authTypeQuery,jwt_location,location_key,jwt_key,tokenCaptureGroup,userIdLocationString,roleLocationString,user_role_claim,user_id_claim,scopedQuery)
	}else{
		return fmt.Errorf("Expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
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

func resourceUserAttributionRuleJwtAuthRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)
	readQuery:="{userAttributionRules{results{id scopeType rank name type disabled customScope{environmentScopes{environmentName __typename} urlScopes{urlMatchRegex __typename} __typename} jwt{authentication{type __typename} location{type headerName cookieName parsingTarget{regexCaptureGroup type __typename} __typename} roleClaim userIdClaim userIdLocation{type jsonPath __typename} roleLocation{type jsonPath __typename} __typename}}}}"
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}
	log.Printf("Response from read %s",responseStr)
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
			log.Printf("adityaaa %s",urlScope.([]interface{})[0].(map[string]interface{})["urlMatchRegex"])
			d.Set("url_regex",urlScope.([]interface{})[0].(map[string]interface{})["urlMatchRegex"])
			// d.Set("environment",nil)
		}else{
			d.Set("scope_type","CUSTOM")
			d.Set("environment",envScope.([]interface{})[0].(map[string]interface{})["environmentName"])
			// d.Set("url_regex",nil)
		}
	}
	if ruleDetails["type"].(string)!="JWT"{
		d.Set("auth_type",nil)
		d.Set("jwt_location",nil)
		d.Set("jwt_key",nil)
		d.Set("token_capture_group",nil)
		d.Set("user_id_location_json_path",nil)
		d.Set("user_role_location_json_path",nil)
		d.Set("user_id_claim",nil)
		d.Set("user_role_claim",nil)
		return nil
	}
	auth_type:=ruleDetails["jwt"].(map[string]interface{})["authentication"]
	if auth_type!=nil{
		auth_type=auth_type.(map[string]interface{})["type"]
		d.Set("auth_type",auth_type)
	}else{
		d.Set("auth_type",nil)
	}
	

	jwt_location:=ruleDetails["jwt"].(map[string]interface{})["location"].(map[string]interface{})["type"]
	d.Set("jwt_location",jwt_location)
	if jwt_location=="HEADER"{
		jwt_key:=ruleDetails["jwt"].(map[string]interface{})["location"].(map[string]interface{})["headerName"]
		d.Set("jwt_key",jwt_key)
	}else{
		jwt_key:=ruleDetails["jwt"].(map[string]interface{})["location"].(map[string]interface{})["cookieName"]
		d.Set("jwt_key",jwt_key)
	}
	token_capture_group_details:=ruleDetails["jwt"].(map[string]interface{})["location"].(map[string]interface{})["parsingTarget"]
    if token_capture_group_details!=nil{
		token_capture_group:=token_capture_group_details.(map[string]interface{})["regexCaptureGroup"]
		d.Set("token_capture_group",token_capture_group)
	}

	userIdLocationDetails:=ruleDetails["jwt"].(map[string]interface{})["userIdLocation"]
	if userIdLocationDetails!=nil{
		user_id_location_json_path:=userIdLocationDetails.(map[string]interface{})["jsonPath"]
		d.Set("user_id_location_json_path",user_id_location_json_path)
	}

	userRoleLocationDetails:=ruleDetails["jwt"].(map[string]interface{})["roleLocation"]
	if userRoleLocationDetails!=nil{
		user_role_location_json_path:=userRoleLocationDetails.(map[string]interface{})["jsonPath"]
		d.Set("user_role_location_json_path",user_role_location_json_path)
	}

	user_id_claim:=ruleDetails["jwt"].(map[string]interface{})["userIdClaim"]
	d.Set("user_id_claim",user_id_claim)
	
	user_role_claim:=ruleDetails["jwt"].(map[string]interface{})["roleClaim"]
	d.Set("user_role_claim",user_role_claim)

	return nil
}

func resourceUserAttributionRuleJwtAuthUpdate(d *schema.ResourceData, meta interface{}) error {
	id:=d.Id()
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
	name := d.Get("name").(string)
	scopeType:=d.Get("scope_type").(string)
	environment:=d.Get("environment").(string)
	url_regex:=d.Get("url_regex").(string)
	token_capture_group:=d.Get("token_capture_group").(string)
	jwt_location:=d.Get("jwt_location").(string)
	jwt_key:=d.Get("jwt_key").(string)
	user_id_claim:=d.Get("user_id_claim").(string)
	user_role_claim:=d.Get("user_role_claim").(string)
	user_role_location_json_path:=d.Get("user_role_location_json_path").(string)
	user_id_location_json_path:=d.Get("user_id_location_json_path").(string)
	auth_type:=d.Get("auth_type").(string)
	disabled := d.Get("disabled").(bool)
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

	tokenCaptureGroup:=""
	if token_capture_group!=""{
		tokenCaptureGroup=fmt.Sprintf(`parsingTarget: {
										type: REGEX_CAPTURE_GROUP
										regexCaptureGroup: "%s"
									}`,token_capture_group)
	}
	userIdLocationString:=""
	if user_id_location_json_path!=""{
		userIdLocationString=fmt.Sprintf(`userIdLocation: { type: JSON_PATH, jsonPath: "%s" }`,user_id_location_json_path)
	}
	roleLocationString:=""
	if user_role_location_json_path!=""{
		roleLocationString=fmt.Sprintf(`roleLocation: { type: JSON_PATH, jsonPath: "%s" }`,user_role_location_json_path)
	}
	location_key:="headerName"
	if jwt_location=="COOKIE"{
		location_key="cookieName"
	}
	var query string
	if scopeType == "SYSTEM_WIDE" {
		query = fmt.Sprintf(`mutation {
			updateUserAttributionRule(
			  rule: {
				id:"%s",
			  	rank:%d, 
				name: "%s",
				disabled: %t,
				type: JWT,
				scopeType: %s,
				jwt: {
					%s
					location: {
						type: %s
						%s: "%s"
						%s
					}
					%s
					%s
					roleClaim: "%s"
					userIdClaim: "%s"
				}
			}
			) {
				id
				scopeType
				rank
				name
			}
		  }`,id,rank,name,disabled,scopeType,authTypeQuery,jwt_location,location_key,jwt_key,tokenCaptureGroup,userIdLocationString,roleLocationString,user_role_claim,user_id_claim)
	} else if scopeType== "CUSTOM" {
		
		if scopedQuery==""{
			return fmt.Errorf("Provide enviroment or url regex for custom scoped user attribution or remove one of them")
		}
		query = fmt.Sprintf(`mutation {
			updateUserAttributionRule(
			  rule: {
			  id:"%s",
			  rank:%d, 
			  name: "%s", 
			  type: JWT, 
			  scopeType: CUSTOM, 
			  jwt: {
				%s
				location: {
					type: %s
					%s: "%s"
					%s
				}
				%s
				%s
				roleClaim: "%s"
				userIdClaim: "%s"
			  },
			  %s
			}
			) {
				id
				scopeType
				rank
				name
				type
			}
		  }`,id,rank,name,authTypeQuery,jwt_location,location_key,jwt_key,tokenCaptureGroup,userIdLocationString,roleLocationString,user_role_claim,user_id_claim,scopedQuery)
	}else{
		return fmt.Errorf("Expected values are CUSTOM or SYSTEM_WIDE for user attribution scope type")
	}
	

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
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

func resourceUserAttributionRuleJwtAuthDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := executeQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}