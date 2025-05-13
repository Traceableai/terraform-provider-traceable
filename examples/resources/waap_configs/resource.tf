resource "traceable_waap_config" "waap" {
  environment = "env_name"
  rule_configs=[
    { 
      rule_name="XSS"
      enabled=true
      subrules=[
        {
          name="IE XSS Filters - Attack (A href with a link)"
          action="BLOCK"
        },
        {
          name="IE XSS Filters - Attack (Applet Tag)"
          action="MONITOR"
        }
      ]
    },
    { 
      rule_name="HTTPProtocolAttack"
      enabled=true
      subrules=[
        {
          name="HTTP Header Injection Attack via headers"
          action="MONITOR"
        },
        {
          name="IIS 6.0 WebDAV buffer overflow: (CVE-2017-7269)"
          action="DISABLE"
        }
      ]
    }
  ]
}