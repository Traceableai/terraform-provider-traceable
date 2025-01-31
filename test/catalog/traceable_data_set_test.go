package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestDataSetsRuleCreation(t *testing.T) {
	t.Parallel() // it makes the test case running paralledl

	// Generate a unique name for parallel test runs
	uniqueName := fmt.Sprintf("test-dataset-%d", time.Now().Unix())

	traceAPIKEy := os.Getenv("TRACEABLE_API_KEY")
	pLATFORMURL := os.Getenv("PLATFORM_URL")

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/resources/traceable_data_set",

		Vars: map[string]interface{}{
			"name":              uniqueName,
			"description":       "Test dataset for integration testing",
			"icon_type":         "Setting",
			"traceable_api_key": traceAPIKEy,
			"platform_url":      pLATFORMURL,
		},
	}

	// it is doing apply and create
	terraform.InitAndApply(t, terraformOptions)

	//after create getting state
	datasetName := terraform.Output(t, terraformOptions, "dataset_name")
	datasetDescription := terraform.Output(t, terraformOptions, "dataset_description")
	datasetIconType := terraform.Output(t, terraformOptions, "dataset_icon_type")

	// check the actual state with desired state
	assert.Equal(t, uniqueName, datasetName)
	assert.Equal(t, "Test dataset for integration testing", datasetDescription)
	assert.Equal(t, "Setting", datasetIconType)

	newDescription := "Updated test dataset description"
	terraformOptions.Vars["description"] = newDescription

	// it is doing update
	terraform.Apply(t, terraformOptions)

	updatedDescription := terraform.Output(t, terraformOptions, "dataset_description")

	assert.Equal(t, newDescription, updatedDescription)

	terraform.Destroy(t, terraformOptions)

	_ = terraform.Output(t, terraformOptions, "dataset_id")

}

// this type of test cases we are using for negative test case
func TestDataSetsRuleValidation(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"name":              "hi how are you",
			"description":       "Test dataset for integration testing",
			"icon_type":         "Setting",
			"traceable_api_key": os.Getenv("TRACEABLE_API_KEY"),
			"platform_url":      os.Getenv("PLATFORM_URL"),
		},
	}

	_, err := terraform.InitAndApplyE(t, terraformOptions)
	assert.Error(t, err)
}
