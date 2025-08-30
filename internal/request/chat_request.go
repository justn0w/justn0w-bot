package request

type ChatRequest struct {
	Question  string `json:"question"`
	EnableRag bool   `json:"enableRag"`
}
