resource "traceable_user_attribution_rule_req_header" "test2" {
  name = "traceable_user_attribution_rule_req_header"
  scope_type = "CUSTOM"
  url_regex = "abcd"
  auth_type = "test"
  user_id_location = "test"
  user_role_location="yes"
  role_location_regex_capture_group="test"
}