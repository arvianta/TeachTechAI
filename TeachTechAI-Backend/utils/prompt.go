package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Request struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters"`
}

type Response struct {
    GeneratedText string  `json:"generated_text"`
    Details       Details `json:"details"`
}

type Token struct {
	ID      int     `json:"id"`
	Text    string  `json:"text"`
	LogProb float64 `json:"logprob"`
	Special bool    `json:"special"`
}

type Details struct {
	FinishReason    string `json:"finish_reason"`
	GeneratedTokens int    `json:"generated_tokens"`
	Seed            uint64 `json:"seed"`
	Prefill         []Token `json:"prefill"`
    Tokens          []Token `json:"tokens"`
}

var (
	ENDPOINT = MustGetenv("ML_ENDPOINT")
)

func PromptAI(inputs string) (*Response, error) {
	requestBody := Request{
		Inputs: inputs,
		Parameters: map[string]interface{}{
			"best_of": 					1,
			"decoder_input_details": 	false,
			"details": 					true,
			"do_sample": 				false,
			"frequency_penalty":		0.5,
			"max_new_tokens": 			1500,
			"repetition_penalty": 		1.2,	
			"return_full_text": 		false,
			"seed": 					nil,
			"temperature": 				0.7,
			"top_k": 					5,
			"top_n_tokens": 			5,
			"top_p": 					0.9,
			"truncate": 				nil,
			"typical_p":				0.8,
			"watermark": 				true,
			"length_penalty": 			1.0,
			"early_stopping": 			true,
		},
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
		Timeout: 60 * time.Second,
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