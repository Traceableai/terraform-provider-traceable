package resources

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/traceableai/terraform-provider-traceable/internal/generated"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

var attackRules = map[string]map[string]string{
	"contentSize": {},

	"contentType": {},

	"httpStatus": {},

	"contentExplosion": {},

	"unexpectedUserAgent": {},

	"invalidEnumerations": {},

	"unknownParam": {},

	"missingParam": {},

	"specialCharacter": {},

	"typeAnomaly": {},

	"valueOutofRange": {},

	"ssrf": {},

	"bola": {},

	"userIdBola": {},

	"bfla": {},

	"volumetric": {
		"Spike in API Call Count": "volumetric_apiCallSpike",
	},

	"AccountTakeOver": {
		"Credential Stuffing": "ato_credStuff",
	},

	"sesionVoilation": {
		"Land Speed Violation": "sessionv_landspeed",
		"Session Expired":      "sessionv_expiry",
	},

	"jwt": {
		"JWT AlBeast Anomaly":         "jwt_albeast",
		"JWT Algorithm Confusion":     "jwt_alg",
		"JWT Expired":                 "jwt_exp",
		"JWT Invalid Signature":       "jwt_sign",
		"JWT JKU Integrity Violation": "jwt_jku",
		"JWT Premature Usage":         "jwt_nbf",
		"Unauthorized JWT audience":   "jwt_aud",
		"Unknown JWT Issuer":          "jwt_iss",
	},

	"XSS": {
		"IE XSS Filters - Attack (A href with a link)":                 "crs_941420",
		"JavaScript global variable":                                   "crs_941370",
		"JSFuck / Hieroglyphy obfuscation":                             "crs_941360",
		"UTF-7 Encoding IE XSS - Attack":                               "crs_941350",
		"US-ASCII Malformed Encoding XSS Filter - Attack":              "crs_941310",
		"IE XSS Filters - Attack (Object Tag)":                         "crs_941300",
		"IE XSS Filters - Attack (Applet Tag)":                         "crs_941290",
		"IE XSS Filters - Attack (Base href)":                          "crs_941280",
		"IE XSS Filters - Attack (Link href)":                          "crs_941270",
		"IE XSS Filters - Attack (Meta charset)":                       "crs_941260",
		"IE XSS Filters - Attack (Meta http-equiv)":                    "crs_941250",
		"IE XSS Filters - Attack (Import or Implementation Attribute)": "crs_941240",
		"IE XSS Filters - Attack (Embed Tag)":                          "crs_941230",
		"IE XSS Filters - Attack (Obfuscated VB Script)":               "crs_941220",
		"IE XSS Filters - Attack (Obfuscated Javascript)":              "crs_941210",
		"IE XSS Filters - Attack (VML Frames)":                         "crs_941200",
		"IE XSS Filters - Attack (Style Sheets)":                       "crs_941190",
		"Node-Validator Blacklist Keywords":                            "crs_941180",
		"NoScript XSS InjectionChecker: Attribute Injection (170)":     "crs_941170",
		"NoScript XSS InjectionChecker: HTML Injection":                "crs_941160",
		"XSS Filter - Category 4: Javascript URI Vector":               "crs_941140",
		"XSS Filter - Category 1: Script Tag Vector":                   "crs_941110",
		"XSS Attack (100)":                                             "crs_941100",
		"NoScript XSS InjectionChecker: Attribute Injection (T170)":    "crs_9410170",
		"XSS InjectionChecker: HTML Injection":                         "crs_9410160",
	},
	"HTTPProtocolAttack": {
		"LDAP Injection Attack":                                                      "crs_921200",
		"HTTP Splitting":                                                             "crs_921190",
		"HTTP Header Injection Attack via payload":                                   "crs_921160",
		"HTTP Header Injection Attack via headers":                                   "crs_921140",
		"HTTP Response Splitting Attack":                                             "crs_921130",
		"HTTP Response Splitting Attack (120) ":                                      "crs_921120",
		"Mail Injection: Protocol Manipulation":                                      "crs_9210400",
		"Suspicious URL Evasion Attempt: Possible Encoding or Obfuscation":           "crs_9210330",
		"Null Byte Injection":                                                        "crs_9210305",
		"URL Encoding Abuse Attack Attempt":                                          "crs_9210300",
		"Suspicious URL Evasion Attempt: Malformed URL Patterns":                     "crs_9210210",
		"Possible CSRF Attack: Host and Origin Mismatch":                             "crs_9210200",
		"HTTP Header Injection Attack via payload (CR/LF) (T150)":                    "crs_9210150",
		"HTTP Request Smuggling Attack":                                              "crs_9210120",
		"HTTP Request Smuggling Attack (Content-Length/Transfer-Encoding Confusion)": "crs_9210115",
		"HTTP Request Smuggling Attack (Request Splitting)":                          "crs_9210110",
		"Potential CORS Misconfiguration Exploitation Attempt":                       "crs_9210011",
		"IIS 6.0 WebDAV buffer overflow: (CVE-2017-7269)":                            "crs_9210010",
	},

	"LFI": {
		"Restricted File Access Attempt":                        "crs_930130",
		"OS File Access Attempt (120)":                          "crs_930120",
		"Path Traversal Attack (/../) (110)":                    "crs_930110",
		"Path Traversal Attack (/../) (100)":                    "crs_930100",
		"Possible Information Disclosure in ownCloud Graph API": "crs_9300150",
		"Rails Action View LFI: Accept header":                  "crs_9300010",
	},

	"RFI": {
		"Possible Remote File Inclusion (RFI) Attack: URL Payload Used w/Trailing Question Mark Character (?)": "crs_931120",
	},

	"NodeJsInjection": {
		"Node.js Injection Attack": "crs_932100",
	},

	"SQLInjection": {
		"MySQL in-line comment": "crs_942500",
		"Concatenated basic SQL injection and SQLLFI attempts (360)":                        "crs_942360",
		"MySQL UDF injection and other data/structure manipulation attempts":                "crs_942350",
		"MySQL and PostgreSQL stored procedure/function injections":                         "crs_942320",
		"Basic MongoDB SQL injection attempts":                                              "crs_942290",
		"Postgres pg_sleep injection, waitfor delay attacks and database shutdown attempts": "crs_942280",
		"Basic sql injection. Common attack string for mysql, oracle and others":            "crs_942270",
		"MATCH AGAINST, MERGE and EXECUTE IMMEDIATE injections":                             "crs_942250",
		"MySQL charset switch and MSSQL DoS attempts":                                       "crs_942240",
		"Conditional SQL injection attempts (230)":                                          "crs_942230",
		"SQL code execution and information gathering attempts":                             "crs_942190",
		"SQL benchmark and sleep injection attempts including conditional queries":          "crs_942170",
		"Blind sqli tests using sleep() or benchmark()":                                     "crs_942160",
		"SQL Injection Attack: Common DB Names":                                             "crs_942140",
		"SQL Injection Attack (100)":                                                        "crs_942100",
		"Concatenated basic SQL injection and SQLLFI attempts (T360)":                       "crs_9420360",
		"Basic MongoDB SQL injection attempts (Nested JSON)":                                "crs_9420291",
		"Basic MongoDB SQL injection attempts (Relaxed)":                                    "crs_9420290",
		"Conditional SQL injection attempts (T230)":                                         "crs_9420230",
		"DB code execution and information gathering attempts":                              "crs_9420190",
	},

	"XMLInjection": {
		"XML External Entity Injection: Entity Tags (110)":            "crs_1020110",
		"XML External Entity Injection: Local/Remote Includes (T100)": "crs_1020100",
	},

	"JavaAppAttack": {
		"Suspicious Java class":                                           "crs_944130",
		"Remote Command Execution: Java serialization":                    "crs_944120",
		"Remote Command Execution: Java process spawn":                    "crs_944110",
		"Remote Command Execution: Suspicious Java class":                 "crs_944100",
		"Apache Struts path traversal to RCE vulnerability":               "crs_9440960",
		"PAN-OS GlobalProtect Portal: RCE Exploitation":                   "crs_9440950",
		"Confluence Server OGNL Injection: Webwork":                       "crs_9440930",
		"Java Spring Core: RCE: Panic Rule":                               "crs_9440921",
		"Java Spring Core: RCE (CVE-2022-22965)":                          "crs_9440920",
		"Java Log4j: RCE DoS Exploitation (CVE-2021-45105) (T910)":        "crs_9440910",
		"Java Log4j: JNDI Exploitation (CVE-2021-44228) (T900)":           "crs_9440900",
		"Java Deserialization: Magic Bytes Base64 Encoded":                "crs_9440210",
		"Apache Commons Text: String Interpolation RCE: (CVE-2022-42889)": "crs_94400910",
	},

	"RCE": {
		"Restricted File Upload Attempt":                                                     "crs_932180",
		"Remote Command Execution: Shellshock (CVE-2014-6271) (171)":                         "crs_932171",
		"Remote Command Execution: Shellshock (CVE-2014-6271) (170)":                         "crs_932170",
		"Remote Command Execution: Unix Shell Code":                                          "crs_932160",
		"Remote Command Execution: Windows FOR/IF Command":                                   "crs_932140",
		"Remote Command Execution: Windows PowerShell Command":                               "crs_932120",
		"Remote Command Execution: Windows Command Injection (115)":                          "crs_932115",
		"Remote Command Execution: Windows Command Injection (110)":                          "crs_932110",
		"Remote Command Execution: Unix Command Injection (105)":                             "crs_932105",
		"Remote Command Execution: Unix Command Injection (100)":                             "crs_932100",
		"Server Side Template Injection (SSTI) Attempt":                                      "crs_9320310",
		"Remote Command Execution: Direct Unix Command Execution (T155)":                     "crs_9320155",
		"Remote Command Execution: Direct Unix Command Execution (T150)":                     "crs_9320150",
		"Remote Command Execution: Unix Command Injection (T105)":                            "crs_9320105",
		"Remote Command Execution: Unix Command Injection (T100)":                            "crs_9320100",
		"Remote Command Execution: Unix Command Injection (T200)":                            "crs_9320200",
		"Remote Command Execution: Command or Script Injection (T090)":                       "crs_9320090",
		"ServiceNow RCE: CVE-2024-4879 And CVE-2024-5217":                                    "crs_9320070",
		"Adobe ColdFusion RCE: Arbitrary File Upload Vulnerability: (CVE-2018-15961)":        "crs_9320060",
		"vBulletin RCE : (CVE-2019-16759)":                                                   "crs_9320050",
		"Confluence Server and Data Center Path Traversal to RCE : (CVE-2019-3398)":          "crs_9320041",
		"Atlassian Confluence SSTI RCE : (CVE-2019-3396)":                                    "crs_9320040",
		"Apache Solr DataImport Handler RCE : (CVE-2019-0193)":                               "crs_9320031",
		"Apache Solr Deserialization RCE : (CVE-2019-0192)":                                  "crs_9320030",
		"Apache Struts RCE: OGNL in URL: CVE-2018-11776":                                     "crs_9320021",
		"Apache Struts RCE: multipart/form-data: (CVE-2017-5638)":                            "crs_9320020",
		"Drupal Core: RCE SA-CORE-2019-003: (CVE-2019-6340)":                                 "crs_93200120",
		"Microsoft Exchange Server: RCE: (CVE-2021-26855)":                                   "crs_93200110",
		"Atlassian Crowd: RCE: (CVE-2019-11580)":                                             "crs_93200100",
		"Spring Cloud Function RCE: malicious Spring Expression: (CVE-2022-22963)":           "crs_9320010",
		"Zoho ManageEngine ServiceDesk Plus: Arbitrary File Upload: (CVE-2019-8394)":         "crs_93200090",
		"Oracle WebLogic Server: RCE: (CVE-2019-2725)":                                       "crs_93200080",
		"Blueimp jQuery-File-Upload: Unauthenticated arbitrary file upload: (CVE-2018-9206)": "crs_93200070",
		"Remote Command Execution:: Common Evasions for Shell Injection":                     "crs_9320510",
		"Remote Command Execution: Windows Command Injection (T110)":                         "crs_9320110",
	},

	"SessionFixation": {
		"Possible Session Fixation Attack: SessionID Parameter Name with No Referer":         "crs_943120",
		"Possible Session Fixation Attack: SessionID Parameter Name with Off-Domain Referer": "crs_943110",
		"Possible Session Fixation Attack: Setting Cookie Values in HTML":                    "crs_943100",
	},

	"SSRF": {
		"Server Side Request Forgery (SSRF): Common Collaborators":             "crs_1010110",
		"Server Side Request Forgery (SSRF): Cloud Provider Metadata Endpoint": "crs_1010100",
	},

	"BasicAuthenticationViolation": {
		"Basic Authentication: Missing Password":     "crs_1030100",
		"Possible Joomla Unauthenticated Acces":      "crs_1030110",
		"Authorization Bypass in Next.js Middleware": "crs_1030120",
	},

	"ScannerDetection": {
		"Request filename/argument associated with security scanner": "crs_913120",
		"Request header associated with security scanner":            "crs_913110",
		"User-Agent associated with security scanner":                "crs_913100",
	},

	"GraphQLAttacks": {
		"GraphQL Introspection Query Detected": "crs_1040100",
	},
}

var CustomSignatureKeyValuesExpressionMap = map[string]generated.CustomSignatureRuleKeyValueTag{
	"HEADER":    generated.CustomSignatureRuleKeyValueTagHeader,
	"PARAMETER": generated.CustomSignatureRuleKeyValueTagParameter,
	"COOKIE":    generated.CustomSignatureRuleKeyValueTagCookie,
}
var WaapConfigSubRuleActionMap = map[string]generated.AnomalySubRuleAction{
	"IGNORE":  generated.AnomalySubRuleActionIgnore,
	"DISABLE": generated.AnomalySubRuleActionDisable,
	"MONITOR": generated.AnomalySubRuleActionMonitor,
	"BLOCK":   generated.AnomalySubRuleActionBlock,
}

var CustomSignatureRuleMatchCategoryMap = map[string]generated.CustomSignatureRuleMatchCategory{
	"REQUEST":  generated.CustomSignatureRuleMatchCategoryRequest,
	"RESPONSE": generated.CustomSignatureRuleMatchCategoryResponse,
}

var CustomSignatureRuleEventTypeMap = map[string]generated.CustomSignatureRuleEventType{
	"TESTING_DETECTION":      generated.CustomSignatureRuleEventTypeTestingDetection,
	"NORMAL_DETECTION":       generated.CustomSignatureRuleEventTypeNormalDetection,
	"DETECTION_AND_BLOCKING": generated.CustomSignatureRuleEventTypeDetectionAndBlocking,
	"ALLOW":                  generated.CustomSignatureRuleEventTypeAllow,
}

var CustomSignatureRuleMatchKeyMap = map[string]generated.CustomSignatureRuleMatchKey{
	"URL":                generated.CustomSignatureRuleMatchKeyUrl,
	"HEADER_NAME":        generated.CustomSignatureRuleMatchKeyHeaderName,
	"HEADER_VALUE":       generated.CustomSignatureRuleMatchKeyHeaderValue,
	"PARAMETER_NAME":     generated.CustomSignatureRuleMatchKeyParameterName,
	"PARAMETER_VALUE":    generated.CustomSignatureRuleMatchKeyParameterValue,
	"HTTP_METHOD":        generated.CustomSignatureRuleMatchKeyHttpMethod,
	"HOST":               generated.CustomSignatureRuleMatchKeyHost,
	"USER_AGENT":         generated.CustomSignatureRuleMatchKeyUserAgent,
	"STATUS_CODE":        generated.CustomSignatureRuleMatchKeyStatusCode,
	"BODY":               generated.CustomSignatureRuleMatchKeyBody,
	"BODY_SIZE":          generated.CustomSignatureRuleMatchKeyBodySize,
	"COOKIE_NAME":        generated.CustomSignatureRuleMatchKeyCookieName,
	"COOKIE_VALUE":       generated.CustomSignatureRuleMatchKeyCookieValue,
	"QUERY_PARAMS_COUNT": generated.CustomSignatureRuleMatchKeyQueryParamsCount,
	"HEADERS_COUNT":      generated.CustomSignatureRuleMatchKeyHeadersCount,
	"COOKIES_COUNT":      generated.CustomSignatureRuleMatchKeyCookiesCount,
}

var RateLimitingRuleEventSeverityMap = map[string]generated.RateLimitingRuleEventSeverity{
	"LOW":      generated.RateLimitingRuleEventSeverityLow,
	"MEDIUM":   generated.RateLimitingRuleEventSeverityMedium,
	"HIGH":     generated.RateLimitingRuleEventSeverityHigh,
	"CRITICAL": generated.RateLimitingRuleEventSeverityCritical,
}

var RateLimitingApiAggregateMap = map[string]generated.RateLimitingRuleApiAggregateType{
	"PER_ENDPOINT":    generated.RateLimitingRuleApiAggregateTypePerEndpoint,
	"ACROSS_ENDPOINT": generated.RateLimitingRuleApiAggregateTypeAcrossEndpoints,
}

var RateLimitingUserAggregateMap = map[string]generated.RateLimitingRuleUserAggregateType{
	"PER_USER":    generated.RateLimitingRuleUserAggregateTypePerUser,
	"ACROSS_USER": generated.RateLimitingRuleUserAggregateTypeAcrossUsers,
}

var valueBasedThresholdConfigMap = map[string]generated.ValueBasedThresholdConfigType{
	"REQUEST_BODY":     generated.ValueBasedThresholdConfigTypeRequestBody,
	"SENSITIVE_PARAMS": generated.ValueBasedThresholdConfigTypeSensitiveParams,
	"PATH_PARAMS":      generated.ValueBasedThresholdConfigTypePathParams,
}

var RateLimitingActionMap = map[string]generated.RateLimitingRuleActionType{
	"ALERT":            generated.RateLimitingRuleActionTypeAlert,
	"BLOCK":            generated.RateLimitingRuleActionTypeBlock,
	"ALLOW":            generated.RateLimitingRuleActionTypeAllow,
	"MARK_FOR_TESTING": generated.RateLimitingRuleActionTypeMarkForTesting,
}

var RateLimitingRuleThresholdConfigMap = map[string]generated.RateLimitingRuleThresholdConfigType{
	"ROLLING_WINDOW": generated.RateLimitingRuleThresholdConfigTypeRollingWindow,
	"VALUE_BASED":    generated.RateLimitingRuleThresholdConfigTypeValueBased,
	"DYNAMIC":        generated.RateLimitingRuleThresholdConfigTypeDynamic,
}

var RateLimitingRuleConditionMap = map[string]generated.RateLimitingRuleConditionType{
	"LOGICAL_CONDITION": generated.RateLimitingRuleConditionTypeLogicalCondition,
	"LEAF_CONDITION":    generated.RateLimitingRuleConditionTypeLeafCondition,
}

var RateLimitingRuleScannerMap = map[string]bool{
	"Traceable AST":        true,
	"Qualys":               true,
	"Rapid7 InsightAppSec": true,
	"Invicti":              true,
	"Tenable":              true,
}

var RateLimitingRuleIpConnectionTypeMap = map[string]generated.RateLimitingRuleIpConnectionType{
	"RESIDENTIAL": generated.RateLimitingRuleIpConnectionTypeResidential,
	"MOBILE":      generated.RateLimitingRuleIpConnectionTypeMobile,
	"CORPORATE":   generated.RateLimitingRuleIpConnectionTypeCorporate,
	"DATA_CENTER": generated.RateLimitingRuleIpConnectionTypeDataCenter,
	"EDUCATION":   generated.RateLimitingRuleIpConnectionTypeEducation,
}

var RateLimitingKeyValueMatchOperatorMap = map[string]generated.RateLimitingRuleKeyValueMatchOperator{
	"EQUALS":          generated.RateLimitingRuleKeyValueMatchOperatorEquals,
	"NOT_EQUAL":       generated.RateLimitingRuleKeyValueMatchOperatorNotEqual,
	"MATCHES_REGEX":   generated.RateLimitingRuleKeyValueMatchOperatorMatchesRegex,
	"NOT_MATCH_REGEX": generated.RateLimitingRuleKeyValueMatchOperatorNotMatchRegex,
	"CONTAINS":        generated.RateLimitingRuleKeyValueMatchOperatorContains,
	"NOT_CONTAIN":     generated.RateLimitingRuleKeyValueMatchOperatorNotContain,
	"GREATER_THAN":    generated.RateLimitingRuleKeyValueMatchOperatorGreaterThan,
	"LESS_THAN":       generated.RateLimitingRuleKeyValueMatchOperatorLessThan,
}

var RateLimitingRuleIpReputationSeverityMap = map[string]generated.RateLimitingRuleIpReputationSeverity{
	"LOW":      generated.RateLimitingRuleIpReputationSeverityLow,
	"MEDIUM":   generated.RateLimitingRuleIpReputationSeverityMedium,
	"HIGH":     generated.RateLimitingRuleIpReputationSeverityHigh,
	"CRITICAL": generated.RateLimitingRuleIpReputationSeverityCritical,
}

var RateLimitingRuleIpLocationTypeMap = map[string]generated.RateLimitingRuleIpLocationType{
	"RESIDENTIAL":      generated.RateLimitingRuleIpLocationTypeResidential,
	"ANONYMOUS":        generated.RateLimitingRuleIpLocationTypeAnonymous,
	"ANONYMOUS_VPN":    generated.RateLimitingRuleIpLocationTypeAnonymousVpn,
	"HOSTING_PROVIDER": generated.RateLimitingRuleIpLocationTypeHostingProvider,
	"PUBLIC_PROXY":     generated.RateLimitingRuleIpLocationTypePublicProxy,
	"TOR_EXIT_NODE":    generated.RateLimitingRuleIpLocationTypeTorExitNode,
	"BOT":              generated.RateLimitingRuleIpLocationTypeBot,
	"SCANNER":          generated.RateLimitingRuleIpLocationTypeScanner,
}

var RateLimitingRuleIpAbuseVelocityMap = map[string]generated.RateLimitingRuleIpAbuseVelocity{
	"LOW":    generated.RateLimitingRuleIpAbuseVelocityLow,
	"MEDIUM": generated.RateLimitingRuleIpAbuseVelocityMedium,
	"HIGH":   generated.RateLimitingRuleIpAbuseVelocityHigh,
}
var RateLimitingRuleKeyValueConditionMetadataTypeMap = map[string]generated.RateLimitingRuleKeyValueConditionMetadataType{
	"URL":                     generated.RateLimitingRuleKeyValueConditionMetadataTypeUrl,
	"HOST":                    generated.RateLimitingRuleKeyValueConditionMetadataTypeHost,
	"HTTP_METHOD":             generated.RateLimitingRuleKeyValueConditionMetadataTypeHttpMethod,
	"USER_AGENT":              generated.RateLimitingRuleKeyValueConditionMetadataTypeUserAgent,
	"QUERY_PARAMETER":         generated.RateLimitingRuleKeyValueConditionMetadataTypeQueryParameter,
	"STATUS_CODE":             generated.RateLimitingRuleKeyValueConditionMetadataTypeStatusCode,
	"REQUEST_BODY":            generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestBody,
	"RESPONSE_BODY":           generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseBody,
	"REQUEST_HEADER":          generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestHeader,
	"RESPONSE_HEADER":         generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseHeader,
	"REQUEST_COOKIE":          generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestCookie,
	"RESPONSE_COOKIE":         generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseCookie,
	"REQUEST_BODY_PARAMETER":  generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestBodyParameter,
	"RESPONSE_BODY_PARAMETER": generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseBodyParameter,
	"TAG":                     generated.RateLimitingRuleKeyValueConditionMetadataTypeTag,
	"REQUEST_BODY_SIZE":       generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestBodySize,
	"RESPONSE_BODY_SIZE":      generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseBodySize,
	"QUERY_PARAMS_COUNT":      generated.RateLimitingRuleKeyValueConditionMetadataTypeQueryParamsCount,
	"REQUEST_HEADERS_COUNT":   generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestHeadersCount,
	"RESPONSE_HEADERS_COUNT":  generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseHeadersCount,
	"REQUEST_COOKIES_COUNT":   generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestCookiesCount,
	"RESPONSE_COOKIES_COUNT":  generated.RateLimitingRuleKeyValueConditionMetadataTypeResponseCookiesCount,
}

var MaliciousIpRangeEventSeverityMap = map[string]generated.IpRangeEventSeverity{
	"LOW":      generated.IpRangeEventSeverityLow,
	"MEDIUM":   generated.IpRangeEventSeverityMedium,
	"HIGH":     generated.IpRangeEventSeverityHigh,
	"CRITICAL": generated.IpRangeEventSeverityCritical,
}

var MaliciousIpRangeActionMap = map[string]generated.IpRangeRuleActionType{
	"BLOCK":            generated.IpRangeRuleActionTypeRuleActionBlock,
	"ALLOW":            generated.IpRangeRuleActionTypeRuleActionAllow,
	"BLOCK_ALL_EXCEPT": generated.IpRangeRuleActionTypeRuleActionBlockAllExcept,
	"ALERT":            generated.IpRangeRuleActionTypeRuleActionAlert,
}
var MaliciousIpRangeIpRangeMapResponse = map[string]string{
	"RULE_ACTION_BLOCK":            "BLOCK",
	"RULE_ACTION_ALLOW":            "ALLOW",
	"RULE_ACTION_BLOCK_ALL_EXCEPT": "BLOCK_ALL_EXCEPT",
	"RULE_ACTION_ALERT":            "ALERT",
}

var MaliciousRegionEventSeverityMap = map[string]generated.RegionRuleEventSeverity{
	"LOW":      generated.RegionRuleEventSeverityLow,
	"MEDIUM":   generated.RegionRuleEventSeverityMedium,
	"HIGH":     generated.RegionRuleEventSeverityHigh,
	"CRITICAL": generated.RegionRuleEventSeverityCritical,
}

var MaliciousRegionActionMap = map[string]generated.RegionRuleActionType{
	"BLOCK":            generated.RegionRuleActionTypeBlock,
	"BLOCK_ALL_EXCEPT": generated.RegionRuleActionTypeBlockAllExcept,
	"ALERT":            generated.RegionRuleActionTypeAlert,
}

var MaliciousIpTypeEventSeverityMap = map[string]generated.MaliciousSourcesRuleEventSeverity{
	"LOW":      generated.MaliciousSourcesRuleEventSeverityLow,
	"MEDIUM":   generated.MaliciousSourcesRuleEventSeverityMedium,
	"HIGH":     generated.MaliciousSourcesRuleEventSeverityHigh,
	"CRITICAL": generated.MaliciousSourcesRuleEventSeverityCritical,
}

var MaliciousIpTypeActionMap = map[string]generated.MaliciousSourcesRuleActionType{
	"BLOCK": generated.MaliciousSourcesRuleActionTypeBlock,
	"ALERT": generated.MaliciousSourcesRuleActionTypeAlert,
}

var MaliciousIpTypeMap = map[string]generated.MaliciousSourcesRuleIpLocationType{
	"ANONYMOUS_VPN":    generated.MaliciousSourcesRuleIpLocationTypeAnonymousVpn,
	"HOSTING_PROVIDER": generated.MaliciousSourcesRuleIpLocationTypeHostingProvider,
	"PUBLIC_PROXY":     generated.MaliciousSourcesRuleIpLocationTypePublicProxy,
	"TOR_EXIT_NODE":    generated.MaliciousSourcesRuleIpLocationTypeTorExitNode,
	"BOT":              generated.MaliciousSourcesRuleIpLocationTypeBot,
}

var RateLimitingDataLocationMap = map[string]generated.RateLimitingRuleDataLocation{
	"REQUEST":  generated.RateLimitingRuleDataLocationRequest,
	"RESPONSE": generated.RateLimitingRuleDataLocationResponse,
}
var RateLimitingValueBasedThresholdConfigTypeMap = map[string]generated.ValueBasedThresholdConfigType{
	"REQUEST_BODY":     generated.ValueBasedThresholdConfigTypeRequestBody,
	"SENSITIVE_PARAMS": generated.ValueBasedThresholdConfigTypeSensitiveParams,
	"PATH_PARAMS":      generated.ValueBasedThresholdConfigTypePathParams,
}
var RateLimitingSensitiveParamsEvaluationTypeMap = map[string]generated.SensitiveParamsEvaluationType{
	"ALL":                 generated.SensitiveParamsEvaluationTypeAll,
	"SELECTED_DATA_TYPES": generated.SensitiveParamsEvaluationTypeSelectedDataTypes,
}

var MaliciousEmailDomainMinEmailFraudScoreLevel = map[string]generated.MaliciousSourcesRuleEmailFraudScoreLevel{
	"HIGH":     generated.MaliciousSourcesRuleEmailFraudScoreLevelHigh,
	"CRITICAL": generated.MaliciousSourcesRuleEmailFraudScoreLevelCritical,
}

// var RateLimitingRequestResponseMultipleMap = map[string]bool{
// 	"QUERY_PARAMETER":         true,
// 	"REQUEST_BODY_PARAMETER":  true,
// 	"RESPONSE_BODY_PARAMETER": true,
// 	"REQUEST_COOKIE":          true,
// 	"RESPONSE_COOKIE":         true,
// 	"REQUEST_HEADERS":         true,
// 	"RESPONSE_HEADERS":        true,
// }

// HasValue checks if a field has a concrete value
// if string not empty
func HasValue(field interface{}) bool {

	if field == nil {
		return false
	}

	val := reflect.ValueOf(field)
	if (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && val.IsNil() {
		return false
	}

	switch v := field.(type) {
	case attr.Value:
		if v.IsNull() || v.IsUnknown() {
			return false
		}

		// Handle specific framework types
		switch concrete := v.(type) {
		case types.String:
			return concrete.ValueString() != ""
		case types.List, types.Set:
			elements := reflect.ValueOf(concrete).MethodByName("Elements").Call(nil)[0]
			for i := 0; i < elements.Len(); i++ {
				element := elements.Index(i).Interface().(attr.Value)
				if HasValue(element) {
					return true
				}
			}
			return false
		case types.Map:
			for _, v := range concrete.Elements() {
				if HasValue(v) {
					return true
				}
			}
			return false

		case types.Object:
			attrs := concrete.Attributes()
			if len(attrs) == 0 {
				return false
			}
			for _, attrVal := range attrs {
				if HasValue(attrVal) {
					return true
				}
			}
			return false
		default:
			return true
		}
	default:
		return true
	}

}

func convertToRuleConfigScope(environments types.Set) (*generated.InputRuleConfigScope, error) {
	if !HasValue(environments) {
		return nil, nil
	}

	var scope *generated.InputRuleConfigScope
	var envIds []*string
	for _, env := range environments.Elements() {
		if env, ok := env.(types.String); ok {
			env1 := env.ValueString()
			envIds = append(envIds, &env1)
		}

	}
	scope = &generated.InputRuleConfigScope{
		EnvironmentScope: &generated.InputEnvironmentScope{
			EnvironmentIds: envIds,
		},
	}
	return scope, nil
}

func GetCountriesId(isoCodes []*string, ctx context.Context, r graphql.Client) ([]*string, error) {
	countriesId := []*string{}
	response, err := generated.GetCountries(ctx, r)
	countriesPresent := map[string]bool{}
	for _, isoCode := range isoCodes {
		countriesPresent[*isoCode] = true
	}
	if err != nil {
		return nil, err
	}
	for _, country := range response.GetCountries().Results {
		if countriesPresent[country.Country.IsoCode] {
			countriesId = append(countriesId, &country.Id)
			countriesPresent[country.Country.IsoCode] = false
		}
	}
	for _, isoCode := range isoCodes {
		if countriesPresent[*isoCode] {
			return nil, utils.NewInvalidError("regions", fmt.Sprintf("%s is not a supported ISO Code", *isoCode))
		}
	}
	return countriesId, nil
}

func GetConfigType(ruleName string) (string, error) {
	ruleCrsId, err := GetRuleId(ruleName)
	if err != nil {
		return "", err
	}
	if strings.Contains(ruleCrsId, "crs") {
		return string(generated.AnomalyDetectionConfigTypeModsecurity), nil
	}
	if strings.Contains(ruleCrsId, "sessionv") || strings.Contains(ruleCrsId, "bola") || strings.Contains(ruleCrsId, "userIdBola") {
		return string(generated.AnomalyDetectionConfigTypeSessionDefinitionMetadata), nil
	}
	if strings.Contains(ruleCrsId, "ato") {
		return string(generated.AnomalyDetectionConfigTypeAccountTakeover), nil
	}
	if strings.Contains(ruleCrsId, "volumetric") {
		return string(generated.AnomalyDetectionConfigTypeVolumetric), nil
	}
	return string(generated.AnomalyDetectionConfigTypeApiDefinitionMetadata), nil
}

func GetRuleId(ruleName string) (string, error) {
	subrules, exists := attackRules[ruleName]
	if !exists {
		return "", utils.NewInvalidError("rule_configs rule_name", fmt.Sprintf("policy %s does not exist", ruleName))
	}
	if len(subrules) == 0 {
		return ruleName, nil
	}
	for _, crsID := range subrules {
		parts := strings.Split(crsID, "_")
		if len(parts) == 2 && len(parts[1]) >= 3 {
			return "crs_" + parts[1][:3], nil
		}
	}
	return ruleName, nil
}

func GetSubRuleId(mainRuleName string, subRuleName string) (string, error) {
	fmt.Printf("This is the input %s %s", mainRuleName, subRuleName)
	subrules, exists := attackRules[mainRuleName]
	if !exists {
		return "", utils.NewInvalidError("rule_configs rule_name", fmt.Sprintf("policy %s does not exist", mainRuleName))
	}
	if len(subrules) > 0 {
		for expectedSubRuleName, subRuleCrsID := range subrules {
			if expectedSubRuleName == subRuleName {
				return subRuleCrsID, nil
			}
		}
	}
	return "", utils.NewInvalidError("rule_configs subrules", fmt.Sprintf("subrule %s does not exist for rule_name %s", subRuleName, mainRuleName))
}