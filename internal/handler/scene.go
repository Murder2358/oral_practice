package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oral_practice/internal/model"
	"oral_practice/internal/repository"
)

type SceneHandler struct {
	repo *repository.Repository
}

func NewSceneHandler(repo *repository.Repository) *SceneHandler {
	return &SceneHandler{repo: repo}
}

func (h *SceneHandler) ListScenes(c *gin.Context) {
	scenes := []map[string]interface{}{
		{
			"id":          "interview",
			"name":        "求职面试",
			"description": "模拟英文面试场景，涵盖自我介绍、项目经验、行为面试等",
			"icon":        "briefcase",
			"difficulty":  []string{"beginner", "intermediate", "advanced"},
		},
		{
			"id":          "restaurant",
			"name":        "餐厅点餐",
			"description": "模拟餐厅用餐场景，练习落座、点菜、特殊需求、结账等",
			"icon":        "coffee",
			"difficulty":  []string{"beginner", "intermediate"},
		},
		{
			"id":          "meeting",
			"name":        "商务会议",
			"description": "模拟商务会议场景，练习开场白、议程讨论、意见表达等",
			"icon":        "team",
			"difficulty":  []string{"intermediate", "advanced"},
		},
		{
			"id":          "social",
			"name":        "日常社交",
			"description": "模拟日常社交场景，练习寒暄、兴趣爱好、邀请与邀约等",
			"icon":        "smile",
			"difficulty":  []string{"beginner", "intermediate"},
		},
		{
			"id":          "travel",
			"name":        "旅行出行",
			"description": "模拟旅行场景，练习酒店入住、交通问路、紧急情况处理等",
			"icon":        "global",
			"difficulty":  []string{"beginner", "intermediate", "advanced"},
		},
		{
			"id":          "free",
			"name":        "自由对话",
			"description": "自定义话题，自由练习英语口语",
			"icon":        "message",
			"difficulty":  []string{"beginner", "intermediate", "advanced"},
		},
	}
	c.JSON(http.StatusOK, gin.H{"scenes": scenes})
}

func (h *SceneHandler) CreateSession(c *gin.Context) {
	var req struct {
		SceneID    string `json:"scene_id" binding:"required"`
		Difficulty string `json:"difficulty" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := model.Session{
		SceneID:    req.SceneID,
		Difficulty: req.Difficulty,
		Status:     "active",
	}
	if err := h.repo.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session": session})
}

func (h *SceneHandler) GetSession(c *gin.Context) {
	id := c.Param("id")
	var session model.Session
	if err := h.repo.First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	var messages []model.Message
	h.repo.Where("session_id = ?", session.ID).Order("created_at asc").Find(&messages)

	c.JSON(http.StatusOK, gin.H{
		"session":  session,
		"messages": messages,
	})
}
