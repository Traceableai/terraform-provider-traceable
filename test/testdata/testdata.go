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

  "HEADER_NAME":MatchOperatorOptions,
  "HEADER_VALUE":MatchOperatorOptions,
  "BODY":MatchOperatorOptions,
  "COOKIE_NAME":MatchOperatorOptions,
  "COOKIE_VALUE":MatchOperatorOptions,
  "HEADERS_COUNT":{"EQUALS","NOT_EQUAL","GREATER_THAN","LESS_THAN"},
  "COOKIES_COUNT": {"EQUALS","NOT_EQUAL","GREATER_THAN","LESS_THAN"},
   "PARAMETER_VALUE":MatchOperatorOptions,
   "PARAMETER_NAME":MatchOperatorOptions,
   "HTTP_METHOD":MatchOperatorOptions,
   "BODY_SIZE":{"EQUALS","NOT_EQUAL","GREATER_THAN","LESS_THAN"},
   "Key":{"EQUALS","NOT_EQUAL","GREATER_THAN","LESS_THAN"},
   "Value":MatchOperatorOptions,
   "Header":MatchOperatorOptions,
   "Parameter":MatchOperatorOptions,
   "Cookie":MatchOperatorOptions,

  

}

var DataLocation=[]string{
  "REQUEST","RESPONSE","REQUEST_RESPONSE",
}

var EnumerationThresholdConfigtype=[]string{
  "PATH_PARAMS","REQUEST_BODY","SENSITIVE_PARAMS",
}

var EnumerationRuleType=[]string{
  "ALERT","BLOCK",
}

var DataSetName=[]string{
  "ui_automation_data_set","test_dataset_a",
}


var DlpCustomLocationAttr=map[string][]string{
  "REQUEST_HEADER":{"MATCHES_REGEX","EQUALS"},
  "REQUEST_COOKIE":{"MATCHES_REGEX","EQUALS"},
  "QUERY_PARAMETER":{"MATCHES_REGEX","EQUALS"},
  "REQUEST_BODY_PARAMETER":{"MATCHES_REGEX","EQUALS"},
  "REQUEST_BODY":{""},
 
}

var UrlRegexes=[][]string{{"hello"},{"hello","hello"}}


var DlpServiceName = "crapi_jacob"


var DlpRRSingle =[]string {
	"URL","HOST","HTTP_METHOD","USER_AGENT","REQUEST_BODY","REQUEST_BODY_SIZE","QUERY_PARAMS_COUNT","REQUEST_COOKIES_COUNT","REQUEST_HEADERS_COUNT","RESPONSE_COOKIE","RESPONSE_BODY","STATUS_CODE","RESPONSE_COOKIES_COUNT",
}


var CustomKeys = map[string][]string{
  "REQUEST":{"URL","HEADER_NAME","HEADER_VALUE","PARAMETER_VALUE","PARAMETER_NAME","HTTP_METHOD","HOST","USER_AGENT","BODY","COOKIE_NAME","COOKIE_VALUE","BODY_SIZE","QUERY_PARAMS_COUNT","HEADERS_COUNT","COOKIES_COUNT"},
  "RESPONSE":{"HEADER_NAME","HEADER_VALUE","STATUS_CODE","BODY","BODY_SIZE","COOKIE_NAME","COOKIE_VALUE","HEADERS_COUNT","COOKIES_COUNT"},
}

var CustomSecRuleOptions=`SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
     "id:92100120,\
     phase:2,\
     block,\
     msg:'Test sec Rule',\
     logdata:'Matched Data: %{TX.0} found within %{MATCHED_VAR_NAME}: %{MATCHED_VAR}',\
     tag:'attack-protocol',\
     tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
     tag:'traceable/severity/HIGH',\
     tag:'traceable/type/safe,block',\
     severity:'CRITICAL',\
     setvar:'tx.anomaly_score_pl1=+%{tx.critical_anomaly_score}'"`


    
var CustomMultiKeys = map[string][]string{
  "REQUEST":{"HEADER","PARAMETER","COOKIE"},
  "RESPONSE":{"HEADER","COOKIE"},
} 








