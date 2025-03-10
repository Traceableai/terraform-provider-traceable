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
  api_token    ="Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkxPbDdCcnhCVzUweUYxVERNYWRpZyJ9.eyJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlc192MiI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvY3VzdG9tZXJfaWQiOiIzZTc2MTg3OS1jNzdiLTRkOGYtYTA3NS02MmZmMjhlOGZhOGEiLCJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlcyI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvanRpIjoiYzVjMWY2OTItOWFmYi00YmMxLWIwMTItYTY5ZDRmN2IzMTNmIiwiaHR0cHM6Ly90cmFjZWFibGUuYWkvcmljaF9yb2xlcyI6W3siZW52cyI6W10sImlkIjoidHJhY2VhYmxlIn1dLCJuaWNrbmFtZSI6InRlc3QrZTJlIiwibmFtZSI6InRlc3QrZTJlQHRyYWNlYWJsZS5haSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8wZWFiNjdkMTBlZTlhNTYwZGY2ZDJiODE4NTQ3NGM5ZT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRnRlLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDI1LTAzLTEwVDEwOjQyOjQyLjMxMVoiLCJlbWFpbCI6InRlc3QrZTJlQHRyYWNlYWJsZS5haSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczovL3RyYWNlYWJsZS1kZXYudXMuYXV0aDAuY29tLyIsImF1ZCI6IkloeHg4QVBlc3dEaUlrb3JxSzIzYkt0OHVUa0pkaDA3Iiwic3ViIjoiYXV0aDB8NjMxMmJmMDRiZGVhODBmYzY1ZWFkMmE5IiwiaWF0IjoxNzQxNjAzMzYzLCJleHAiOjE3NDE2MzkzNjMsInNpZCI6IktDYWlCUTJzZG8tMHdYNTMtdmVWaWZaY2VRTm9rcVhpIiwibm9uY2UiOiJ5NUNFbzNZbEYxOE5OeWVJZHF4SEl4ZDB4clpXbjZ4S1U1NUpocFM3V2RZIn0.y6JZBIArMe-JsIRN4ExyaSpnU9ZRzkWgHDhiJUzQnjo_jlUKB0HMBxiY_yzWesU1z8wnYq4F-lZ9DYONjcmfy6Byi7wD5thuayovvMc8PXGuQ9hFJAcip4BG8UzRTr2if17OgyaVBpFQA1NS12eVCCpQ6yjVT1pAu9kty8qjeoSHuss2EevJZ-95ei-6QJGgK4EiWDzCj75_Imy38f5sCHyv4RuIfn6SyezCFtrnS_zJIFRLxZ4q8UJbHFgKEkTzHyicIVKeWZpsjENQySY5F7do4OgzxqwaTjTYHKxW5hIlb3W_BbKID_j_9IC0av2vNawJ3omEDiRlOZvnIIJSzw"
  
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

