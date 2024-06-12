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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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

// type StreamResponse struct {
//     ID                string         `json:"id"`
// 	Object            string         `json:"object"`
// 	Created           int64          `json:"created"`
// 	Model             string         `json:"model"`
// 	SystemFingerprint string         `json:"system_fingerprint"`
// 	Choices           []StreamChoice `json:"choices"`
// }

// type StreamChoice struct {
// 	Index           int       `json:"index"`
// 	Message         Message   `json:"delta"`
//     Logprobs        *LogProbs `json:"logprobs"`
//     FinishReason    string   `json:"finish_reason"`
// }

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

// func PromptAIStream(ctx context.Context, inputs string, model string, responseChan chan string) (string, int, string, error) {
// 	requestBody := Request{
// 		FrequencyPenalty: 1,
// 		MaxTokens:        2048,
// 		Messages: []map[string]string{
// 			{
// 				"role":    "system",
// 				"content": SYSTEM_CONTENT,
// 			},
// 			{
// 				"role":    "user",
// 				"content": inputs,
// 			},
// 		},
// 		Model:             model,
// 		Stream:            true,
// 		RepetitionPenalty: 1.2,
// 		Temperature:       0.6,
// 		DoSample:          true,
// 		TopP:              0.95,
// 		Watermark:         false,
// 		LengthPenalty:     1,
// 		EarlyStopping:     true,
// 	}

// 	var (
// 		completeMessage string
// 		numOfTokens     int
// 		finishReason    string
// 	)

// 	body, err := json.Marshal(requestBody)
// 	if err != nil {
// 		return completeMessage, numOfTokens, finishReason, fmt.Errorf("failed to marshal request body: %v", err)
// 	}

// 	req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(body))
// 	if err != nil {
// 		return completeMessage, numOfTokens, finishReason, fmt.Errorf("failed to create request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	client := &http.Client{
// 		Timeout: 0, // No timeout for streaming requests
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return completeMessage, numOfTokens, finishReason, fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check if the response is of type text/event-stream
// 	contentType := resp.Header.Get("Content-Type")
// 	if contentType != "text/event-stream" {
// 		return completeMessage, numOfTokens, finishReason, fmt.Errorf("unexpected Content-Type: %s", contentType)
// 	}

// 	// Use bufio.Reader to read the response line by line
// 	reader := bufio.NewReader(resp.Body)

// 	for {
// 		line, err := reader.ReadBytes('\n')
// 		if err != nil && err != io.EOF {
// 			return completeMessage, numOfTokens, finishReason, fmt.Errorf("failed to read response: %v", err)
// 		}

// 		// Check for Heartbeat or comment lines
// 		if bytes.HasPrefix(line, []byte(":")) {
// 			continue
// 		}

// 		// Trim newline and carriage return characters
// 		line = bytes.TrimSpace(line)

// 		// Check for empty lines
// 		if len(line) == 0 {
// 			continue
// 		}

// 		// Process the SSE message
// 		fields := bytes.SplitN(line, []byte(":"), 2)
// 		if len(fields) < 2 {
// 			return completeMessage, numOfTokens, finishReason, fmt.Errorf("malformed SSE message: %s", string(line))
// 		}

// 		event := string(fields[0])
// 		data := fields[1]

// 		switch event {
// 		case "data":
// 			// Process the data payload
// 			var chunk Response
// 			if err := json.Unmarshal(data, &chunk); err != nil {
// 				return completeMessage, numOfTokens, finishReason, fmt.Errorf("failed to decode JSON response: %v", err)
// 			}

// 			// Send each chunk's content to the response channel
// 			content := chunk.Choices[0].Message.Content
// 			responseChan <- content

// 			// Build complete message
// 			completeMessage += content

// 			// Update num of tokens
// 			numOfTokens += chunk.Usage.CompletionTokens

// 			// Update finish reason
// 			if chunk.Choices[0].FinishReason != "" {
// 				finishReason = chunk.Choices[0].FinishReason
// 				fmt.Println(completeMessage, numOfTokens, finishReason)
// 				return completeMessage, numOfTokens, finishReason, nil
// 			}

// 		case "ping":
// 			// Ignore ping messages
// 			continue

// 		default:
// 			// Unknown event type
// 			return completeMessage, numOfTokens, finishReason, fmt.Errorf("unknown event type: %s", event)
// 		}
// 	}
// }