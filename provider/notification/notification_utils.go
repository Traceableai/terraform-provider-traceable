package notification

func IsCustomThreatEvent(event_type string) bool {
	customThreatActivities := map[string]struct{}{
		"RATE_LIMITING":              {},
		"DATA_ACCESS":                {},
		"ENUMERATION":                {},
		"MALICIOUS_SOURCES_IP_TYPE":  {},
		"CUSTOM_SIGNATURE":           {},
		"MALICIOUS_SOURCES_REGION":   {},
		"MALICIOUS_SOURCES_EMAIL":    {},
		"MALICIOUS_SOURCES_IP_RANGE": {},
	}
	if _, exists := customThreatActivities[event_type]; exists {
		return true
	}
	return false
}

func FindThreatByCrsId(crsid string) string {
	crsToThreatType := map[string]string{
		"bola":                   "bola",
		"userIdBola":             "userIdBola",
		"bfla":                   "bfla",
		"sessionv":               "sesionVoilation",
		"volumetricApiCallSpike": "volumetricApiCallSpike",
		"credentialStuffing":     "credentialStuffing",
		"contentSize":            "contentSize",
		"contentType":            "contentType",
		"httpStatus":             "httpStatus",
		"contentExplosion":       "contentExplosion",
		"device":                 "unexpectedUserAgent",
		"enum":                   "invalidEnumerations",
		"unknownParam":           "unknownParam",
		"missingParam":           "missingParam",
		"specialCharacter":       "specialCharacter",
		"type":                   "typeAnomaly",
		"integer":                "valueOutofRange",
		"crs_941":                "XSS",
		"crs_930":                "LFI",
		"crs_931":                "RFI",
		"crs_921":                "HTTPProtocolAttack",
		"crs_934":                "NodeJsInjection",
		"crs_942":                "SQLInjection",
		"crs_102":                "XMLInjection",
		"crs_944":                "JavaAppAttack",
		"crs_932":                "RCE",
		"crs_943":                "SessionFixation",
		"crs_101":                "SSRF",
		"ssrf":                   "ssrf",
		"crs_103":                "BasicAuthenticationViolation",
		"jwt":                    "jwt",
		"crs_913":                "Scanner Detection",
		"crs_104":                "GraphQLAttacks",
	}
	return crsToThreatType[crsid]
}
func IsPreDefinedThreatEvent(event_type string) (bool, string) {
	preDefinedThreatActivities := map[string]string{
		"bola":                         "bola",
		"userIdBola":                   "userIdBola",
		"bfla":                         "bfla",
		"sesionVoilation":              "sessionv",
		"volumetricApiCallSpike":       "volumetricApiCallSpike",
		"credentialStuffing":           "credentialStuffing",
		"contentSize":                  "contentSize",
		"contentType":                  "contentType",
		"httpStatus":                   "httpStatus",
		"contentExplosion":             "contentExplosion",
		"unexpectedUserAgent":          "device",
		"invalidEnumerations":          "enum",
		"unknownParam":                 "unknownParam",
		"missingParam":                 "missingParam",
		"specialCharacter":             "specialCharacter",
		"typeAnomaly":                  "type",
		"valueOutofRange":              "integer",
		"XSS":                          "crs_941",
		"LFI":                          "crs_930",
		"RFI":                          "crs_931",
		"HTTPProtocolAttack":           "crs_921",
		"NodeJsInjection":              "crs_934",
		"SQLInjection":                 "crs_942",
		"XMLInjection":                 "crs_102",
		"JavaAppAttack":                "crs_944",
		"RCE":                          "crs_932",
		"SessionFixation":              "crs_943",
		"SSRF":                         "crs_101",
		"ssrf":                         "ssrf",
		"BasicAuthenticationViolation": "crs_103",
		"jwt":                          "jwt",
		"Scanner Detection":            "crs_913",
		"GraphQLAttacks":               "crs_104",
	}
	if _, exists := preDefinedThreatActivities[event_type]; exists {
		return true, preDefinedThreatActivities[event_type]
	}
	return false, ""
}
