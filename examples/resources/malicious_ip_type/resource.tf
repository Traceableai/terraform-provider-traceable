resource "traceable_malicious_ip_type" "sample"{
    name = "mal_ip_type"
    description = "traceable-des"
    enabled = true
    event_severity = "HIGH"
    duration = "PT1M"
    action = "ALERT"
    environments = ["env1","env2"]
    ip_type = ["ANONYMOUS_VPN","BOT",]
}