package main

import "promptbox.tyfacey.net/internal/models"

type templateData struct {
	Prompt  models.Prompt
	Prompts []models.Prompt
}
