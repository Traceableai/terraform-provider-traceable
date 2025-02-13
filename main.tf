terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
    # aws = {
    #   source  = "hashicorp/aws"
    #   version = "5.35.0"
    # }
  }
  # backend "s3" {
  #   bucket = "traceable-provider-store"
  #   key    = "traceable-provider-store"
  #   region = "us-west-2"
  # }
# }
}

# 

