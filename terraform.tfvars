# ip_type_rule_block = [
#   {
#     name = "amazing_rule_7"
#   description = "nice rule"
#   event_severity = "LOW"
#   environment = ["utkarsh_crapi"]
#   ip_types = ["BOT","ANONYMOUS_VPN"]
#   },
#   {
#     name = "amazing_rule_6"
#   description = "very nice rule"
#   event_severity = "MEDIUM"
#   environment = ["utkarsh_21"]
#   ip_types = ["ANONYMOUS_VPN"]
#   },
#   {
#   name = "amazing_rule_5"
#   description = "very very nice rule"
#   event_severity = "HIGH"
#   environment = ["utkarsh_21", "utkarsh_crapi"]
#   ip_types = ["BOT"]
#   }
# ]
# ip_type_rule_block = [
#   {
#     name = "amazing_rule_a"
#     environment = ["utkarsh_crapi"]
#     inject_request_headers = [
#       {
#         header_key   = "X-Trace-ID_09"
#         header_value = "JAI MODI"
#       }
#     ]
#   },
#   {
#     name = "amazing_rule_b"
#     event_severity = "HIGH"
#   },
#   {
#   name = "amazing rule_c"
#   description = "extremely nice rule"
#   },
#   {
#   name = "rule_11"
#   },
#   {
#   name = "rule_12"
#   ip_types = ["ANONYMOUS_VPN"]
#   },
#   {
#   name = "rule_13"
#   description = "what a rule"
#   },
#   {
#   name = "rule_14"
#   event_severity = "LOW"
#   inject_request_headers = [
#       {
#         header_key   = "X-Trace-ID_1"
#         header_value = "123456789"
#       }
#     ]
#   },
#   {
#   name = "rule_15"
#   ip_types = ["HOSTING_PROVIDER"]
#   inject_request_headers = [
#       {
#         header_key   = "X-Trace-ID_2"
#         header_value = "12345612"
#       }
#     ]
#   },
#   {
#   name = "rule_16"
#   ip_types = ["HOSTING_PROVIDER", "TOR_EXIT_NODE"]
#   inject_request_headers = [
#       {
#         header_key   = "X-Trace-ID_3"
#         header_value = "12345613"
#       }
#     ]
#   },
#   {
#   name = "rule_17"
#   ip_types = ["PUBLIC_PROXY"]
#   inject_request_headers = [
#       {
#         header_key   = "X-Trace-ID_4"
#         header_value = "123456793"
#       }
#     ]
#   },
#   {
#   name = "rule_19"
#   ip_types = ["HOSTING_PROVIDER", "TOR_EXIT_NODE", "BOT"]
#   inject_request_headers = [
#       {
#         header_key   = "X-Trace-ID"
#         header_value = "123456"
#       }
#     ]
#   }
# ]
# email_domain_block = [
#   {
#     name                    = "email_rule_1"
#     description             = "Blocks known fraud domains"
#     rule_action             = "BLOCK"
#     event_severity          = "HIGH"
#     expiration              = "PT2000S"
#     environment             = ["utkarsh_21"]
#     data_leaked_email       = true
#     disposable_email_domain = false
#     email_domains           = ["fraud.com", "scam.com"]
#     email_regexes           = [".*fraud.*", ".*scam.*"]
#     email_fraud_score       = "HIGH"
#   },
#   {
#     name = "email_rule_2"
#   }
# ]
# email_domain_rules = [
#   {
#     name                    = "kya"
#     description             = ""
#     event_severity          = "MEDIUM"
#     environment             = ["utkarsh_crapi", "xml-modsec-automation-blocking-enabled"]
#     data_leaked_email       = false
#     disposable_email_domain = false
#     email_domains           = ["example.com"]
#     email_fraud_score_type  = "MIN_SEVERITY"
#     min_email_fraud_score_level = "HIGH"
#   }
# ]

