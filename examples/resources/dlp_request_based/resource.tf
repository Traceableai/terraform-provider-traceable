resource "traceable_data_loss_prevention_request_based" "example" {
  name         = "example_dlp_request_rule_2"
  description  = "Example DLP request-based rule"
  enabled      = true
  environments = ["dev", "prod"]

  action = {
    action_type    = "BLOCK"   # ALERT | BLOCK | ALLOW
    event_severity = "LOW"     # LOW | MEDIUM | HIGH | CRITICAL (only BLOCK / ALERT)
    # duration     = "PT60S"   # Allowed only with ALLOW / BLOCK
  }

  sources = {
    service_scope = {
      service_ids = ["172d77d2-ae7a-3355-9a09-2db0803f3523"]
    }

    url_regex_scope = {
      url_regexes = ["/private/.*"]
    }

    regions = {
      region_ids = ["US", "CA"]
    }

    ip_location_type = {
      ip_location_types = ["BOT", "ANONYMOUS_VPN"]
    }

    ip_address = {
      ip_address_list = ["192.0.2.10","1.1.1.1"]
    }

    request_payload = [
      {
        metadata_type   = "REQUEST_HEADER" 
        key_operator    = "EQUALS"
        key_value       = "Content-Type"
        value_operator  = "CONTAINS"
        value           = "json"
      },
      {
        metadata_type  = "HTTP_METHOD"
        value_operator = "EQUALS"
        value          = "POST"
      }
    ]

    dateset_datatype_filter = {
      dateset_datatype_id = {
        data_sets_ids  = ["4c41bbe3-92a3-42c2-aa78-3d69cef1ef49"]
      }


      data_type_matching = {
        metadata_type = "REQUEST_BODY"
      }
    }
  }
}