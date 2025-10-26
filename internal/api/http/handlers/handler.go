package handlers

import (
	"context"
	"net/http"

	request_models "github.com/NarzhanProduction/Geography/internal/api/models"
	"github.com/NarzhanProduction/Geography/internal/service"
	"github.com/labstack/echo/v4"
)

type handler_srv interface {
	Chat(ctx context.Context, w http.ResponseWriter, r *http.Request, request *request_models.ChatRequest) error
	Login(ctx context.Context, request *request_models.UserLoginRequest) (*request_models.LoginResponse, error)
	//Callback(ctx context.Context, request *service.CallbackOperation) error
}

type Handler struct {
	handler handler_srv
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{handler: service}
}

func (h Handler) Login(c echo.Context) error {
	return c.String(http.StatusOK, "accepted")
}

func (h Handler) Chat(c echo.Context) error {
	var req request_models.ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	return h.handler.Chat(c.Request().Context(), c.Response(), c.Request(), &req)
}
