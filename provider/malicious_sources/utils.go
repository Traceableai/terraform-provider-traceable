package malicious_sources

import (
	"encoding/json"
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

func MapCountryNameToRegionId(regions []interface{}, meta interface{}) ([]interface{}, error) {
	regionResponseStr, fetchRegionQueryErr := common.CallExecuteQuery(FETCH_REGION_ID, meta)
	regionIds := []interface{}{}
	if fetchRegionQueryErr != nil {
		return regionIds, fetchRegionQueryErr
	}
	var response map[string]interface{}
	err := json.Unmarshal([]byte(regionResponseStr), &response)
	if err != nil {
		return regionIds, err
	}
	responseData, ok := response["data"].(map[string]interface{})
	if !ok {
		return regionIds, fmt.Errorf("unexpected response in countries graphql")
	}

	countries := responseData["countries"].(map[string]interface{})
	countiresList := countries["results"].([]interface{})
	countryMap := make(map[string]string)

	for _, country := range countiresList {
		if countryData, ok := country.(map[string]interface{}); ok {
			if cName, ok := countryData["name"].(string); ok {
				if cId, ok := countryData["id"].(string); ok {
					countryMap[cName] = cId
				}
			}
		}
	}
	for _, region := range regions {
		if id, exists := countryMap[region.(string)]; exists {
			regionIds = append(regionIds, id)
		} else {
			return []interface{}{}, fmt.Errorf("incorrect region name provided")
		}
	}
	return regionIds, nil
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
func ReturnEnvScopedQuery(environments []interface{}) string {
	envQuery := ""
	if len(environments) != 0 {
		envQuery = fmt.Sprintf(ENVIRONMENT_SCOPE_QUERY, common.InterfaceToStringSlice(environments))
	}
	return envQuery
}

func RegionRuleExpiryString(expiration string) string {
	exipiryDurationString := ""
	if expiration != "" {
		exipiryDurationString = fmt.Sprintf(`duration : "%s"`, expiration)
	}
	return exipiryDurationString
}
