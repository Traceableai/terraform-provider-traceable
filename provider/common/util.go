package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type GraphqlProviderConfig struct {
	GQLServerUrl string
	ApiToken     string
}

func CallExecuteQuery(query string, meta interface{}) (string, error) {
	url := meta.(*GraphqlProviderConfig).GQLServerUrl
	api_token := meta.(*GraphqlProviderConfig).ApiToken

	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GraphQL request failed with status: %d - %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func ReturnQuotedStringList(inputList []interface{}) []string {
	var envList []string
	for _, env := range inputList {
		envList = append(envList, fmt.Sprintf(`"%s",`, env.(string)))
	}
	lastValue := envList[len(envList)-1]
	lastValue = lastValue[:len(lastValue)-1]
	envList[len(envList)-1] = lastValue
	return envList
}

func CallGetRuleDetailsFromRulesListUsingIdName(response map[string]interface{}, arrayJsonKey string, args ...string) map[string]interface{} {
	var res map[string]interface{}
	rules := response["data"].(map[string]interface{})[arrayJsonKey].(map[string]interface{})
	results := rules["results"].([]interface{})
	id_name := args[0]
	if len(args) == 1 {
		args = append(args, "id")
		args = append(args, "name")
	}
	// log.Println(id_name)
	// log.Println(results)
	for _, rule := range results {
		ruleData := rule.(map[string]interface{})
		// log.Println(ruleData)
		rule_id := ruleData[args[1]].(string)
		var rule_name string
		var ok bool
		if rule_name, ok = ruleData[args[2]].(string); ok {
			// fmt.Println("Rule Name:", rule_name)
		} else {
			rule_name = ""
		}
		if rule_id == id_name || rule_name == id_name {
			// log.Println("Inside if block %s",rule)
			return rule.(map[string]interface{})
		}
	}
	return res
}

func ConvertDurationToSeconds(duration string) (string, error) {
	re := regexp.MustCompile(`P(T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) == 0 {
		return "", fmt.Errorf("invalid duration format: %s", duration)
	}

	totalSeconds := 0

	// Parse hours, minutes, and seconds if present
	if matches[2] != "" { // Hours
		hours, _ := strconv.Atoi(matches[2])
		totalSeconds += hours * 3600
	}
	if matches[3] != "" { // Minutes
		minutes, _ := strconv.Atoi(matches[3])
		totalSeconds += minutes * 60
	}
	if matches[4] != "" { // Seconds
		seconds, _ := strconv.Atoi(matches[4])
		totalSeconds += seconds
	}

	return fmt.Sprintf("PT%dS", totalSeconds), nil
}

func InterfaceToStringSlice(arr []interface{}) []string {
	var envList []string
	for _, env := range arr {
		envList = append(envList, fmt.Sprintf(`"%s"`, env.(string)))
	}
	return envList
}
func InterfaceToEnumStringSlice(arr []interface{}) []string {
	var envList []string
	for _, env := range arr {
		envList = append(envList, fmt.Sprintf(`%s`, env))
	}
	return envList
}

func GetIdFromResponse(responseStr string, graphQlResponseKey string) (string, error) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	responseData, ok := response["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("some error eccorred while fetching response")
	}
	updatedId, _ := responseData[graphQlResponseKey].(map[string]interface{})["id"].(string)
	return updatedId, nil
}
