package dlp

const (
	DLP_KEY       = `DATA_EXFILTRATION`
	SERVICE_SCOPE = `{
                    leafCondition: {
                        conditionType: SCOPE
                        scopeCondition: {
                            scopeType: ENTITY
                            entityScope: {
                                entityIds: %s
                                entityType: SERVICE
                            }
                        }
                    }
                }`
	URL_REGEX_SCOPE = `{
                    leafCondition: {
                        conditionType: SCOPE
                        scopeCondition: {
                            scopeType: URL
                            urlScope: { urlRegexes: %s }
                        }
                    }
                }`
	DLP_VALUE_BASED_THRESHOLD_CONFIG_QUERY = `{
		apiAggregateType: %s
		userAggregateType: %s
		thresholdConfigType: VALUE_BASED
		valueBasedThresholdConfig: {
			valueType: SENSITIVE_PARAMS
			uniqueValuesAllowed: %d
			duration: "%s"
			sensitiveParamsEvaluationType: %s
		}
	}`

	MULTI_VALUES_REQ_CONDITIONS = `{
                    leafCondition: {
                        conditionType: KEY_VALUE
                        keyValueCondition: {
                            metadataType: %s
                            keyCondition: { operator: %s, value: "%s" }
							%s
                        }
                    }
                }`
	DATATYPE_VALUE_CONDITIONS            = `valueCondition: { operator: %s, value: "%s" }`
	DATATYPE_KEY_CONDITIONS              = `keyCondition: { operator: %s, value: "%s" }`
	ALERT_TRANSACTION_CONFIGS            = ` action: { actionType: %s, alert: { eventSeverity: %s } }`
	BLOCK_INDEFINITE_TRANSACTION_CONFIGS = `action: {
                    actionType: %s
                    block: { eventSeverity: %s }
                }`
	BLOCK_WITH_EXPIRY = `action: {
                    actionType: %s
                    block: { eventSeverity: %s, duration: "%s" }
                }`
	ALLOW_INDEFINITE_TRANSACTION_CONFIGS = `action: { actionType: %s }`
	ALLOW_WITH_EXPIRY                    = `action: { actionType: %s, allow: { duration: "%s" } }`
	DLP_REQUEST_BASED_QUERY_CREATE       = `mutation {
    createRateLimitingRule(
        rateLimitingRuleData: {
            category: DATA_EXFILTRATION
            conditions: [
               %s
            ]
            enabled: %t
            name: "%s"
            description: "%s"
            labels: []
            transactionActionConfigs: {
                %s
            }
            %s
        }
    ) {
        category
        id
        name
    }
}`
	DLP_REQUEST_BASED_QUERY_UPDATE = `mutation {
    updateRateLimitingRule(
        rateLimitingRule: {
			id: "%s"
            category: DATA_EXFILTRATION
            conditions: [
               %s
            ]
            enabled: %t
            name: "%s"
            description: "%s"
            labels: []
            transactionActionConfigs: {
                %s
            }
            %s
        }
    ) {
        category
        id
        name
    }
}`
	DLP_REQUEST_BASED_DATATYPE_CONDITIONS = `{
                    leafCondition: {
                        conditionType: DATATYPE
                        datatypeCondition: {
                            datasetIds: [%s]
                            datatypeIds: [%s]
                            dataLocation: REQUEST
                            datatypeMatching: {
                                datatypeMatchingType: REGEX_BASED_MATCHING
                                regexBasedMatching: {
                                    customMatchingLocation: {
                                        metadataType: %s
                                        %s
                                    }
                                }
                            }
                        }
                    }
                }`
)
