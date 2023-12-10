package mocks

import (
	"time"

	"promptbox.tyfacey.net/internal/models"
)

var mockPrompts = models.Prompt{
	ID:      1,
	Title:   "Explain to me like a 5 year-old",
	Content: "What the first law of thermodynamics?",
	Created: time.Now(),
	Expires: time.Now(),
}

type PromptModel struct{}

func (m *PromptModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *PromptModel) Get(id int) (models.Prompt, error) {
	switch id {
	case 1:
		return mockPrompts, nil
	default:
		return models.Prompt{}, models.ErrNoRecord
	}

}

func (m *PromptModel) Latest() ([]models.Prompt, error) {
	return []models.Prompt{mockPrompts}, nil
}
