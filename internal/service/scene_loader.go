package service

import (
	"encoding/json"
	"os"
)

type SceneConfig struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Difficulty  []string `json:"difficulty"`
	Prompt      string   `json:"prompt"`
	Greeting    string   `json:"greeting"`
}

type SceneLoader struct {
	scenes []SceneConfig
}

func NewSceneLoader(path string) *SceneLoader {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to load scenes config: " + err.Error())
	}

	var cfg struct {
		Scenes []SceneConfig `json:"scenes"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic("Failed to parse scenes config: " + err.Error())
	}

	return &SceneLoader{scenes: cfg.Scenes}
}

func (l *SceneLoader) GetScene(sceneID string) *SceneConfig {
	for i := range l.scenes {
		if l.scenes[i].ID == sceneID {
			return &l.scenes[i]
		}
	}
	return nil
}

func (l *SceneLoader) AllScenes() []map[string]any {
	items := make([]map[string]any, 0, len(l.scenes))
	for _, s := range l.scenes {
		items = append(items, map[string]any{
			"id":          s.ID,
			"name":        s.Name,
			"description": s.Description,
			"icon":        s.Icon,
			"difficulty":  s.Difficulty,
		})
	}
	return items
}
