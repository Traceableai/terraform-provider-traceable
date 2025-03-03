package testdata

var SeverityOptions=[]string{
	"LOW","MEDIUM","HIGH","CRITICAL",
}

var RequestResponsePairs=[]string{
	"REQUEST","RESPONSE",
}
var ExcludeOptions=[]bool{
	true,false,
}
var IPReputationOptions=[]string{
  "LOW","MEDIUM","HIGH","CRITICAL",
}

var ApiAggregateType=[]string{
	"PER_ENDPOINT","ACROSS_ENDPOINTS",
}
var UserAggregateType=[]string{
	"PER_USER","ACROSS_USERS",
}

var ThresholdConfigType=[]string{
	"DYNAMIC","ROLLING_WINDOW",
}

var RRSingle=[]string {
	"RESPONSE_BODY_SIZE","URL","HOST","HTTP_METHOD","USER_AGENT","REQUEST_BODY","REQUEST_BODY_SIZE","QUERY_PARAMS_COUNT","REQUEST_COOKIES_COUNT","REQUEST_HEADERS_COUNT","RESPONSE_COOKIE","RESPONSE_BODY","STATUS_CODE","RESPONSE_COOKIES_COUNT","RESPONSE_HEADERS_COUNT",
}
var MatchOperatorOptions=[]string{"EQUALS","NOT_EQUAL"," CONTAINS","NOT_CONTAIN","MATCHES_REGEX","NOT_MATCH_REGEX"}

var RRMulti=[]string{
	"RESPONSE_HEADER","QUERY_PARAMETER","REQUEST_BODY_PARAMETER","REQUEST_COOKIE","REQUEST_HEADER","RESPONSE_BODY_PARAMETER",
}

var IPAbuseVelocity=[]string{
	"LOW","MEDIUM","HIGH",
}

var Enabled=[]bool{
	true,false,
}
//Every key can have different match operators to support different operator
var MatchOperatorskey=map[string][]string{
  "QUERY_PARAMETER":{"NOT_EQUAL","CONTAINS","NOT_CONTAIN","MATCHES_REGEX","NOT_MATCH_REGEX"},
  "REQUEST_BODY_PARAMETER":MatchOperatorOptions,
  "RESPONSE_HEADER":MatchOperatorOptions,
  "REQUEST_COOKIE":MatchOperatorOptions,
  "REQUEST_HEADER":MatchOperatorOptions,
  "RESPONSE_BODY_PARAMETER":MatchOperatorOptions,
  "RESPONSE_BODY_SIZE":{"EQUALS","NOT_EQUAL","GREATER_THAN","LESS_THAN"},
  "URL":MatchOperatorOptions,
  "HOST":MatchOperatorOptions,
  "USER_AGENT":MatchOperatorOptions,
  "REQUEST_BODY":MatchOperatorOptions,
  "REQUEST_BODY_SIZE":{"EQUALS","NOT_EQUAL","GREATER_THAN","LESS_THAN"},
  "QUERY_PARAMS_COUNT":MatchOperatorOptions,
  "REQUEST_COOKIES_COUNT":MatchOperatorOptions,
  "RESPONSE_COOKIE":MatchOperatorOptions,
  "STATUS_CODE":MatchOperatorOptions,
  "RESPONSE_COOKIES_COUNT":MatchOperatorOptions,
  "RESPONSE_HEADERS_COUNT":MatchOperatorOptions,
}

var DataLocation=[]string{
  "REQUEST","RESPONSE","REQUEST_RESPONSE",
}

var EnumerationThresholdConfigtype=[]string{
  "PATH_PARAMS","REQUEST_BODY","SENSITIVE_PARAMS",
}

