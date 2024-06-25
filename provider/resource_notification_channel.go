package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationChannelRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationChannelCreate,
		Read:   resourceNotificationChannelRead,
		Update: resourceNotificationChannelUpdate,
		Delete: resourceNotificationChannelDelete,

		Schema: map[string]*schema.Schema{
			"channel_name": {
				Type:        schema.TypeString,
				Description: "Name of the notification channel",
				Required:    true,
			},
			"email": {
				Type:        schema.TypeSet,
				Description: "Email address for notification channel",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"slack_webhook": {
				Type:        schema.TypeString,
				Description: "Slack webhook config for notification channel",
				Optional:    true,
			},
			"splunk_id": {
				Type:        schema.TypeString,
				Description: "Your splunk integration id",
				Optional:    true,
			},
			"syslog_id": {
				Type:        schema.TypeString,
				Description: "Your syslog integration id",
				Optional:    true,
			},
			"custom_webhook": {
				Type: schema.TypeList,	
				Description: "Your custom webhook url",
				Optional: true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"webhook_url": {
							Type:     schema.TypeString,
							Description: "Url of custom webhook",
							Required: true,
						},
						"custom_webhook_headers": {
							Type:        schema.TypeSet,
							Description: "Headers for custom webhook",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Description: "Test header key",
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Description: "Test header value",
										Required: true,
									},
									"is_secret": {
										Type:     schema.TypeBool,
										Description: "Header is secret or not",
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"s3_webhook": {
				Type:        schema.TypeSet,
				Description: "S3 bucket configuration for notification channel",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Description: "Name of your s3 bucket",
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Description: "Region of your s3 bucket",
							Required: true,
						},
						"bucket_arn": {
							Type:     schema.TypeString,
							Description: "S3 bucket arn",
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNotificationChannelCreate(d *schema.ResourceData, meta interface{}) error {
	channelName := d.Get("channel_name").(string)

	var emails []string
	if v, ok := d.GetOk("email"); ok {
		for _, email := range v.(*schema.Set).List() {
			emails = append(emails, email.(string))
		}
	}

	slackWebhook := d.Get("slack_webhook").(string)
	splunkID := d.Get("splunk_id").(string)
	syslogID := d.Get("syslog_id").(string)

	var customWebhookURL string
	var customWebhookHeaders []map[string]interface{}

	if v, ok := d.GetOk("custom_webhook"); ok {
		customWebhook := v.([]interface{})[0].(map[string]interface{})
		customWebhookURL = customWebhook["webhook_url"].(string)

		log.Println(customWebhook["custom_webhook_headers"])
		if headers, ok := customWebhook["custom_webhook_headers"]; ok {
			if headersSet, ok := headers.(*schema.Set); ok {
				for _, header := range headersSet.List() {
					h := header.(map[string]interface{})
					customWebhookHeader := make(map[string]interface{})
					customWebhookHeader["key"] = h["key"].(string)
					customWebhookHeader["value"] = h["value"].(string)
					customWebhookHeader["is_secret"] = h["is_secret"].(bool)
					customWebhookHeaders = append(customWebhookHeaders, customWebhookHeader)
				}
			}
		}
	}

	var s3WebhookBucketName, s3WebhookRegion, s3WebhookBucketArn string
	if v, ok := d.GetOk("s3_webhook"); ok {
		s3Webhook := v.(*schema.Set).List()[0].(map[string]interface{})
		s3WebhookBucketName = s3Webhook["bucket_name"].(string)
		s3WebhookRegion = s3Webhook["region"].(string)
		s3WebhookBucketArn = s3Webhook["bucket_arn"].(string)
	}

	slackWebhookChannelConfigsQuery:=""
	if slackWebhook!="" {
		slackWebhookChannelConfigsQuery=fmt.Sprintf(`slackWebhookChannelConfigs: [{ url: "%s" }]`,slackWebhook)
	}

	s3BucketChannelConfigsQuery:=""
	if s3WebhookBucketName!=""{
		s3BucketChannelConfigsQuery=fmt.Sprintf(`s3BucketChannelConfigs: [
			{
				bucketName: "%s"
				region: "%s"
				authenticationCredentialType: WEB_IDENTITY
				webIdentityAuthenticationCredential: { roleArn: "%s" }
			}
		]`,s3WebhookBucketName,s3WebhookRegion,s3WebhookBucketArn)
	}

	splunkIntegrationChannelConfigsQuery:=""
	if splunkID!=""{
		splunkIntegrationChannelConfigsQuery=fmt.Sprintf(`splunkIntegrationChannelConfigs: [{ splunkIntegrationId: "%s" }]`,splunkID)
	}

	syslogIntegrationChannelConfigsQuery:=""
	if syslogID!=""{
		syslogIntegrationChannelConfigsQuery=fmt.Sprintf(`syslogIntegrationChannelConfigs: [{ syslogIntegrationId: "%s" }]`,syslogID)
	}
	
	customWebhookChannelConfigsQuery:=""
	if customWebhookURL!="" && len(customWebhookHeaders)==0{
		customWebhookChannelConfigsQuery=fmt.Sprintf(`customWebhookChannelConfigs: [
			{
				url: "%s"
				headers: []
				id: ""
			}
		]`,customWebhookURL)
	}else if customWebhookURL!="" && len(customWebhookHeaders)>0{
		headerString:="["
		for _, headers := range customWebhookHeaders {
			tmp := fmt.Sprintf(`{key: "%s", value: "%s", isSecret: %t}`, headers["key"], headers["value"],headers["is_secret"])
			headerString+=tmp
			headerString+=","
		}
		headerString=headerString[:len(headerString)-1]
		headerString+="]"

		customWebhookChannelConfigsQuery=fmt.Sprintf(`customWebhookChannelConfigs: [
			{
				url: "%s"
				headers: %s
				id: ""
			}
		]`,customWebhookURL,headerString)
	}

	emailChannelConfigsQuery:=""
	if len(emails)>0{
		emailString:="["
		for _,em :=range emails{
			tmp:=fmt.Sprintf(`{ address: "%s" }`,em)
			emailString+=tmp
			emailString+=","
		}
		emailString=emailString[:len(emailString)-1]
		emailString+="]"
		emailChannelConfigsQuery=fmt.Sprintf(`emailChannelConfigs: %s`,emailString)
	}

	if slackWebhookChannelConfigsQuery=="" && emailChannelConfigsQuery=="" && s3BucketChannelConfigsQuery=="" && splunkIntegrationChannelConfigsQuery=="" && syslogIntegrationChannelConfigsQuery=="" && customWebhookChannelConfigsQuery=="" && emailChannelConfigsQuery==""{
		return fmt.Errorf("No channel configuration provided")
	}

	query:=fmt.Sprintf(`mutation {
		createNotificationChannel(
			input: {
				channelName: "%s"
				notificationChannelConfig: {
					%s
					%s
					%s
					%s
					%s
					%s
				}
			}
		) {
			channelId
		}
	}
	`,channelName,slackWebhookChannelConfigsQuery,emailChannelConfigsQuery,s3BucketChannelConfigsQuery,splunkIntegrationChannelConfigsQuery,syslogIntegrationChannelConfigsQuery,customWebhookChannelConfigsQuery)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}
	rules := response["data"].(map[string]interface{})["createNotificationChannel"].(map[string]interface{})
	log.Println(rules)
	id:=rules["channelId"].(string)
	d.SetId(id)
 	return nil
}

func resourceNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	readQuery:=`{
		notificationChannels {
			results {
				channelId
				channelName
				notificationChannelConfig {
					emailChannelConfigs {
						address
					}
					customWebhookChannelConfigs {
						headers {
							key
							value
							isSecret
						}
						customWebhookChannelIndex: id
						url
					}
					slackWebhookChannelConfigs {
						url
					}
					splunkIntegrationChannelConfigs {
						splunkIntegrationId
					}
					syslogIntegrationChannelConfigs {
						syslogIntegrationId
					}
					s3BucketChannelConfigs {
						bucketName
						region
						authenticationCredentialType
						webIdentityAuthenticationCredential {
							roleArn	
						}	
					}	
				}	
			}
		}
	}`
	var response map[string]interface{}
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql query %s", readQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}
	id:=d.Id()
	ruleDetails:=getRuleDetailsFromRulesListUsingIdName(response,"notificationChannels" ,id,"channelId","channelName")
	if len(ruleDetails)==0{
		d.SetId("")
		return nil
	}
	log.Printf("channels %s",ruleDetails)

	channelName:=ruleDetails["channelName"]
	d.Set("channel_name",channelName)

	notificationChannelConfig:=ruleDetails["notificationChannelConfig"].(map[string]interface{})

	splunkIntegrationChannelConfigs:=notificationChannelConfig["splunkIntegrationChannelConfigs"].([]interface{})
	if len(splunkIntegrationChannelConfigs)>0{

		splunk_id:=splunkIntegrationChannelConfigs[0].(map[string]interface{})["splunkIntegrationId"]
		log.Printf("this is splunkkkk %s",splunk_id)
		d.Set("splunk_id",splunk_id)
	}

	syslogIntegrationChannelConfigs:=notificationChannelConfig["syslogIntegrationChannelConfigs"].([]interface{})
	if len(syslogIntegrationChannelConfigs)>0{

		syslog_id:=syslogIntegrationChannelConfigs[0].(map[string]interface{})["syslogIntegrationId"]
		log.Printf("this is syslogkkk %s",syslog_id)
		d.Set("syslog_id",syslog_id)
	}

	slackWebhookChannelConfigs:=notificationChannelConfig["slackWebhookChannelConfigs"].([]interface{})
	if len(slackWebhookChannelConfigs)>0{
		slack_url:=slackWebhookChannelConfigs[0].(map[string]interface{})["url"]
		log.Printf("this is slack_webhookkk %s",slack_url)
		d.Set("slack_webhook",slack_url)
	}

	emailChannelConfigs:=notificationChannelConfig["emailChannelConfigs"].([]interface{})
	if len(emailChannelConfigs)>0{
		var emails []interface{}
		for _,em := range emailChannelConfigs{
			emails=append(emails,em.(map[string]interface{})["address"])
		}
		d.Set("email", schema.NewSet(schema.HashString, emails))
	}
	s3BucketChannelConfigs := notificationChannelConfig["s3BucketChannelConfigs"].([]interface{})
	if len(s3BucketChannelConfigs) > 0 {
		s3Bucket := s3BucketChannelConfigs[0].(map[string]interface{})
		s3Data := make(map[string]interface{})
		s3Data["bucket_name"] = s3Bucket["bucketName"]
		s3Data["region"] = s3Bucket["region"]
		s3Data["bucket_arn"] = s3Bucket["webIdentityAuthenticationCredential"].(map[string]interface{})["roleArn"]
		d.Set("s3_webhook", []interface{}{s3Data})
	}

	customWebhookChannelConfigs := notificationChannelConfig["customWebhookChannelConfigs"].([]interface{})
	if len(customWebhookChannelConfigs) > 0 {
		customWebhook := customWebhookChannelConfigs[0].(map[string]interface{})
		headers := customWebhook["headers"].([]interface{})
		headerData := make([]interface{}, len(headers))
		for i, header := range headers {
			headerMap := header.(map[string]interface{})
			headerData[i] = map[string]interface{}{
				"key":       headerMap["key"],
				"value":     headerMap["value"],
				"is_secret": headerMap["isSecret"],
			}
		}
		
		customWebhookData := map[string]interface{}{
			"webhook_url":            customWebhook["url"],
			"custom_webhook_headers": headerData,
		}
		d.Set("custom_webhook",[]interface{}{customWebhookData})
	}
	return nil
}

func resourceNotificationChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	channel_id := d.Id()
	channelName := d.Get("channel_name").(string)

	var emails []string
	if v, ok := d.GetOk("email"); ok {
		for _, email := range v.(*schema.Set).List() {
			emails = append(emails, email.(string))
		}
	}

	slackWebhook := d.Get("slack_webhook").(string)
	splunkID := d.Get("splunk_id").(string)
	syslogID := d.Get("syslog_id").(string)

	var customWebhookURL string
	var customWebhookHeaders []map[string]interface{}

	if v, ok := d.GetOk("custom_webhook"); ok {
		customWebhook := v.([]interface{})[0].(map[string]interface{})
		customWebhookURL = customWebhook["webhook_url"].(string)

		log.Println(customWebhook["custom_webhook_headers"])
		if headers, ok := customWebhook["custom_webhook_headers"]; ok {
			if headersSet, ok := headers.(*schema.Set); ok {
				for _, header := range headersSet.List() {
					h := header.(map[string]interface{})
					customWebhookHeader := make(map[string]interface{})
					customWebhookHeader["key"] = h["key"].(string)
					customWebhookHeader["value"] = h["value"].(string)
					customWebhookHeader["is_secret"] = h["is_secret"].(bool)
					customWebhookHeaders = append(customWebhookHeaders, customWebhookHeader)
				}
			}
		}
	}

	var s3WebhookBucketName, s3WebhookRegion, s3WebhookBucketArn string
	if v, ok := d.GetOk("s3_webhook"); ok {
		s3Webhook := v.(*schema.Set).List()[0].(map[string]interface{})
		s3WebhookBucketName = s3Webhook["bucket_name"].(string)
		s3WebhookRegion = s3Webhook["region"].(string)
		s3WebhookBucketArn = s3Webhook["bucket_arn"].(string)
	}

	slackWebhookChannelConfigsQuery:=""
	if slackWebhook!="" {
		slackWebhookChannelConfigsQuery=fmt.Sprintf(`slackWebhookChannelConfigs: [{ url: "%s" }]`,slackWebhook)
	}

	s3BucketChannelConfigsQuery:=""
	if s3WebhookBucketName!=""{
		s3BucketChannelConfigsQuery=fmt.Sprintf(`s3BucketChannelConfigs: [
			{
				bucketName: "%s"
				region: "%s"
				authenticationCredentialType: WEB_IDENTITY
				webIdentityAuthenticationCredential: { roleArn: "%s" }
			}
		]`,s3WebhookBucketName,s3WebhookRegion,s3WebhookBucketArn)
	}

	splunkIntegrationChannelConfigsQuery:=""
	if splunkID!=""{
		splunkIntegrationChannelConfigsQuery=fmt.Sprintf(`splunkIntegrationChannelConfigs: [{ splunkIntegrationId: "%s" }]`,splunkID)
	}

	syslogIntegrationChannelConfigsQuery:=""
	if syslogID!=""{
		syslogIntegrationChannelConfigsQuery=fmt.Sprintf(`syslogIntegrationChannelConfigs: [{ syslogIntegrationId: "%s" }]`,syslogID)
	}
	
	customWebhookChannelConfigsQuery:=""
	if customWebhookURL!="" && len(customWebhookHeaders)==0{
		customWebhookChannelConfigsQuery=fmt.Sprintf(`customWebhookChannelConfigs: [
			{
				url: "%s"
				headers: []
				id: ""
			}
		]`,customWebhookURL)
	}else if customWebhookURL!="" && len(customWebhookHeaders)>0{
		headerString:="["
		for _, headers := range customWebhookHeaders {
			tmp := fmt.Sprintf(`{key: "%s", value: "%s", isSecret: %t}`, headers["key"], headers["value"],headers["is_secret"])
			headerString+=tmp
			headerString+=","
		}
		headerString=headerString[:len(headerString)-1]
		headerString+="]"

		customWebhookChannelConfigsQuery=fmt.Sprintf(`customWebhookChannelConfigs: [
			{
				url: "%s"
				headers: %s
				id: ""
			}
		]`,customWebhookURL,headerString)
	}

	emailChannelConfigsQuery:=""
	if len(emails)>0{
		emailString:="["
		for _,em :=range emails{
			tmp:=fmt.Sprintf(`{ address: "%s" }`,em)
			emailString+=tmp
			emailString+=","
		}
		emailString=emailString[:len(emailString)-1]
		emailString+="]"
		emailChannelConfigsQuery=fmt.Sprintf(`emailChannelConfigs: %s`,emailString)
	}

	if slackWebhookChannelConfigsQuery=="" && emailChannelConfigsQuery=="" && s3BucketChannelConfigsQuery=="" && splunkIntegrationChannelConfigsQuery=="" && syslogIntegrationChannelConfigsQuery=="" && customWebhookChannelConfigsQuery=="" && emailChannelConfigsQuery==""{
		return fmt.Errorf("No channel configuration provided")
	}

	query:=fmt.Sprintf(`mutation {
		updateNotificationChannel(
			input: {
				channelId: "%s"
				channelName: "%s"
				notificationChannelConfig: {
					%s
					%s
					%s
					%s
					%s
					%s
				}
			}
		) {
			channelId
		}
	}
	`,channel_id,channelName,slackWebhookChannelConfigsQuery,emailChannelConfigsQuery,s3BucketChannelConfigsQuery,splunkIntegrationChannelConfigsQuery,syslogIntegrationChannelConfigsQuery,customWebhookChannelConfigsQuery)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error:%s", err)
	}
	rules := response["data"].(map[string]interface{})["updateNotificationChannel"].(map[string]interface{})
	log.Println(rules)
	id:=rules["channelId"].(string)
	d.SetId(id)
 	return nil
}

func resourceNotificationChannelDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(`mutation {
		deleteNotificationChannel(
		  input: {channelId: "%s"}
		) {
		  success
		  
		}
	  }
	  `, id)
	_, err := executeQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	d.SetId("")
	return nil
}