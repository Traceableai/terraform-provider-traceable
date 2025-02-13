package test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
	// "log
	"os"
	"testing"
	"fmt"
	"github.com/traceableai/terraform-provider-traceable/test/testlog"

	
)

var names = []string{"test-dataset-a","test-dataset-b"}
var descriptions = []string{"Test dataset 1","Test dataset 2"}
var iconTypes = []string{"Setting","dashboard"}

func generateTestCases() []struct {
	name        string
	description string
	iconType    string
} {

	var testCases []struct {
		name        string
		description string
		iconType    string
	}

	for _, n := range names {
		for _, d := range descriptions {
			for _, i := range iconTypes {
				testCases = append(testCases, struct {
					name        string
					description string
					iconType    string
				}{
					name:      n,
					description: d,
					iconType:    i,
				})
			}
		}
	}

	return testCases
}

func TestDataSetsRuleCreation(t *testing.T) {
  logger:=testlog.Logger

	





	t.Parallel() // enables parallel execution of the test case
	t.Log("Test TestDataSetsRuleCreation is running")
	t.Log(t.Name())
	traceAPIKey := os.Getenv("TRACEABLE_API_KEY")
	platformURL := os.Getenv("PLATFORM_URL")

	tests := generateTestCases()


	

	// Define table-driven test cases

	// Iterate through the test cases
	for _, tt := range tests {

		
		t.Run(fmt.Sprintf(" Name:%s,Description:%s,IconType:%s",tt.name, tt.description, tt.iconType), func(t *testing.T) {
			// Apply Terraform to create resources
			logger.SetTestName(fmt.Sprintf("%s Name:%s,Description:%s,IconType:%s", t.Name(),tt.name, tt.description, tt.iconType))


			terraformOptions := &terraform.Options{
				TerraformDir: "../../examples/resources/traceable_data_set",
				Vars: map[string]interface{}{
					"name":              tt.name,
					"description":       tt.description,
					"icon_type":         tt.iconType,
					"traceable_api_key": traceAPIKey,
					"platform_url":      platformURL,
					
				},
			}

			logger.Log("Starting terraform init and apply")
			terraform.InitAndApply(t, terraformOptions)

			// Verify created resources
			logger.Log("Verifying created resources")
			datasetName := terraform.Output(t, terraformOptions, "dataset_name")
			datasetDescription := terraform.Output(t, terraformOptions, "dataset_description")
			datasetIconType := terraform.Output(t, terraformOptions, "dataset_icon_type")

			// Assert the created resource properties
			logger.Log(fmt.Sprintf("Verifying outputs - Name: %s, Description: %s, IconType: %s", datasetName, datasetDescription, datasetIconType))
			require.Equal(t, tt.description, datasetDescription,"description should be same after creating")
			require.Equal(t, tt.iconType, datasetIconType,"icon type should be same after creating")

			// Apply updates to resources
			// logger.Printf("Starting update phase")
			terraformOptions.Vars["description"] = "new description"
			terraform.Apply(t, terraformOptions)
    
			// Verify the updated description
			updatedDescription := terraform.Output(t, terraformOptions, "dataset_description")
			require.Equal(t, "new description", updatedDescription,"description should be same after updating")

			// Destroy the resources after the test
			terraform.Destroy(t, terraformOptions)
			// logger.Log("Verifying updated description: %s", updatedDescription)
			
		})
	}
}

// this type of test cases we are using for negative test case
func TestDataSetsRuleValidation1(t *testing.T) {
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

	// _, err := terraform.InitAndApplyE(t, terraformOptions)
	// require.Error(t, err)
	terraform.InitAndApply(t, terraformOptions)
}


