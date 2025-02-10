# variable "ip_type_rule_block" {
#   type = list(object({
#     name                    = string
#     description             = optional(string)
#     rule_action             = optional(string)
#     event_severity          = string
#     environment             = list(string)
#     ip_types                = list(string)
#   }))
# }
# variable "ip_type_rule_block" {
#   description = "List of IP Type Rule Blocks"
#   type = list(object({
#     name           = string
#     description    = optional(string, "Default rule description")
#     rule_action    = optional(string, "ALERT")
#     event_severity = optional(string, "MEDIUM")
#     environment    = optional(list(string), ["utkarsh_21"])
#     ip_types       = optional(list(string), ["BOT"])

#     inject_request_headers = optional(list(object({
#       header_key   = string
#       header_value = string
#     })), [
#       { header_key = "X-Default-Header", header_value = "DefaultValue" }
#     ])
#   }))
# }
# variable "email_domain_block" {
#   description = "List of Email Domain Block Rules"
#   type = list(object({
#     name                    = string
#     description             = optional(string, "Default email domain policy description")
#     rule_action             = optional(string, "BLOCK")
#     event_severity          = optional(string, "MEDIUM")
#     expiration              = optional(string, "PT2000S")
#     environment             = optional(list(string), ["utkarsh_21"])
#     data_leaked_email       = optional(bool, false)
#     disposable_email_domain = optional(bool, false)
#     email_domains           = optional(list(string), ["example.com"])  # ✅ Ensure this is set
#     email_regexes           = optional(list(string), [".*example.*"])  # ✅ Ensure this is set
#     email_fraud_score       = optional(string, "CRITICAL")
#   }))
# }
# variable "email_domain_rules" {
#   type = list(object({
#     name                    = string
#     description             = string
#     event_severity          = string
#     environment             = list(string)
#     data_leaked_email       = bool
#     disposable_email_domain = bool
#     email_domains           = list(string)
#     email_fraud_score_type  = string
#     min_email_fraud_score_level = string
#   }))
# }
