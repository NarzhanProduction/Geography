package request_models

type ChatRequest struct {
	Session string `json:"session"`
	Message string `json:"message"`
}
