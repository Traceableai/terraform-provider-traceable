package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiAbuseRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiAbuseCreate,
		Read:   resourceApiAbuseRead,
		Update: resourceApiAbuseUpdate,
		Delete: resourceApiAbuseDelete,

		Schema: map[string]*schema.Schema{
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Description: "category of the policy",
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Description: "name of the policy",
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Description: "discription of the policy",
				Optional: true,
			},
			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Description: "env of the policy",
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeString,
				Description: "true or false",
				Optional: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeMap,
				Description: "event severity",
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
							Type:        schema.TypeString,
							Required:    true,
						},
						"severity": {
							Type:        schema.TypeString,
							Required:    true,
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"conditions": &schema.Schema{
				Type:     schema.TypeList,
				Description: "event severity",
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeMap,
							Description: "event severity",
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scope": {
										Type:        schema.TypeString,
										Required:    true,
										Elem: 
									},
									"severity": {
										Type:        schema.TypeString,
										Required:    true,
									},
									"expiration_time": {
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"severity": {
							Type:        schema.TypeString,
							Required:    true,
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}
