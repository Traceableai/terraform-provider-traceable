resource "traceable_user_attribution_rule_custom_token" "test6" {
  name = "traceable_user_attribution_rule_custom_token"
  scope_type="SYSTEM_WIDE"
  auth_type="test"
  location="REQUEST_COOKIE"
  token_name="test"
}