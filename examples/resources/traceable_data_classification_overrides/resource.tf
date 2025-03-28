resource "traceable_data_classification_overrides" "example_override" {
  name                      = "example_override"
  description               = "This is an example override for data classification"
  data_suppression_override = "REDACT"
  environments              = []

  span_filter {
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
