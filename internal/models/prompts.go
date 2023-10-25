package models

import (
	"database/sql"
	"errors"
	"time"
)

type Prompt struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type PromptModel struct {
	DB *sql.DB // Databse handle
}

func (m *PromptModel) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO prompts (title, content, created, expires) 
			  VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(query, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *PromptModel) Get(id int) (Prompt, error) {
	query := `SELECT id, title, content, created, expires FROM prompts
			  WHERE expires > UTC_TIMESTAMP() AND id = ?`

	var p Prompt

	err := m.DB.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Prompt{}, ErrNoRecord
		} else {
			return Prompt{}, err
		}
	}

	return p, nil
}

func (m *PromptModel) Latest() ([]Prompt, error) {
	return nil, nil
}
