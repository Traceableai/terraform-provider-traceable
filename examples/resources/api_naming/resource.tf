resource "traceable_api_naming" "test" {
  name = "sample-api-naming-rule"
  disabled = false
  service_names=["nginx"]
  environment_names=[] #empty for all env
  values=["someval"]
  regexes=["namedep"]
}