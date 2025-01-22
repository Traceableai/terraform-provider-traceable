package waap

import (
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/provider/notification"
	"strings"
)

func GetConfigScope(environment string) string {
	if environment != "" {
		return fmt.Sprintf(ENV_CONFIG_SCOPE, environment)
	}
	return ALL_ENV_CONFIG_SCOPE
}

func GetConfigType(ruleId string) (string, error) {
	ruleCrsId, err := GetRuleCrsId(ruleId)
	if err != nil {
		return "", err
	}
	if strings.Contains(ruleCrsId, "crs") {
		return "MODSECURITY", nil
	}
	return "API_DEFINITION_METADATA", nil
}

func GetRuleConfig(ruleId string, configType string, disabled bool, subRuleId string, subRuleAction string) (string, error) {
	ruleCrsId, err := GetRuleCrsId(ruleId)
	if err != nil {
		return "", err
	}
	ruleConfigs := fmt.Sprintf(RULE_CONFIG, ruleCrsId, configType, disabled)
	if subRuleId != "" {
		ruleConfigs = fmt.Sprintf(SUB_RULE_CONFIG, ruleCrsId, configType, subRuleId, subRuleAction)
	}
	return ruleConfigs, nil
}

func GetRuleCrsId(ruleId string) (string, error) {
	exists, ruleCrsId := notification.IsPreDefinedThreatEvent(ruleId)
	if !exists {
		return "", fmt.Errorf("no rule found with supplied id")
	}
	return ruleCrsId, nil
}
