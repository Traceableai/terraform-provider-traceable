package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserAttributionCustomTokenRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributionRuleCustomTokenCreate,
		Read:   resourceUserAttributionRuleCustomTokenRead,
		Update: resourceUserAttributionRuleCustomTokenUpdate,
		Delete: resourceUserAttributionRuleCustomTokenDelete,

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
			"auth_type": {
				Type:        schema.TypeString,
				Description: "auth type",
				Required:    true,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "custom token location (REQUEST_HEADER or REQUEST_BODY or REQUEST_COOKIE)",
				Required:    true,
			},
			"token_name": {
				Type:        schema.TypeString,
				Description: "token name",
				Required:    true,
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
				Default:     "CUSTOM_TOKEN",
			},
		},
	}
}

func resourceUserAttributionRuleCustomTokenCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	scope_type := d.Get("scope_type").(string)
	environment := d.Get("environment").(string)
	url_regex:=d.Get("url_regex").(string)
	auth_type := d.Get("auth_type").(string)
	location:=d.Get("location").(string)
	token_name:=d.Get("token_name").(string)
	category:=d.Get("category").(string)

	if scope_type!="SYSTEM_WIDE" && scope_type!="CUSTOM"{
		return fmt.Errorf("scope_type supported string is SYSTEM_WIDE or CUSTOM")
	}
	if location!="REQUEST_HEADER" && location!="REQUEST_BODY" && location!="REQUEST_COOKIE"{
		return fmt.Errorf("location supported string is REQUEST_BODY or REQUEST_HEADER")
	}

	customScopeString:=""
	if scope_type=="CUSTOM"{
		if environment!="" && url_regex==""{
			customScopeString=fmt.Sprintf(`customScope: { environmentScopes: [{ environmentName: "%s" }] }`,environment)
		}else if environment=="" && url_regex!=""{
			customScopeString=fmt.Sprintf(`customScope: { urlScopes: [{ urlMatchRegex: "%s" }] }`,url_regex)
		}else{
			return fmt.Errorf("Required environment or url_regex")
		}
	}

	tokenLocationString:=""
	if location=="REQUEST_HEADER"{
		tokenLocationString=fmt.Sprintf(`requestHeaderLocation: { headerName: "%s", type: HEADER }`,token_name)
	}else if location=="REQUEST_BODY"{
		tokenLocationString=fmt.Sprintf(`requestBodyLocation: { jsonPath: "%s", type: JSON_PATH }`,token_name)
	}else if location=="REQUEST_COOKIE"{
		location="REQUEST_HEADER"
		tokenLocationString=fmt.Sprintf(`requestHeaderLocation: { cookieName: "%s", type: COOKIE }`,token_name)
	}
	

	var query string
	query = fmt.Sprintf(`mutation {
		createUserAttributionRule(
		  input: {
		  name: "%s", 
		  type: %s, 
		  scopeType: %s,
		  customToken:{
			authentication: { type: "%s" }
			customTokenLocation: %s
			%s
		  }
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
	  }`,name,category,scope_type,auth_type,location,tokenLocationString,customScopeString)
	
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

func resourceUserAttributionRuleCustomTokenRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)
	readQuery:="{ userAttributionRules { results { id scopeType rank name type disabled customScope { environmentScopes { environmentName __typename } urlScopes { urlMatchRegex __typename } __typename } customToken { authentication { type __typename } customTokenLocation requestBodyLocation { jsonPath type __typename } requestHeaderLocation { cookieName headerName type __typename } __typename } } } }"
	responseStr, err := ExecuteQuery(readQuery, meta)
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
	scopeType:=ruleDetails["scopeType"].(string)
	d.Set("name",name)
	category:=ruleDetails["type"].(string)
	d.Set("category",category)
	disabled:=ruleDetails["disabled"].(bool)
	d.Set("disabled",disabled)
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
	customTokenDetails:=ruleDetails["customToken"]
	if customTokenDetails!=nil{
		log.Printf("This is details %s",customTokenDetails)
		auth_type:=customTokenDetails.(map[string]interface{})["authentication"]
		customTokenLocation:=customTokenDetails.(map[string]interface{})["customTokenLocation"]
		if auth_type!=nil{
			auth_type=auth_type.(map[string]interface{})["type"]
			d.Set("auth_type",auth_type)
		}else{
			d.Set("auth_type",nil)
		}
		if customTokenLocation=="REQUEST_HEADER"{
			requestHeaderLocation:=customTokenDetails.(map[string]interface{})["requestHeaderLocation"]
			headerType:=requestHeaderLocation.(map[string]interface{})["type"]
			if headerType=="COOKIE"{
				d.Set("location","REQUEST_COOKIE")
				d.Set("token_name",requestHeaderLocation.(map[string]interface{})["cookieName"])
	
			}else{
				d.Set("location","REQUEST_HEADER")
				d.Set("token_name",requestHeaderLocation.(map[string]interface{})["headerName"])
			}
		}else{
			d.Set("location","REQUEST_BODY")
			requestBodyLocation:=customTokenDetails.(map[string]interface{})["requestBodyLocation"]
			d.Set("token_name",requestBodyLocation.(map[string]interface{})["jsonPath"])
		}
	}

	return nil
}

func resourceUserAttributionRuleCustomTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	id:=d.Id()
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
	name := d.Get("name").(string)
	scope_type := d.Get("scope_type").(string)
	environment := d.Get("environment").(string)
	url_regex:=d.Get("url_regex").(string)
	auth_type := d.Get("auth_type").(string)
	location:=d.Get("location").(string)
	token_name:=d.Get("token_name").(string)
	disabled := d.Get("disabled").(bool)
	category := d.Get("category").(string)
	if scope_type!="SYSTEM_WIDE" && scope_type!="CUSTOM"{
		return fmt.Errorf("scope_type supported string is SYSTEM_WIDE or CUSTOM")
	}
	if location!="REQUEST_HEADER" && location!="REQUEST_BODY" && location!="REQUEST_COOKIE"{
		return fmt.Errorf("location supported string is REQUEST_BODY or REQUEST_HEADER")
	}

	customScopeString:=""
	if scope_type=="CUSTOM"{
		if environment!="" && url_regex==""{
			customScopeString=fmt.Sprintf(`customScope: { environmentScopes: [{ environmentName: "%s" }] }`,environment)
		}else if environment=="" && url_regex!=""{
			customScopeString=fmt.Sprintf(`customScope: { urlScopes: [{ urlMatchRegex: "%s" }] }`,url_regex)
		}else{
			return fmt.Errorf("Required environment or url_regex")
		}
	}

	tokenLocationString:=""
	if location=="REQUEST_HEADER"{
		tokenLocationString=fmt.Sprintf(`requestHeaderLocation: { headerName: "%s", type: HEADER }`,token_name)
	}else if location=="REQUEST_BODY"{
		tokenLocationString=fmt.Sprintf(`requestBodyLocation: { jsonPath: "%s", type: JSON_PATH }`,token_name)
	}else if location=="REQUEST_COOKIE"{
		location="REQUEST_HEADER"
		tokenLocationString=fmt.Sprintf(`requestHeaderLocation: { cookieName: "%s", type: COOKIE }`,token_name)
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
			customToken:{
				authentication: { type: "%s" }
				customTokenLocation: %s
				%s
			}
			%s
		}
		) {
			id
			scopeType
			rank
			name
			type
		}
	  }`,id,rank,name,disabled,category,scope_type,auth_type,location,tokenLocationString,customScopeString)

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

func resourceUserAttributionRuleCustomTokenDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation { deleteUserAttributionRule(input: {id: \"%s\"}) { results { id scopeType rank name type disabled } } }", id)
	_, err := ExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}