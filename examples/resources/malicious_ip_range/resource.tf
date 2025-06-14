

resource "traceable_malicious_ip_range" "block_all_except_sample"{
    name = "traceable-source"
    description = "traceable-des"
    enabled = true
    event_severity = "LOW"
    duration = "PT1M"
    action = "BLOCK_ALL_EXCEPT"
    ip_range = ["192.168.1.1","192.168.1.2"]
    environments = ["env1","env2"]
}

resource "traceable_malicious_ip_range" "allow_sample"{
    name = "traceable-source"
    description = "traceable-des"
    enabled = true
    duration = "PT1M"
    action = "ALLOW"
    ip_range = ["192.168.1.1","192.168.1.2"]
    environments = ["env1","env2"]
}

resource "traceable_malicious_ip_range" "alert_sample"{
    name = "traceable-source"
    description = "traceable-des"
    enabled = true
    event_severity = "LOW"
    action = "Alert"
    ip_range = ["192.168.1.1","192.168.1.2"]
    environments = ["env1","env2"]
}

resource "traceable_malicious_ip_range" "block_sample"{
    name = "traceable-source"
    description = "traceable-des"
    enabled = true
    event_severity = "LOW"
    duration = "PT1M"
    action = "BLOCK_ALL_EXCEPT"
    ip_range = ["192.168.1.1","192.168.1.2"]
    environments = ["env1","env2"]
}