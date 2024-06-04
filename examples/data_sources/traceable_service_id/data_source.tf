data "traceable_service_id" "endpoint" {
  service_name="nginx-automation-test"
  enviroment_name="fintech-1"
}

output "traceable_service_id" {
  value = data.traceable_service_id.endpoint.service_id
}