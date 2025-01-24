package malicious_sources

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func DeleteIPRangeRule(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_QUERY_IP_RANGE, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId("")
	return nil
}
func DeleteRegionRule(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_REGION, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId("")
	return nil
}

func DeleteMaliciousSourcesRule(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_MALICIOUS_SOURCES, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId("")
	return nil
}

func ReturnEmailFraudScoreQuery(sev string) string {
	finalQuery := ""
	if sev != "" {
		finalQuery = fmt.Sprintf(EMAIL_FRAUD_SCORE_QUERY, sev)
	}
	return finalQuery
}

func ReturnExipiryDuration(exipiry string) string {
	exipiryDurationString := ""
	if exipiry != "" {
		exipiryDurationString = fmt.Sprintf(`expirationDuration: "%s"`, exipiry)
	}
	return exipiryDurationString
}
func ReturnMalicousSourcesExipiryDuration(exipiry string) string {
	exipiryDurationString := ""
	if exipiry != "" {
		exipiryDurationString = fmt.Sprintf(`expirationDetails: { expirationDuration: "%s" }`, exipiry)
	}
	return exipiryDurationString
}

func SetInjectedHeaders(ruleDetails map[string]interface{}) []map[string]interface{} {
	injectedHeaders := []map[string]interface{}{}
	if ruleEffect, ok := ruleDetails["effects"].(map[string]interface{}); ok {
		if agentEffect, ok := ruleEffect["agentEffect"].(map[string]interface{}); ok {
			if agentModifications, ok := agentEffect["agentModifications"].([]interface{}); ok {
				for _, agentModification := range agentModifications {
					agentModificationMap := agentModification.(map[string]interface{})
					injectedHeader := map[string]interface{}{
						"header_key":   agentModificationMap["headerInjection"].(map[string]interface{})["key"].(string),
						"header_value": agentModificationMap["headerInjection"].(map[string]interface{})["value"].(string),
					}
					injectedHeaders = append(injectedHeaders, injectedHeader)
				}
			}
		}
	}
	return injectedHeaders
}
