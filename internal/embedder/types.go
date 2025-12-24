package embedder

import "net/http"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type embeddingRequest struct {
	Input    []string `json:"input"`
	Model    string   `json:"model"`
	Encoding string   `json:"encoding_format"`
}

type data struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type embeddingResponse struct {
	Model string `json:"model"`
	Data  []data `json:"data"`
	Usage usage  `json:"usage"`
}
