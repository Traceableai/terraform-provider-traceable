build-local:
	go build -o terraform-provider-traceable
	mv terraform-provider-traceable ~/.terraform.d/plugins/registry.terraform.io/traceableai/traceable/0.0.1/darwin_arm64/
	rm -rf .terraformrc terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl .terraform
	terraform init
