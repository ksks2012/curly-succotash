package main

import (
	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/ai"
	"curly-succotash/backend/pkg/setting"
	"fmt"
	"log"
	"strings"
)

var (
	cfg string
)

func init() {
	cfg = "etc/"
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(cfg, ",")...)
	err = s.ReadSection("AI", &global.AISetting)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	aiClient, err := ai.NewGeminiClient()
	if err != nil {
		log.Fatalf("failed to initialize AI client: %s", err)
		return
	}
	defer aiClient.Close()

	log.Println("Gemini AI client initialized successfully")

	prompt := fmt.Sprintf(global.StoryPromptTemplate, "Fantasy Adventure")
	story, err := aiClient.GenerateContent(prompt)
	if err != nil {
		log.Fatalf("failed to generate content: %s", err)
		return
	}
	log.Printf("Generated content: %s", story)

	rolePrompt := fmt.Sprintf(global.RolePrompt, 1, story)
	roleText, err := aiClient.GenerateContent(rolePrompt)
	if err != nil {
		log.Fatalf("failed to generate role text: %s", err)
		return
	}
	log.Printf("Generated role text: %s", roleText)

	eventPrompt := fmt.Sprintf(global.EventPrompt, 1, story)
	eventText, err := aiClient.GenerateContent(eventPrompt)
	if err != nil {
		log.Fatalf("failed to generate event text: %s", err)
		return
	}
	log.Printf("Generated event text: %s", eventText)
}
