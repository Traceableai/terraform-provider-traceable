terraform {
  required_providers {
    traceable = {
      source  = "traceableai/traceable"
      version = "0.0.1"
    }
    # aws = {
    #   source  = "hashicorp/aws"
    #   version = "5.35.0"
    # }
  }
}


provider "traceable" {
  platform_url ="https://app-dev.traceable.ai/graphql"
  api_token    ="Bearer "
  
}

resource "traceable_data_set" "sampledataset"{
             description = "hello I am good"
            icon_type = "Password"
            name= "shreyansh123"
}


# resource "traceable_data_set" "sampledataset"{
#   # name = "PII India"
#   name = "shreyanshgupta123"
#   icon_type = "Password"
#   # icon_type = "Financial"
#   description = "create by improved version of provider"

# }

# resource "aws_instance" "example" {
 
# }

