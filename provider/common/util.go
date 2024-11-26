package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	req.Header.Set("Authorization", api_token)

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

func ReturnQuotedStringList(inputList []interface{}) []string{
	var envList []string
    for _, env := range inputList {
        envList = append(envList, fmt.Sprintf(`"%s"`, env.(string)))
    }
	return envList
}