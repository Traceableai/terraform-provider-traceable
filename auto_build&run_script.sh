set -e 

PLUGIN_NAME="terraform-provider-traceable"
PLUGIN_VERSION="0.0.1"
PLUGIN_PATH="$HOME/.terraform.d/plugins/terraform.local/local/traceable/$PLUGIN_VERSION/darwin_amd64"

echo "Building Terraform provider..."
go build -o "$PLUGIN_NAME"

echo "Cleaning up Terraform state and cache..."
rm -rf .terraformrc terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl .terraform

echo "Moving provider binary..."
mkdir -p "$PLUGIN_PATH"
mv "$PLUGIN_NAME" "$PLUGIN_PATH"

echo "Reinitializing Terraform..."
terraform init

echo "Applying Terraform configuration with debug logs..."
TF_LOG=DEBUG terraform apply -auto-approve

echo "Done!"
