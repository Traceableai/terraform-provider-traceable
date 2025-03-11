package custom_signature

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"strings"
)

func EscapeString(input string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		line = strings.ReplaceAll(line, `\`, `\\`)
		line = strings.ReplaceAll(line, `"`, `\"`)
		lines[i] = line
	}
	return strings.Join(lines, `\n`)
}

func ReturnEnvScopedQuery(environments []interface{}) string {
	envQuery := ""
	if len(environments) != 0 {
		envQuery = fmt.Sprintf(ENVIRONMENT_SCOPE_QUERY, common.InterfaceToStringSlice(environments))
	}
	return envQuery
}

func ReturnReqResConditionsQuery(requestResponseSingleValuedConditions []interface{}, requestResponseMultiValuedConditions []interface{}) string {
	finalConditionsQuery := ""
	if len(requestResponseSingleValuedConditions) > 0 {
		for _, requestPayloadSingleValuedCondition := range requestResponseSingleValuedConditions {
			requestPayloadSingleValuedConditionData := requestPayloadSingleValuedCondition.(map[string]interface{})
			matchKey := requestPayloadSingleValuedConditionData["matchKey"].(string)
			matchCategory := requestPayloadSingleValuedConditionData["match_category"].(string)
			matchOp := requestPayloadSingleValuedConditionData["match_operator"].(string)
			matchVal := requestPayloadSingleValuedConditionData["match_value"].(string)
			finalConditionsQuery += fmt.Sprintf(SINGLE_VALUED_REQ_RES_CONDITION_QUERY, matchKey, matchCategory, matchOp, matchVal)
		}
	}
	if len(requestResponseMultiValuedConditions) > 0 {
		for _, requestPayloadMultiValuedCondition := range requestResponseMultiValuedConditions {
			requestPayloadMultiValuedConditionData := requestPayloadMultiValuedCondition.(map[string]interface{})
			matchCategory := requestPayloadMultiValuedConditionData["match_category"].(string)
			keyValueTag := requestPayloadMultiValuedConditionData["key_value_tag"].(string)
			keyMatchOperator := requestPayloadMultiValuedConditionData["key_match_operator"].(string)
			matchKey := requestPayloadMultiValuedConditionData["match_key"].(string)
			valueMatchOperator := requestPayloadMultiValuedConditionData["value_match_operator"].(string)
			matchValue := requestPayloadMultiValuedConditionData["match_value"].(string)
			finalConditionsQuery += fmt.Sprintf(MULTI_VALUED_REQ_RES_CONDITION_QUERY, keyValueTag, matchCategory, keyMatchOperator, matchKey, valueMatchOperator, matchValue)
		}
	}
	return finalConditionsQuery
}

func ReturnAttributeBasedConditionsQuery(attribute_based_conditions []interface{}) (string, error) {
	finalAttributeBasedConditionsQuery := ""
	for _, att_based_cond := range attribute_based_conditions {
		att_based_cond_data := att_based_cond.(map[string]interface{})
		if att_based_cond_data["value_condition_operator"].(string) == "" {
			finalAttributeBasedConditionsQuery += fmt.Sprintf(ATTRIBUTES_BASED_QUERY, att_based_cond_data["key_condition_operator"], att_based_cond_data["key_condition_value"], "")
		} else {
			if att_based_cond_data["value_condition_value"].(string) == "" {
				return "", fmt.Errorf("required both value_condition_operator and value_condition_value")
			} else {
				valueConditionString := fmt.Sprintf(ATTRIBUTE_VALUE_CONDITION_QUERY, att_based_cond_data["value_condition_operator"].(string), att_based_cond_data["value_condition_value"].(string))
				finalAttributeBasedConditionsQuery += fmt.Sprintf(ATTRIBUTES_BASED_QUERY, att_based_cond_data["key_condition_operator"], att_based_cond_data["key_condition_value"], valueConditionString)
			}
		}
	}
	return finalAttributeBasedConditionsQuery, nil
}

func ReturnCustomSecRuleQuery(custom_sec_rule string) string {
	customSecRuleQuery := ""
	if custom_sec_rule != "" {
		customSecRuleQuery = fmt.Sprintf(CUSTOM_SEC_RULE_QUERY, custom_sec_rule)
	}
	return customSecRuleQuery
}

func ReturnExipiryDuration(allow_expiry_duration string) string {
	exipiryDurationString := ""
	if allow_expiry_duration != "" {
		exipiryDurationString = fmt.Sprintf(`blockingExpirationDuration: "%s"`, allow_expiry_duration)
	}
	return exipiryDurationString
}

func ReturnfinalAgentEffectQuery(inject_request_headers []interface{}) string {
	finalAgentEffectQuery := ""

	for _, req_header := range inject_request_headers {
		req_header_key := req_header.(map[string]interface{})["header_key"]
		req_header_value := req_header.(map[string]interface{})["header_value"]
		finalAgentEffectQuery += fmt.Sprintf(AGENT_EFFECT_QUERY_TEMPLATE, req_header_key, req_header_value)
	}

	if finalAgentEffectQuery != "" {
		finalAgentEffectQuery = fmt.Sprintf(CUSTOM_HEADER_INJECTION_QUERY, finalAgentEffectQuery)
	}
	return finalAgentEffectQuery
}

func DeleteCustomSignatureRule(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
