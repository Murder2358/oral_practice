package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50" json:"username"`
	Nickname  string    `gorm:"size:100" json:"nickname"`
	Level     string    `gorm:"size:20;default:intermediate" json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	SceneID    string    `gorm:"size:50" json:"scene_id"`
	Difficulty string    `gorm:"size:20" json:"difficulty"`
	Score      float64   `json:"score"`
	Summary    string    `gorm:"type:text" json:"summary"`
	Status     string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	EndedAt    *time.Time `json:"ended_at"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index" json:"session_id"`
	Role      string    `gorm:"size:20" json:"role"` // user / assistant / system
	Content   string    `gorm:"type:text" json:"content"`
	AudioURL  string    `gorm:"size:500" json:"audio_url"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
