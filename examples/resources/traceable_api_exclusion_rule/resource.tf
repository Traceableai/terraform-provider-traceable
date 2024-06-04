resource "traceable_api_exclusion_rule" "example_exclusion_rule" {
  name =  "test-rule-exclusion"
  disabled= true
  regexes=  "hello/test/6785"
  service_names=  ["example-svc"]
  environment_names=  ["example-env"]
}