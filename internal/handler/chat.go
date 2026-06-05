package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"oral_practice/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	convService *service.ConversationService
}

func NewChatHandler(convService *service.ConversationService) *ChatHandler {
	return &ChatHandler{convService: convService}
}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	sessionID := c.Query("session_id")
	sceneID := c.Query("scene_id")
	difficulty := c.Query("difficulty")

	h.convService.HandleConnection(conn, sessionID, sceneID, difficulty)
}
