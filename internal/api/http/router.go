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

	v1 := c.Group("/api/v1")
	v1.POST("/chat", handler.Chat)
	v1.POST("/login", handler.Login)
	//v1.POST("/callback/:orderID", handler.Callback)
}
