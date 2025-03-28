# atmost one of rule_config or subrule_config
# rule_config case
resource "traceable_detection_policies" "sample_rule" {
  environment = "env"
  waap_config {
    rule_id = "XSS"
    rule_config  {
      disabled = false
    }
  }
}

# subrule_config case
resource "traceable_detection_policies" "sample_rule" {
  environment = "utkarsh_crapi"
  waap_config {
    rule_id = "XSS"
    subrule_config {
      sub_rule_id = "crs_941280"
      sub_rule_action = "MONITOR"
    }
  }
}