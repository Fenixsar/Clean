package handler

import (
	"fmt"
	"net/http"

	"gitlab.q123123.net/ligmar/console-back/pkg/telegram"

	"gitlab.q123123.net/ligmar/boot"
	"gitlab.q123123.net/ligmar/console-back/pkg/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *auth.Auth
	telegram telegram.Service
}

func NewHandler(services *auth.Auth, telegram telegram.Service) *Handler {
	return &Handler{
		services,
		telegram,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	if boot.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.GET("/", h.index)
	router.GET("/ws", h.ws)

	fmt.Printf("Server is up! Port: %s\n", boot.ConsoleServerPort)

	return router
}

func (h *Handler) index(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
