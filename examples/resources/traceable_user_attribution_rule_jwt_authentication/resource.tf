resource "traceable_user_attribution_rule_jwt_authentication" "test3" {
  name = "jwtauth"
  scope_type = "CUSTOM"
  url_regex="sfdsf"
  jwt_location = "COOKIE"
  jwt_key = "abcd"
  user_id_claim = "aditya"
}