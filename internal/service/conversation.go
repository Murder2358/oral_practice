package service

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"oral_practice/internal/model"
	"oral_practice/internal/repository"
	"oral_practice/pkg/llm"
)

type ConversationService struct {
	llmClient   *llm.Client
	repo        *repository.Repository
	sceneLoader *SceneLoader
}

func NewConversationService(llmClient *llm.Client, repo *repository.Repository, sceneLoader *SceneLoader) *ConversationService {
	return &ConversationService{
		llmClient:   llmClient,
		repo:        repo,
		sceneLoader: sceneLoader,
	}
}

func (s *ConversationService) HandleConnection(conn *websocket.Conn, sessionID, sceneID, difficulty string) {
	var session model.Session
	if err := s.repo.First(&session, "id = ?", sessionID).Error; err != nil {
		log.Printf("Session not found: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","content":"Session not found"}`))
		return
	}

	sceneCfg := s.sceneLoader.GetScene(sceneID)
	if sceneCfg == nil {
		sceneCfg = s.sceneLoader.GetScene("free")
	}

	// 构建系统 prompt
	difficultyDesc := map[string]string{
		"beginner":     "beginner level (simple vocabulary, short sentences)",
		"intermediate": "intermediate level (moderate vocabulary, compound sentences)",
		"advanced":     "advanced level (rich vocabulary, complex sentences, idioms)",
	}
	level := difficultyDesc[difficulty]
	if level == "" {
		level = difficultyDesc["intermediate"]
	}

	systemPrompt := sceneCfg.Prompt + "\n\nThe user is at " + level + ". Adjust your language complexity accordingly.\nRules:\n- Stay in character throughout the conversation\n- Keep responses concise (2-3 sentences max)\n- If the user makes grammar mistakes, naturally use the correct form in your response\n- Encourage the user to speak more"

	// 加载历史消息
	var history []model.Message
	s.repo.Where("session_id = ?", session.ID).Order("created_at asc").Find(&history)

	messages := []llm.Message{{Role: "system", Content: systemPrompt}}
	for _, m := range history {
		messages = append(messages, llm.Message{Role: m.Role, Content: m.Content})
	}

	// 新会话发送开场白
	if len(history) == 0 {
		opening := sceneCfg.Greeting
		messages = append(messages, llm.Message{Role: "assistant", Content: opening})
		s.saveMessage(session.ID, "assistant", opening)

		respBytes, _ := json.Marshal(map[string]any{
			"type":    "assistant",
			"content": opening,
		})
		conn.WriteMessage(websocket.TextMessage, respBytes)
	}

	// 对话循环
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var msg map[string]any
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		userText, _ := msg["content"].(string)
		if userText == "" {
			continue
		}

		s.saveMessage(session.ID, "user", userText)
		messages = append(messages, llm.Message{Role: "user", Content: userText})

		// 保持上下文窗口在 10 轮以内
		if len(messages) > 21 {
			messages = append(messages[:1], messages[len(messages)-20:]...)
		}

		reply, err := s.llmClient.Chat(messages)
		if err != nil {
			log.Printf("LLM error: %v", err)
			errResp, _ := json.Marshal(map[string]any{
				"type":    "error",
				"content": "AI service error, please try again",
			})
			conn.WriteMessage(websocket.TextMessage, errResp)
			continue
		}

		messages = append(messages, llm.Message{Role: "assistant", Content: reply})
		s.saveMessage(session.ID, "assistant", reply)

		respBytes, _ := json.Marshal(map[string]any{
			"type":    "assistant",
			"content": reply,
		})
		conn.WriteMessage(websocket.TextMessage, respBytes)
	}
}

func (s *ConversationService) saveMessage(sessionID uint, role, content string) {
	msg := model.Message{
		SessionID: sessionID,
		Role:      role,
		Content:   content,
	}
	if err := s.repo.Create(&msg).Error; err != nil {
		log.Printf("Failed to save message: %v", err)
	}
}
