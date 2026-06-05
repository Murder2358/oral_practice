package service

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"oral_practice/pkg/llm"
)

type ConversationService struct {
	llmClient *llm.Client
}

func NewConversationService(llmClient *llm.Client) *ConversationService {
	return &ConversationService{llmClient: llmClient}
}

func (s *ConversationService) HandleConnection(conn *websocket.Conn, sessionID, sceneID, difficulty string) {
	systemPrompt := s.buildSystemPrompt(sceneID, difficulty)
	messages := []llm.Message{
		{Role: "system", Content: systemPrompt},
	}

	// 发送开场白
	opening, err := s.llmClient.Chat(messages)
	if err != nil {
		log.Printf("LLM opening error: %v", err)
		return
	}
	messages = append(messages, llm.Message{Role: "assistant", Content: opening})

	resp := map[string]interface{}{
		"type":    "assistant",
		"content": opening,
	}
	respBytes, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, respBytes)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Message parse error: %v", err)
			continue
		}

		userText, _ := msg["content"].(string)
		if userText == "" {
			continue
		}

		messages = append(messages, llm.Message{Role: "user", Content: userText})

		// 保持上下文窗口在 10 轮以内
		if len(messages) > 21 {
			messages = append(messages[:1], messages[len(messages)-20:]...)
		}

		reply, err := s.llmClient.Chat(messages)
		if err != nil {
			log.Printf("LLM error: %v", err)
			errResp, _ := json.Marshal(map[string]interface{}{
				"type":    "error",
				"content": "AI service error, please try again",
			})
			conn.WriteMessage(websocket.TextMessage, errResp)
			continue
		}

		messages = append(messages, llm.Message{Role: "assistant", Content: reply})

		respBytes, _ = json.Marshal(map[string]interface{}{
			"type":    "assistant",
			"content": reply,
		})
		conn.WriteMessage(websocket.TextMessage, respBytes)
	}
}

func (s *ConversationService) buildSystemPrompt(sceneID, difficulty string) string {
	difficultyDesc := map[string]string{
		"beginner":     "beginner level (simple vocabulary, short sentences)",
		"intermediate": "intermediate level (moderate vocabulary, compound sentences)",
		"advanced":     "advanced level (rich vocabulary, complex sentences, idioms)",
	}

	level := difficultyDesc[difficulty]
	if level == "" {
		level = difficultyDesc["intermediate"]
	}

	scenePrompts := map[string]string{
		"interview": "You are an English interviewer conducting a job interview. Ask questions about self-introduction, work experience, behavioral questions, and technical topics. Be professional and encouraging.",
		"restaurant": "You are a waiter/waitress at an English restaurant. Help the customer with seating, menu recommendations, taking orders, handling special dietary needs, and billing. Be friendly and helpful.",
		"meeting":    "You are a colleague in a business meeting. Discuss project updates, share opinions, handle disagreements professionally, and summarize action items. Use professional business English.",
		"social":     "You are a friendly English speaker at a social gathering. Engage in casual conversation about hobbies, interests, current events, and make plans together. Be warm and approachable.",
		"travel":     "You are a hotel receptionist/tour guide/travel assistant. Help the traveler with check-in, directions, recommendations, and handle any travel issues. Be helpful and patient.",
		"free":       "You are a friendly English conversation partner. Discuss any topic the user wants to talk about. Help them practice natural spoken English.",
	}

	prompt, ok := scenePrompts[sceneID]
	if !ok {
		prompt = scenePrompts["free"]
	}

	return prompt + "\n\nThe user is at " + level + ". Adjust your language complexity accordingly.\nRules:\n- Stay in character throughout the conversation\n- Keep responses concise (2-3 sentences max)\n- If the user makes grammar mistakes, naturally use the correct form in your response\n- Encourage the user to speak more\n- Start the conversation with an appropriate greeting"
}
