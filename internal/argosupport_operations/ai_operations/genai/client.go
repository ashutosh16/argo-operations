package genai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

func (client *HttpClient) PostRequest(ctx context.Context, requestData string, endpointSuffix string) (interface{}, error) {
	logger := log.FromContext(ctx)
	logger.Info("Rollout seems to be healthy and should not be included in the genai analysis")

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

func (client *HttpClient) GetRequest(fullUrl string, params map[string]string) (*Application, error) {
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	// Adding headers as per the curl command
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "argocd.token="+client.AppSecret+"; Secure; HttpOnly")
	httpClient := &http.Client{

		Timeout: time.Minute * 5,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-OK status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var app Application
	if err := json.Unmarshal(body, &app); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &app, nil
}
