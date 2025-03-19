terraform {
  required_providers {
    traceable = {
      source  ="terraform.local/local/traceable"  
      version = "0.0.1"
    }
  }
  
}

variable "platform_url" {
  type        = string
  description = "Traceable Platform URL"
}


variable "traceable_api_key" {
  type        = string
  description = "Traceable API Key"
  sensitive   = true
}


variable "name" {
  description = "Name of the dataset"
  type        = string
}

variable "description" {
  description = "Description of the dataset"
  type        = string
}

variable "icon_type" {
  description = "Icon type for the dataset"
  type        = string
}





provider "traceable" {
  platform_url =var.platform_url
  api_token    =var.traceable_api_key
}


resource "traceable_data_sets" "test_dataset" {
  name        = var.name
  description = var.description
  icon_type   = var.icon_type
}

# Output values for testing
output "dataset_id" {
  value = traceable_data_sets.test_dataset.id
}

output "dataset_name" {
  value = traceable_data_sets.test_dataset.name
}

output "dataset_description" {
  value = traceable_data_sets.test_dataset.description
}

output "dataset_icon_type" {
  value = traceable_data_sets.test_dataset.icon_type
}