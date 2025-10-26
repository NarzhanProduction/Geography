package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NarzhanProduction/Geography/internal/pkg/logger"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatbotService struct {
	logger logger.Logger
	client openai.Client
}

func NewChatbotService(apiKey string, lgr logger.Logger) *ChatbotService {
	return &ChatbotService{
		logger: lgr,
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

func (s *ChatbotService) StreamResponse(w http.ResponseWriter, r *http.Request, msg string) error {
	ctx := context.Background()
	s.logger.Info(ctx, "starting gpt streaming")

	stream := s.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(msg),
		},
		Model: openai.ChatModelGPT4o,
	})

	defer stream.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		s.logger.Error(ctx, "streaming unsupported by response writer")
		return fmt.Errorf("streaming unsupported by response writer")
	}

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if _, err := fmt.Fprintf(w, "%s", content); err != nil {
				break
			}
			flusher.Flush()
			s.logger.Info(ctx, fmt.Sprintf("got chunk of content: %s", content))
		}
	}

	if err := stream.Err(); err != nil {
		s.logger.Error(ctx, fmt.Sprintf("stream error: %v", err.Error()))
		http.Error(w, "stream error: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
