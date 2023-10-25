package models

import (
	"database/sql"
	"time"
)

type Prompt struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expired time.Time
}

type PromptModel struct {
	DB *sql.DB
}

func (m *PromptModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *PromptModel) Get(id int) (*Prompt, error) {
	return nil, nil
}

func (m *PromptModel) Latest() ([]Prompt, error) {
	return nil, nil
}
