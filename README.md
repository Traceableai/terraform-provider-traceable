# Traceable Terraform Provider

## Requirements

- Terraform 1.5.7
- Go 1.22.2

## Building the provider

First make the correct directory, cd to it, and checkout the repo.
```markdown
git clone https://github.com/Traceableai/traceable-terraform-provider.git
cd traceable-terraform-provider
```

Convert the existing Go project to use modules by setting up a new module by creating the `go.mod` file that describes the module's properties including its name and dependencies.
Clean up the module by adding missing and removing unused modules so that the `go.mod` file matches the source code in the module
Compile the Go program in the current directory into a binary
```markdown
go mod init
go mod tidy
go build -o terraform-provider-example
```

## Traceable provider

#### Example usage:
```markdown
provider "example" {
  platform_url=""
  api_token=""
}
```
#### Schema:

##### Required:

- `platform_url`: The platform url to be used
- `api_token`: API token for accessing the platform

##### Optional:
- none

## Resources

### IP range resource

#### Example usage:

```markdown
resource "example_ip_range_rule" "my_ip_range" {
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
- `raw_ip_range_data`: (set of string) list of the ips that the rule would apply on
- `environment`: (set of string) list of the env for which the rule would be applicable

##### Optional:
- `expiration`: (string) expiration of the rule
- `description`: (string) description of the rule

## Plugin set up (to run locally)

cd into the `.terraform.d` directory under root - this directory is used by Terraform to store global configurations and data, including custom provider plugins.
Create a new directory path for a custom Terraform provider, where Terraform will look for local provider plugins.

```markdown
mkdir -p plugins/terraform.local/local/example/0.0.1/darwin_amd64/
```

- `terraform.local/local/example` simulates a local namespace for the provider.
- `0.0.1` is the version of the provider.
- `darwin_amd64` indicates the build is for macOS architecture, which is common for Mac users (darwin stands for macOS, and amd64 indicates the 64-bit architecture).

Move the compiled Terraform provider binary to the directory within `.terraform.d` that was created above

```markdown
mv /Users/<USER>/Desktop/traceable-terraform-provider/terraform-provider-example .terraform.d/plugins/terraform.local/local/example/0.0.1/darwin_amd64/
```

Create a file .terraformrc that is used to configure various behaviors of Terraform, including provider installation paths and methods. This one Specifies custom provider installation settings for Terraform.
The of these settings is to have Terraform look locally in /Users/<USER>/.terraform.d/plugins for any provider plugins before checking online, and to specifically avoid online checks for any providers that are intended to be managed locally (as indicated by the namespace pattern terraform.local/*/*).

```markdown
vi .terraformrc
```
paste the following in the file:

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

Initialize a Terraform working directory and Apply the changes required to reach the desired state of the configuration

```markdown
terraform init
terraform apply
```





