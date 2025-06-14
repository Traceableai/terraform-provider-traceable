fetch-schema:
	@if [ -z "$$JWT_TOKEN" ]; then \
		echo "Error: JWT_TOKEN environment variable is not set"; \
		exit 1; \
	fi
	@echo "Installing GraphQL CLI tools..."
	npm install -g graphql-cli graphql-introspection-json-to-sdl
	@echo "GraphQL tools installed successfully"
	@echo "Fetching GraphQL schema..."
	curl -X POST https://app-dev.traceable.ai/graphql \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $$JWT_TOKEN" \
		--data-raw '{"query": "query IntrospectionQuery { __schema { queryType { name } mutationType { name } subscriptionType { name } types { ...FullType } directives { name description locations args { ...InputValue } } } } fragment FullType on __Type { kind name description fields(includeDeprecated: true) { name description args { ...InputValue } type { ...TypeRef } isDeprecated deprecationReason } inputFields { ...InputValue } interfaces { ...TypeRef } enumValues(includeDeprecated: true) { name description isDeprecated deprecationReason } possibleTypes { ...TypeRef } } fragment InputValue on __InputValue { name description type { ...TypeRef } defaultValue } fragment TypeRef on __Type { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name } } } } } } }"}' > schema.json
	@echo "Converting schema to SDL format..."
	npx graphql-introspection-json-to-sdl schema.json > helper/schema.graphql
	@echo "Schema generated: helper/schema.graphql"

build-local:
	go build -o terraform-provider-traceable
	mv terraform-provider-traceable ~/.terraform.d/plugins/registry.terraform.io/traceableai/traceable/0.0.1/darwin_arm64/
	rm -rf .terraformrc terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl .terraform
	terraform init
