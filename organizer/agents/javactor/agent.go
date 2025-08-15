package javactor

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

var (
	//go:embed prompt_template.md
	promptTemplate string

	thinkingBudget int32 = 8000
)

const (
	// the job is complicated, require at least 2.5 flash to work on.
	model string = "gemini-2.5-flash"
)

func Run(apiKey string, name string) ([]string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(promptTemplate, name)

	resp, err := client.Models.GenerateContent(ctx, model, genai.Text(prompt), &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget,
		},
	})
	if err != nil {
		return nil, err
	}

	s := resp.Text()
	s = strings.ReplaceAll(s, "```json", "")
	s = strings.ReplaceAll(s, "```", "")

	names := []string{}
	if err := json.Unmarshal([]byte(s), &names); err != nil {
		return nil, err
	}

	return names, nil
}
