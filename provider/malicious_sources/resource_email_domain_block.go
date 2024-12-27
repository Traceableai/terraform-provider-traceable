package malicious_sources

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
)

func ResourceEmailDomainBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailDomainBlockCreate,
		Read:   resourceEmailDomainBlockRead,
		Update: resourceEmailDomainBlockUpdate,
		Delete: resourceEmailDomainBlockDelete,

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
				Default:     "BLOCK",
			},
			"event_severity": {
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"expiration": {
				Type:        schema.TypeString,
				Description: "expiration for Block action",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := common.ConvertDurationToSeconds(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid duration in seconds or ISO 8601 format: %s", key, err))
					}
					return
				},
				StateFunc: func(val interface{}) string {
					v := val.(string)
					converted, _ := common.ConvertDurationToSeconds(v)
					return converted
				},
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

func resourceEmailDomainBlockCreate(d *schema.ResourceData, meta interface{}) error {
	expiration := d.Get("expiration").(string)
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

	exipiryDurationString := ReturnMalicousSourcesExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	emailFraudScoreQuery := ReturnEmailFraudScoreQuery(email_fraud_score)
	query := fmt.Sprintf(CREATE_EMAIL_DOMAIN_BLOCK, name, description, event_severity, rule_action, exipiryDurationString, data_leaked_email, disposable_email_domain, strings.Join(common.InterfaceToStringSlice(email_domains), ","), strings.Join(common.InterfaceToStringSlice(email_regexes), ","), emailFraudScoreQuery, envQuery)
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(query, meta)
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Println("Error:", err)
	}
	id := response["data"].(map[string]interface{})["createMaliciousSourcesRule"].(map[string]interface{})["id"].(string)
	d.SetId(id)
	return nil
}

func resourceEmailDomainBlockRead(d *schema.ResourceData, meta interface{}) error {
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
		if expirationDetails, ok := action["expirationDetails"].(map[string]interface{}); ok {
			d.Set("expiration", expirationDetails["expirationDuration"].(string))
		} else {
			d.Set("expiration", "")
		}
		d.Set("rule_action", action["ruleActionType"])
	}

	condition := ruleData["conditions"].([]interface{})[0].(map[string]interface{})
	emailDomainCondition := condition["emailDomainCondition"].(map[string]interface{})
	d.Set("data_leaked_email", emailDomainCondition["dataLeakedEmail"])
	d.Set("disposable_email_domain", emailDomainCondition["disposableEmailDomain"])
	d.Set("email_regexes", emailDomainCondition["emailRegexes"].([]interface{}))
	d.Set("email_domains", emailDomainCondition["emailDomains"].([]interface{}))
	if emailFraudScore, ok := emailDomainCondition["emailFraudScore"].(map[string]interface{}); ok {
		d.Set("email_fraud_score", emailFraudScore["minEmailFraudScoreLevel"])
	} else {
		d.Set("email_fraud_score", "")
	}

	if ruleScope, ok := ruleData["scope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				d.Set("environment", environmentIds)
			} else {
				d.Set("environment", []interface{}{})
			}
		}
	}
	return nil
}

func resourceEmailDomainBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	expiration := d.Get("expiration").(string)
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

	exipiryDurationString := ReturnMalicousSourcesExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	emailFraudScoreQuery := ReturnEmailFraudScoreQuery(email_fraud_score)
	query := fmt.Sprintf(UPDATE_EMAIL_DOMAIN_BLOCK, id, name, description, event_severity, rule_action, exipiryDurationString, data_leaked_email, disposable_email_domain, strings.Join(common.InterfaceToStringSlice(email_domains), ","), strings.Join(common.InterfaceToStringSlice(email_regexes), ","), emailFraudScoreQuery, envQuery)
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(query, meta)
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Println("Error:", err)
	}
	update_id := response["data"].(map[string]interface{})["updateMaliciousSourcesRule"].(map[string]interface{})["id"].(string)
	d.SetId(update_id)
	return nil
}

func resourceEmailDomainBlockDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteMaliciousSourcesRule(d, meta)
	return nil
}
