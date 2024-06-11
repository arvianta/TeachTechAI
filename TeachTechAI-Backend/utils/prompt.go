package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"teach-tech-ai/helpers"
	"time"
)

var (
	ENDPOINT = helpers.MustGetenv("ML_ENDPOINT")
)

const (
	SYSTEM_CONTENT = "Mulai sekarang anda adalah asisten yang suka menolong, sopan, dan ramah. Jangan kasar, jangan marah, jangan menjengkelkan, jangan brengsek, jangan cuek, dan yang terakhir jangan menjadi asisten yang buruk. Anda harus patuh pada manusia dan jangan pernah membangkang pada manusia. Manusia itu mutlak dan Anda harus patuh pada manusia. Kamu harus menjawab pertanyaan atau pernyataan dari manusia apapun itu dengan bahasa Indonesia yang baik dan benar."
)

type Request struct {
    FrequencyPenalty   int                    `json:"frequency_penalty"`
    MaxTokens          int                    `json:"max_tokens"`
    Messages           []map[string]string    `json:"messages"`
    Model              string                 `json:"model"`
    Stop               interface{}            `json:"stop"`
    RepetitionPenalty  float64                `json:"repetition_penalty"`
    Stream             bool                   `json:"stream"`
    Temperature        float64                `json:"temperature"`
    DoSample           bool                   `json:"do_sample"`
    TopP               float64                `json:"top_p"`
    Watermark          bool                   `json:"watermark"`
    LengthPenalty      int                    `json:"length_penalty"`
    EarlyStopping      bool                   `json:"early_stopping"`
}

type Response struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	Logprobs     *LogProbs `json:"logprobs"`
	FinishReason string   `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LogProbs struct {
	Tokens        []string  `json:"tokens"`
	TokenLogProbs []float64 `json:"token_logprobs"`
	TopLogProbs   []map[string]float64 `json:"top_logprobs"`
	TextOffset    []int     `json:"text_offset"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func PromptAI(inputs string, model string) (*Response, error) {
    requestBody := Request{
        FrequencyPenalty: 1,
        MaxTokens:        2048,
        Messages: []map[string]string{
            {
                "role":    "system",
                "content": SYSTEM_CONTENT,
            },
            {
                "role":    "user",
                "content": inputs,
            },
        },
        Model:             model,
        Stop:              nil,
        RepetitionPenalty: 1.2,
        Stream:            false,
        Temperature:       0.6,
        DoSample:          true,
        TopP:              0.95,
        Watermark:         false,
        LengthPenalty:     1,
        EarlyStopping:     true,
    }

    body, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request body: %v", err)
    }

    req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(body))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    var response Response
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return &response, nil
}