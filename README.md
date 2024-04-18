# Tf provider set up
Check in into the repo (i.e. go into the folder traceable-terraform-provider)
Execute the following. Commands:
```markdown
go mod init
go mod tidy
go build -o terraform-provider-example
```
Now go to the root directory by doing cd (I.e ~)
```markdown
cd
```
Do ls to check if .terraform.d folder is present or not, if yes then proceed with the steps after
```markdown
ls -al ~
cd .terraform.d 
mkdir -p plugins/terraform.local/local/example/0.0.1/darwin_amd64/
cd
mv /Users/<USER>/Desktop/traceable-terraform-provider/terraform-provider-example .terraform.d/plugins/terraform.local/local/example/0.0.1/darwin_amd64 
```

```markdown
vi .terraformrc
```
Paste the following. In the vi file:
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
Open the repo in vscode and open the terminal, 
Your current directory should be the repo - traceable-terraform-provider, Run the following.
```markdown
terraform init 
terraform apply 
```
After running apply, you will be asked if you want to perform the actions, type yes.

To see logs, export the following. And run apply
```markdown
export TF_LOG="DEBUG"
terraform apply
```


