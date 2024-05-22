package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func executeQuery(query string, meta interface{}) (string, error) {
	url := meta.(*graphqlProviderConfig).GQLServerUrl
	api_token := meta.(*graphqlProviderConfig).ApiToken

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
