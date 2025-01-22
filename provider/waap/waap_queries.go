package waap

const (
	ALL_ENV_CONFIG_SCOPE = `anomalyScope: { scopeType: CUSTOMER }`
	ENV_CONFIG_SCOPE     = `anomalyScope: {
                                                 scopeType: ENVIRONMENT
                                                 environmentScope: { id: "%s" }
                                             }`
	RULE_CONFIG = `ruleConfig: {
		ruleId: "%s"
		configType: %s
		configStatus: { disabled: %t }
	}`
	SUB_RULE_CONFIG = `ruleConfig: {
			ruleId: "%s"
			configType: %s
			subRuleConfigs: [{ subRuleId: "%s", anomalySubRuleAction: %s }]
		}`
	UPDATE_WAAP_CONFIG = `mutation {
                            updateAnomalyRuleConfig(
                                update: {
                                    %s
                                    %s
                                }
                            ) {
                                ruleId
                                configStatus { disabled }
                                subRuleConfigs{
                                            subRuleId
                                            blockingEnabled
                                        }
                                __typename
                            }
                        }`
	READ_QUERY = `{
                  anomalyDetectionRuleConfigs(
                    %s
                  ) {
                    count
                    total
                    results {
                      configStatus {
                        disabled
                        internal
                        __typename
                      }
                      hidden
                      eventFamily
                      configType
                      ruleCategory
                      ruleId
                      ruleName

                      subRuleConfigs {
					  anomalySubRuleAction
                        blockingEnabled
                        configStatus {
                          disabled
                          internal
                          __typename
                        }
                        subRuleId
                        subRuleName


                        anomalyProtectionType
                        anomalyRuleSeverity
                        __typename
                      }
                      anomalyRuleSeverity
                      __typename
                    }
                    __typename
                  }
                }`
)
