package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oral_practice/internal/model"
	"oral_practice/internal/repository"
	"oral_practice/internal/service"
)

type SceneHandler struct {
	repo        *repository.Repository
	sceneLoader *service.SceneLoader
}

func NewSceneHandler(repo *repository.Repository, sceneLoader *service.SceneLoader) *SceneHandler {
	return &SceneHandler{repo: repo, sceneLoader: sceneLoader}
}

func (h *SceneHandler) ListScenes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"scenes": h.sceneLoader.AllScenes()})
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

func (h *SceneHandler) EndSession(c *gin.Context) {
	id := c.Param("id")
	var session model.Session
	if err := h.repo.First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	session.Status = "ended"
	h.repo.Save(&session)

	c.JSON(http.StatusOK, gin.H{"session": session})
}
