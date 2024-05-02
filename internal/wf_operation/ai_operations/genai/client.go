package genai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	BaseURL          string
	AppID            string
	AppSecret        string
	IdentityEndpoint string
	IdentityJobID    string
	APIVersion       string
}

type IdentityResponse struct {
	Data struct {
		IdentitySignInInternalApplicationWithPrivateAuth struct {
			AuthorizationHeader string `json:"authorizationHeader"`
		} `json:"identitySignInInternalApplicationWithPrivateAuth"`
	} `json:"data"`
}

func (client *HttpClient) GetAuthorizationHeaderFromIdentityService() (string, error) {
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Intuit_IAM_Authentication intuit_appid=%s, intuit_app_secret=%s", client.AppID, client.AppSecret)
	headers["Content-Type"] = "application/json"

	requestBody := fmt.Sprintf(`{"query":"mutation identitySignInInternalApplicationWithPrivateAuth($input: Identity_SignInApplicationWithPrivateAuthInput!) { identitySignInInternalApplicationWithPrivateAuth(input: $input) { authorizationHeader }}","variables":{"input":{"profileId":%s}}}`, client.IdentityJobID)

	req, err := http.NewRequest("POST", client.IdentityEndpoint+"/v1/graphql", bytes.NewBufferString(requestBody))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error getting authorization header: %v", resp.Status)
	}

	var identityResponse IdentityResponse
	if err := json.NewDecoder(resp.Body).Decode(&identityResponse); err != nil {
		return "", err
	}

	authorizationHeader := fmt.Sprintf("%s,intuit_appid=%s,intuit_app_secret=%s", identityResponse.Data.IdentitySignInInternalApplicationWithPrivateAuth.AuthorizationHeader, client.AppID, client.AppSecret)

	return authorizationHeader, nil
}

func (client *HttpClient) GetPostRequest(requestData string, endpointSuffix string) (interface{}, error) {
	if !json.Valid([]byte(requestData)) {
		return nil, fmt.Errorf("Unable to generate token for GenAI")
	}
	body := []byte(requestData)

	authorizationHeader, err := client.GetAuthorizationHeaderFromIdentityService()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", client.BaseURL+"/"+client.APIVersion+endpointSuffix, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	resDataBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	var resData interface{}
	err = json.Unmarshal(resDataBytes, &resData)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	return resData, nil
}
