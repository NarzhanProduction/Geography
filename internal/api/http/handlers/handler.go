package handlers

import (
	"context"
	"net/http"
	"strconv"

	request_models "github.com/NarzhanProduction/Geography/internal/api/models"
	"github.com/NarzhanProduction/Geography/internal/database/models"
	"github.com/NarzhanProduction/Geography/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler_srv interface {
	Chat(ctx context.Context, w http.ResponseWriter, r *http.Request, request *request_models.ChatRequest) error
	Login(ctx context.Context, request *request_models.UserLoginRequest) (*request_models.LoginResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateArticle(ctx context.Context, request *request_models.ArticleUpdateRequest) error
	GetArticleByID(ctx context.Context, request *request_models.GetArticleRequest) (*request_models.GetArticleResponse, error)
}

type Handler struct {
	handler   handler_srv
	jwtSecret string
}

func NewHandler(service *service.Service, jwtKey string) *Handler {
	return &Handler{
		handler:   service,
		jwtSecret: jwtKey,
	}
}

func (h Handler) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization needed"})
			}

			tokenStr := cookie.Value
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(h.jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			}

			sub, ok := claims["sub"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid subject claim"})
			}

			userID, err := uuid.Parse(sub)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid UUID"})
			}

			user, err := h.handler.GetUserByID(context.Background(), userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
			}

			// Save user into the context
			c.Set("user", user)
			return next(c)
		}
	}
}

func (h Handler) GetArticle(c echo.Context) error {
	var req request_models.GetArticleRequest
	idStr := c.QueryParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	req.ID = id

	article, err := h.handler.GetArticleByID(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, article)
}

func (h Handler) Login(c echo.Context) error {
	var req request_models.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	resp, err := h.handler.Login(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid login or password"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = resp.Token
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode
	cookie.MaxAge = 60 * 60 * 24 // 1 day
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Succesful login"})
}

func (h Handler) GetUser(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	return c.JSON(http.StatusOK, user.IsAdmin)
}

func (h Handler) Chat(c echo.Context) error {
	var req request_models.ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	return h.handler.Chat(c.Request().Context(), c.Response(), c.Request(), &req)
}

func (h Handler) ArticleEdit(c echo.Context) error {
	userData := c.Get("user")
	if userData == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization needed"})
	}

	user := userData.(*models.User)

	if !user.IsAdmin {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can edit this"})
	}

	var update_req request_models.ArticleUpdateRequest

	if err := c.Bind(&update_req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)
	update_req.ID = id

	if err := h.handler.UpdateArticle(c.Request().Context(), &update_req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error saving an edited article"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully edited article"})
}
