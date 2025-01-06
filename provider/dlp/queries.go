package dlp

const (
	DLP_KEY = `DATA_EXFILTRATION`
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
)