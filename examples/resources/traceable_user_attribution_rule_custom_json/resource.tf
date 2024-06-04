resource "traceable_user_attribution_rule_custom_json" "test5" {
  name = "traceable_user_attribution_rule_custom_json"
  scope_type="CUSTOM"
  url_regex="sfdsf"
  auth_type_json=jsonencode(file("authType.json"))
  user_id_json=jsonencode(file("authType.json"))
}