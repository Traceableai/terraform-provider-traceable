// test/dataset_rule_test.go

package test

// import (
// 	"fmt"
// 	"github.com/gruntwork-io/terratest/modules/terraform"
// 	"github.com/stretchr/testify/assert"
// 	"os"
// 	"testing"
// 	"time"
// )

// func TestCustomSignatureAllowCreation(t *testing.T) {
// 	t.Parallel() // it makes the test case running paralledl

// 	// Generate a unique name for parallel test runs
// 	uniqueName := fmt.Sprintf("test-custom-signature-%d", time.Now().Unix())

// 	terraformOptions := &terraform.Options{
// 		TerraformDir: "../../examples/resources/traceable_custom_signature_allow",

// 		Vars: map[string]interface{}{
// 			"name":                  uniqueName,
// 			"description":           "Test dataset for integration testing",
// 			"environments":          []string{"ui-data-validation"},
// 			"allow_expiry_duration": "PT30M",

// 			"req_res_conditions": []map[string]interface{}{
// 				{
// 					"match_key":      "HEADER_NAME",
// 					"match_category": "REQUEST",
// 					"match_operator": "EQUALS",
// 					"match_value":    "req_header",
// 				},
// 				{
// 					"match_key":      "HEADER_NAME",
// 					"match_category": "REQUEST",
// 					"match_operator": "EQUALS",
// 					"match_value":    "req_header_test",
// 				},
// 			},

// 			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
// 			"platform_url":      os.Getenv("PLATFORM_URL"),
// 		},
// 	}

// 	// it is doing apply and create
// 	terraform.InitAndApply(t, terraformOptions)

// 	//after create getting state
// 	customName := terraform.Output(t, terraformOptions, "custom_name")
// 	customDescription := terraform.Output(t, terraformOptions, "custom_description")

// 	// check the actual state with desired state
// 	assert.Equal(t, uniqueName, customName)
// 	assert.Equal(t, "Test dataset for integration testing", customDescription)

// 	newDescription := "Updated test dataset description"
// 	terraformOptions.Vars["description"] = newDescription

// 	// it is doing update
// 	terraform.Apply(t, terraformOptions)
// 	//check description is successfully updated
// 	updatedDescription := terraform.Output(t, terraformOptions, "custom_description")
// 	assert.Equal(t, newDescription, updatedDescription)
// 	customId := terraform.Output(t, terraformOptions, "custom_id")
// 	fmt.Println(customId)

// 	// it is doing destroy
// 	terraform.Destroy(t, terraformOptions)
// 	//after destroy we can not fetch id we have to get error
// 	_, err := terraform.OutputE(t, terraformOptions, "dataset_id")
// 	assert.Error(t, err, "Expected error while fetching dataset_id after destroy. The resource should no longer exist.")

// }
