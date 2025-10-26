package service

import (
	"context"
	"net/http"

	request_models "github.com/NarzhanProduction/Geography/internal/api/models"
)

type Service struct {
	userSrv *UserService
	gptSrv  *ChatbotService
}

func NewService(usrService *UserService, gptService *ChatbotService) *Service {
	return &Service{
		userSrv: usrService,
		gptSrv:  gptService,
	}
}

func (s *Service) Login(ctx context.Context, request *request_models.UserLoginRequest) (*request_models.LoginResponse, error) {
	return s.userSrv.LoginCheck(ctx, *request)
}

func (s *Service) Chat(ctx context.Context, w http.ResponseWriter, r *http.Request, request *request_models.ChatRequest) error {
	return s.gptSrv.StreamResponse(w, r, request.Message)
}
