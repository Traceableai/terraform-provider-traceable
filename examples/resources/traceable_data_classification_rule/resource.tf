resource "traceable_data_classification_rule" "example_data_classification" {
  name             = "example_data_classification"
  description      = "This is an example data classification rule"
  enabled          = true
  data_sets        = ["dataset-1"]
  data_suppression = "REDACT"
  sensitivity      = "HIGH"

  scoped_patterns {
    scoped_pattern_name = "example_pattern"
    environments        = []
    match_type          = "MATCH"
    locations           = ["REQUEST_BODY"]

    key_patterns {
      operator = "MATCHES_REGEX"
      value    = ".*example.*"
    }

    value_patterns {
      operator = "EQUALS"
      value    = "sensitive"
    }
  }
}
