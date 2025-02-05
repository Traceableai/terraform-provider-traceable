package malicious_sources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"log"
)

func ResourceEmailDomainAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailDomainAlertCreate,
		Read:   resourceEmailDomainAlertRead,
		Update: resourceEmailDomainAlertUpdate,
		Delete: resourceEmailDomainAlertDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the email domain policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the policy",
				Optional:    true,
			},
			"rule_action": {
				Type:        schema.TypeString,
				Description: "Need to provide the action to be performed ",
				Optional:    true,
				Default:     "ALERT",
			},
			"event_severity": {
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeSet,
				Description: "environment where it will be applied",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_leaked_email": {
				Type:        schema.TypeBool,
				Description: "Users from leaked email domain",
				Required:    true,
			},
			"disposable_email_domain": {
				Type:        schema.TypeBool,
				Description: "Users from disposable email domain",
				Required:    true,
			},
			"email_domains": {
				Type:        schema.TypeSet,
				Description: "Email domains for rule",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"email_regexes": {
				Type:        schema.TypeSet,
				Description: "Email regexes for rule",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"email_fraud_score": {
				Type:        schema.TypeString,
				Description: "Minimum email fraud score",
				Required:    true,
			},
		},
	}
}

func resourceEmailDomainAlertCreate(d *schema.ResourceData, meta interface{}) error {
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	rule_action := d.Get("rule_action").(string)
	data_leaked_email := d.Get("data_leaked_email").(bool)
	email_fraud_score := d.Get("email_fraud_score").(string)
	disposable_email_domain := d.Get("disposable_email_domain").(bool)
	email_domains := d.Get("email_domains").(*schema.Set).List()
	email_regexes := d.Get("email_regexes").(*schema.Set).List()
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	emailFraudScoreQuery := ReturnEmailFraudScoreQuery(email_fraud_score)
	query := fmt.Sprintf(CREATE_EMAIL_DOMAIN_ALERT, name, description, event_severity, rule_action, data_leaked_email, disposable_email_domain, common.InterfaceToStringSlice(email_domains), common.InterfaceToStringSlice(email_regexes), emailFraudScoreQuery, envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createMaliciousSourcesRule")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId(id)
	return nil
}

func resourceEmailDomainAlertRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(MALICOUS_SOURCES_READ, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "maliciousSourcesRules", id)
	if len(ruleData) == 0 {
		d.SetId("")
		return nil
	}
	ruleDetails := ruleData["info"].(map[string]interface{})
	d.Set("name", ruleDetails["name"].(string))
	d.Set("description", ruleDetails["description"].(string))
	if action, ok := ruleDetails["action"].(map[string]interface{}); ok {
		d.Set("event_severity", action["eventSeverity"])
		d.Set("rule_action", action["ruleActionType"])
	}

	condition := ruleDetails["conditions"].([]interface{})[0].(map[string]interface{})
	emailDomainCondition := condition["emailDomainCondition"].(map[string]interface{})
	d.Set("data_leaked_email", emailDomainCondition["dataLeakedEmail"])
	d.Set("disposable_email_domain", emailDomainCondition["disposableEmailDomain"])
	d.Set("email_regexes", emailDomainCondition["emailRegexes"].([]interface{}))
	d.Set("email_domains", emailDomainCondition["emailDomains"].([]interface{}))
	emailFraudScoreFlag := true
	if emailFraudScore, ok := emailDomainCondition["emailFraudScore"].(map[string]interface{}); ok {
		d.Set("email_fraud_score", emailFraudScore["minEmailFraudScoreLevel"])
		emailFraudScoreFlag = false
	} 
	if emailFraudScoreFlag {
		d.Set("email_fraud_score",false)
	}

	envFlag := true
	if ruleScope, ok := ruleDetails["scope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				d.Set("environment", environmentIds)
				envFlag=false
			} 
		}
	}
	if envFlag{
		d.Set("environment", []interface{}{})
	}
	return nil
}

func resourceEmailDomainAlertUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	rule_action := d.Get("rule_action").(string)
	data_leaked_email := d.Get("data_leaked_email").(bool)
	email_fraud_score := d.Get("email_fraud_score").(string)
	disposable_email_domain := d.Get("disposable_email_domain").(bool)
	email_domains := d.Get("email_domains").(*schema.Set).List()
	email_regexes := d.Get("email_regexes").(*schema.Set).List()
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	envQuery := ReturnEnvScopedQuery(environment)
	emailFraudScoreQuery := ReturnEmailFraudScoreQuery(email_fraud_score)
	query := fmt.Sprintf(UPDATE_EMAIL_DOMAIN_ALERT, id, name, description, event_severity, rule_action, data_leaked_email, disposable_email_domain, common.InterfaceToStringSlice(email_domains), common.InterfaceToStringSlice(email_regexes), emailFraudScoreQuery, envQuery)

	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateMaliciousSourcesRule")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId(updatedId)
	return nil
}

func resourceEmailDomainAlertDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteMaliciousSourcesRule(d, meta)
	return nil
}
