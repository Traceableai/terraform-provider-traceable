data "traceable_service_id" "endpoint" {
  service_name="test-service"
  enviroment_name="test-env"
}

output "traceable_service_id" {
  value = data.traceable_service_id.endpoint.service_id
}