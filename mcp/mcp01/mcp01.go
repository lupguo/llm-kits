package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ConversionRequest struct {
	SourceData   string `json:"source_data"`
	TargetFormat string `json:"target_format"`
}

type ConversionResponse struct {
	ConvertedData string `json:"converted_data,omitempty"`
	Error         string `json:"error,omitempty"`
}

func (c *OpenAIClient) convertWithLLM(prompt string) (string, error) {
	messages := []map[string]string{
		{"role": "system", "content": "你是一个专业的数据格式转换引擎，严格按用户要求输出"},
		{"role": "user", "content": prompt},
	}
	return c.DoRequest(modelName, messages)
}

func conversionHandler(w http.ResponseWriter, r *http.Request) {
	var req ConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prompt := fmt.Sprintf("将以下数据严格转换为 %s 格式，只输出结果:\n%s",
		req.TargetFormat,
		req.SourceData)

	client, err := NewOpenAIClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("PROXY_URL"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	convertedData, err := client.convertWithLLM(prompt)
	if err != nil {
		json.NewEncoder(w).Encode(ConversionResponse{Error: err.Error()})
		return
	}

	// 二次格式验证
	switch req.TargetFormat {
	case "json":
		if !json.Valid([]byte(convertedData)) {
			json.NewEncoder(w).Encode(ConversionResponse{Error: "LLM 返回无效 JSON"})
			return
		}
	case "yaml":
		// 可添加 YAML 解析校验
	}

	json.NewEncoder(w).Encode(ConversionResponse{ConvertedData: convertedData})
}
