resource "traceable_api_naming_rule" "example_naming_rule" {
  name             = "test-rule-naming"
  disabled         = false
  regexes          = ["hello", "test", "123"]
  values           = ["hello", "test", "number"]
  service_names    = ["example-svc"]
  environment_names = ["example-env"]
}