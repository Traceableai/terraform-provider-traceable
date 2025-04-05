package resources

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/generated"
)

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

////////////////

// // HasValue checks if any Terraform attribute value contains actual data
// func HasValue(val interface{}) (bool    ) {
//     // Handle nil case
//     if val == nil {
//         return false
//     }

//     // Type switch for all supported attribute types
//     switch v := val.(type) {
//     case attr.Value:
//         return hasAttrValue(v)
//     case *types.Object:
//         if v == nil {
//             return false
//         }
//         return hasAttrValue(*v)
//     case *types.List:
//         if v == nil {
//             return false
//         }
//         return hasAttrValue(*v)
//     case *types.Map:
//         if v == nil {
//             return false
//         }
//         return hasAttrValue(*v)
//     case *types.Set:
//         if v == nil {
//             return false
//         }
//         return hasAttrValue(*v)
//     default:
//         return false
//     }
// }

// // hasAttrValue handles the actual attribute value checking
// func hasAttrValue(val attr.Value) (bool) {

//     if val == nil || val.IsNull() || val.IsUnknown() {
//         return false
//     }

//     switch v := val.(type) {
//     case types.Object:
//         return hasObjectValue(v)
//     case types.List:
//         return hasCollectionValue(v)
//     case types.Map:
//         return hasCollectionValue(v)
//     case types.Set:
//         return hasCollectionValue(v)
//     default:
//         return true
//     }
// }

// // hasObjectValue checks if an Object has any non-null values
// func hasObjectValue(obj types.Object) (bool) {
//     var objMap map[string]attr.Value

//     obj.As(&objMap, basetypes.ObjectAsOptions{})

//     for _, val := range objMap {
//         hasVal := HasValue(val)
//         if hasVal {
//             return true
//         }
//     }

//     return false
// }

// // hasCollectionValue checks if collections (List/Map/Set) have any non-null values
// func hasCollectionValue(coll attr.Value) (bool) {
//     var elements []attr.Value

//     switch v := coll.(type) {
//     case types.List:
//         v.ElementsAs(&elements, false)
//     case types.Set:
//         v.ElementsAs(&elements, false)
//     case types.Map:
//         v.ElementsAs(&elements, false)
//     default:
//         return false
//     }

//     for _, elem := range elements {
//         hasVal := HasValue(elem)
//         if hasVal {
//             return false
//         }
//     }

//     return true
// }
