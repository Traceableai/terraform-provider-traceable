# For Static Labels
resource "traceable_label_management_label_rule" "example_label_rule_static" {
  name        = "example_label_rule_static"
  description = "This is an example label rule"
  enabled     = true
  
  condition_list {
    key = "http.method"
    operator = "OPERATOR_EQUALS"
    value = "POST"
  }

  action {
    type = "STATIC_LABELS"
    entity_types = ["API", "SERVICE", "BACKEND"]
    operation = "OPERATION_MERGE"
    static_labels = ["Sensitive"]
  }
}

# For Dynamic Labels
resource "traceable_label_management_label_rule" "example_label_rule_dynamic" {
  name        = "example_label_rule_dynamic"
  description = "This is an example label rule"
  enabled     = true
  
  condition_list {
    key = "http.method"
    operator = "OPERATOR_EQUALS"
    value = "POST"
  }

  action {
    type = "DYNAMIC_LABEL"
    entity_types = ["API", "SERVICE", "BACKEND"]
    operation = "OPERATION_MERGE"
    dynamic_labels {
      attribute = "http.enduser.id"
      regex = "abcd123(.*)"
    }
  }
}
