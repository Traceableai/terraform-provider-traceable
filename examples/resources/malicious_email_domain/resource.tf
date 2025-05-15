resource "traceable_malicious_email_domain" "sample" {
    name = "malicious_email_domain_example"
    description = "example rule"
    enabled = true
    event_severity = "HIGH"
    duration = "PT1M"
    action = "ALERT"
    environments = ["env1","env2"]
    email_domains_list = ["traceable.ai", "harness.io"]
    apply_rule_to_data_leaked_email = true
    min_email_fraud_score_level = "CRITICAL"
    email_regexes_list = ["traceable.ai", "harness.io"]
    apply_rule_to_disposable_email_domains = true
}