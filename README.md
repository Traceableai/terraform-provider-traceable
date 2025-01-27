# Traceable Terraform Provider

## Requirements

- Terraform 1.5.7
- Go 1.22.2

## Building the provider

Clone the repo
```markdown
git clone https://github.com/traceableai/terraform-provider-traceable.git
cd traceable-terraform-provider
```

Install packages and build
```markdown
go mod tidy
go build -o terraform-provider-traceable
```

## Traceable provider

#### Example usage:
```markdown
provider "traceable" {
  platform_url=""
  api_token=""
}
```
#### Schema:

##### Required:

- `platform_url`: The platform url to be used
- `api_token`: API token for accessing the platform

## Resources

### IP range resource

#### Example usage:

```markdown
resource "traceable_ip_range_rule" "my_ip_range" {
    name     = "first_rule_2"
    rule_action     = "RULE_ACTION_ALERT"
    event_severity     = "LOW"
    raw_ip_range_data = [
        "1.1.1.1",
        "3.3.3.3"
    ]
    environment=[]
    description="rule created from custom provider"
}
```

#### Schema:

##### Required:

- `name`: (string) name of the rule
- `rule_action`: (string) type of action of the rule to be created [RULE_ACTION_BLOCK, RULE_ACTION_ALERT, RULE_ACTION_ALLOW]
- `event_severity`: (string) severity of the rule [LOW, MEDIUM, HIGH]
- `raw_ip_range_data`: (set of string) list of the ip addresses(IPv4 or IPv6), cidr ranges that the rule would apply on
- `environment`: (set of string) list of the env for which the rule would be applicable

##### Optional:
- `expiration`: (string) expiration time of the rule (this attribute don't apply on `RULE_ACTION_ALERT`, don't pass this attribute if we need to block or allow indefinetly)
- `description`: (string) description of the rule

## Plugin set up (to run locally)

Inside `~/.terraform.d` directory

Create a new directory path for a custom Terraform provider, where Terraform will look for local provider plugins.

```markdown
mkdir -p plugins/terraform.local/local/traceable/0.0.1/darwin_amd64/
```

Move provider binary `terraform-provider-traceable` to this directory ```plugins/terraform.local/local/traceable/0.0.1/darwin_amd64/```
`darwin_amd64` this should match with your os arch

Add provider binary path to your `.terraformrc` (if not exist create and paste below lines) to run it locally

```markdown
provider_installation {
  filesystem_mirror {
    path    = "/Users/<USER>/.terraform.d/plugins"
  }
  direct {
    exclude = ["terraform.local/*/*"]
  }
}
```

Initialize a Terraform working directory and Apply the changes.

```markdown
terraform init
terraform apply
```





