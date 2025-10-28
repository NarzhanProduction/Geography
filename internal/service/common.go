package service

import (
	"context"
	"net/http"

	request_models "github.com/NarzhanProduction/Geography/internal/api/models"
	"github.com/NarzhanProduction/Geography/internal/database/models"
	"github.com/google/uuid"
)

type Service struct {
	userSrv    *UserService
	gptSrv     *ChatbotService
	articleSrv *ArticleService
}

func NewService(usrService *UserService, gptService *ChatbotService, articleService *ArticleService) *Service {
	return &Service{
		userSrv:    usrService,
		gptSrv:     gptService,
		articleSrv: articleService,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userSrv.repo.GetByID(ctx, id)
}

func (s *Service) Login(ctx context.Context, request *request_models.UserLoginRequest) (*request_models.LoginResponse, error) {
	return s.userSrv.LoginCheck(ctx, *request)
}

func (s *Service) Chat(ctx context.Context, w http.ResponseWriter, r *http.Request, request *request_models.ChatRequest) error {
	return s.gptSrv.StreamResponse(w, r, request.Message)
}

func (s *Service) GetArticleByID(ctx context.Context, request *request_models.GetArticleRequest) (*request_models.GetArticleResponse, error) {
	return s.articleSrv.GetArticle(ctx, request.ID)
}

func (s *Service) UpdateArticle(ctx context.Context, request *request_models.ArticleUpdateRequest) error {
	return s.articleSrv.UpdateArticle(ctx, request.ID, request.HTML)
}
