resource "traceable_waap_config" "waap" {
  environment = "env_name"
  rule_configs=[
    { 
      rule_name="XSS"
      disabled=false
      subrules=[
        {
          sub_rule_name="IE XSS Filters - Attack (A href with a link)"
          sub_rule_action="BLOCK"
        },
        {
          sub_rule_name="IE XSS Filters - Attack (Applet Tag)"
          sub_rule_action="MONITOR"
        }
      ]
    },
    { 
      rule_name="HTTPProtocolAttack"
      disabled=false
      subrules=[
        {
          sub_rule_name="HTTP Header Injection Attack via headers"
          sub_rule_action="MONITOR"
        },
        {
          sub_rule_name="IIS 6.0 WebDAV buffer overflow: (CVE-2017-7269)"
          sub_rule_action="DISABLE"
        }
      ]
    }
  ]
}