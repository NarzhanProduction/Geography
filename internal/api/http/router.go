package http

import (
	"github.com/NarzhanProduction/Geography/internal/api/http/handlers"
	"github.com/labstack/echo/v4"
)

func InitRouter() *echo.Echo {
	c := echo.New()
	return c
}

func RegisterHandlers(c *echo.Echo, handler *handlers.Handler) {
	c.GET("/chatbot", func(c echo.Context) error {
		return c.File("static/chatbox.html")
	})
	c.GET("/editor/:articleID", func(c echo.Context) error {
		return c.File("static/editor.html")
	}, handler.JWTAuth())
	c.GET("/login", func(c echo.Context) error {
		return c.File("static/login.html")
	})
	c.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	v1 := c.Group("/api/v1")
	v1.POST("/chat", handler.Chat)
	v1.POST("/login", handler.Login)
	v1.GET("/article", handler.GetArticle)
	v1.GET("/user", handler.GetUser, handler.JWTAuth())
	v1.PUT("/article/:id", handler.ArticleEdit, handler.JWTAuth())
}
