package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	params := requestCompletion{
		Model:            "text-davinci-003",
		Temperature:      0.9,
		MaxTokens:        600,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
	}
	params.Prompt = "Apa itu wibu?"

	res, err := fetchCompletions(params)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(res.Choices[0].Text)
}

type responseChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     string `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type responseUsage struct {
	PromptToken      int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type responseError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type responseCompletion struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int              `json:"created"`
	Model   string           `json:"model"`
	Choices []responseChoice `json:"choices"`
	Usage   responseUsage    `json:"usage"`
	Error   responseError    `json:"error"`
}

type requestCompletion struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             float64 `json:"top_p"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
}

func fetchCompletions(paramReq requestCompletion) (responseCompletion, error) {
	client := &http.Client{}
	dataRes := responseCompletion{}

	payload, err := json.Marshal(paramReq)
	if err != nil {
		fmt.Println(err)
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(payload))
	if err != nil {
		return dataRes, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	response, err := client.Do(request)
	if err != nil {
		return dataRes, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dataRes)
	if err != nil {
		return dataRes, err
	}

	if len(dataRes.Choices) > 0 {
		dataRes.Choices[0].Text = strings.Replace(dataRes.Choices[0].Text, "\n", "", 2)
	}

	return dataRes, nil
}
