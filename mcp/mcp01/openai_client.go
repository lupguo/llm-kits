package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	openAIURL = "https://api.openai.com/v1/chat/completions"
	modelName = "gpt-4"
)

type OpenAIClient struct {
	apiKey     string
	proxyURL   string
	httpClient *http.Client
}

func NewOpenAIClient(apiKey, proxyURL string) (*OpenAIClient, error) {
	client := &http.Client{}
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}
	return &OpenAIClient{apiKey: apiKey, proxyURL: proxyURL, httpClient: client}, nil
}

func (c *OpenAIClient) DoRequest(modelName string, messages []map[string]string) (string, error) {
	requestBody := map[string]interface{}{
		"model":       modelName,
		"messages":    messages,
		"temperature": 0.1,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", openAIURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	body, _ := io.ReadAll(rsp.Body)
	json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from LLM")
}
