package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/solomon/ims/internal/domain"
)

type GenerateDescriptionInput struct {
	Title    string `json:"title"`
	Location string `json:"location"`
}

type ParseEventInput struct {
	Text    string `json:"text"`
	Today   string `json:"today"`
}

type ParsedEvent struct {
	Title    string     `json:"title"`
	Location string     `json:"location"`
	StartsAt *time.Time `json:"starts_at"`
	EndsAt   *time.Time `json:"ends_at"`
}

type SuggestTimesInput struct {
	PreferredTime time.Time      `json:"preferred_time"`
	Conflicts     []*domain.Event `json:"conflicts"`
}

type TimeSuggestion struct {
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	Reasoning string    `json:"reasoning"`
}

type AIService interface {
	GenerateDescription(ctx context.Context, input GenerateDescriptionInput) (string, error)
	ParseEvent(ctx context.Context, input ParseEventInput) (*ParsedEvent, error)
	SuggestTimes(ctx context.Context, input SuggestTimesInput) ([]TimeSuggestion, error)
}

type aiService struct {
	apiKey  string
	model   string
	timeout time.Duration
	client  *http.Client
}

func NewAIService(apiKey, model string, timeout time.Duration) AIService {
	return &aiService{
		apiKey:  apiKey,
		model:   model,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

func (s *aiService) callGemini(ctx context.Context, prompt string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not configured")
	}

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		s.model, s.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gemini request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini error %d: %s", resp.StatusCode, string(b))
	}

	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}

func (s *aiService) GenerateDescription(ctx context.Context, input GenerateDescriptionInput) (string, error) {
	prompt := fmt.Sprintf(
		`Write a concise, engaging 2-3 sentence description for an event called "%s" at location "%s".
Be professional and informative. Do not include the event name or location as headers — just the description paragraph.`,
		input.Title, input.Location,
	)

	return s.callGemini(ctx, prompt)
}

func (s *aiService) ParseEvent(ctx context.Context, input ParseEventInput) (*ParsedEvent, error) {
	prompt := fmt.Sprintf(
		`Today is %s. Parse the following natural language event description and extract event details.
Text: "%s"

Respond with ONLY valid JSON in this exact format (no markdown, no explanation):
{
  "title": "event title",
  "location": "location or empty string",
  "starts_at": "2024-01-15T14:00:00Z",
  "ends_at": "2024-01-15T15:00:00Z"
}

Rules:
- Use ISO 8601 UTC format for times
- If no end time is specified, add 1 hour to start time
- If no specific date, use the nearest upcoming occurrence
- If no time specified, use 12:00 UTC`,
		input.Today, input.Text,
	)

	text, err := s.callGemini(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Strip markdown code blocks if present
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result struct {
		Title    string `json:"title"`
		Location string `json:"location"`
		StartsAt string `json:"starts_at"`
		EndsAt   string `json:"ends_at"`
	}

	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	parsed := &ParsedEvent{
		Title:    result.Title,
		Location: result.Location,
	}

	if result.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, result.StartsAt)
		if err == nil {
			parsed.StartsAt = &t
		}
	}
	if result.EndsAt != "" {
		t, err := time.Parse(time.RFC3339, result.EndsAt)
		if err == nil {
			parsed.EndsAt = &t
		}
	}

	return parsed, nil
}

func (s *aiService) SuggestTimes(ctx context.Context, input SuggestTimesInput) ([]TimeSuggestion, error) {
	var conflictLines strings.Builder
	for _, c := range input.Conflicts {
		conflictLines.WriteString(fmt.Sprintf("- %s: %s to %s\n",
			c.Title,
			c.StartsAt.Format("Jan 2, 2006 3:04 PM"),
			c.EndsAt.Format("3:04 PM MST"),
		))
	}

	prompt := fmt.Sprintf(
		`I want to schedule an event at %s but it conflicts with:
%s

Suggest 3 alternative time slots (same day or nearby days). Each should be at least 1 hour long and avoid the conflicts above.

Respond with ONLY valid JSON array (no markdown):
[
  {
    "starts_at": "2024-01-15T16:00:00Z",
    "ends_at": "2024-01-15T17:00:00Z",
    "reasoning": "brief human-readable reason why this slot works"
  }
]`,
		input.PreferredTime.Format("Jan 2, 2006 3:04 PM MST"),
		conflictLines.String(),
	)

	text, err := s.callGemini(ctx, prompt)
	if err != nil {
		return nil, err
	}

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var raw []struct {
		StartsAt  string `json:"starts_at"`
		EndsAt    string `json:"ends_at"`
		Reasoning string `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse AI time suggestions: %w", err)
	}

	suggestions := make([]TimeSuggestion, 0, len(raw))
	for _, r := range raw {
		s, err1 := time.Parse(time.RFC3339, r.StartsAt)
		e, err2 := time.Parse(time.RFC3339, r.EndsAt)
		if err1 != nil || err2 != nil {
			continue
		}
		suggestions = append(suggestions, TimeSuggestion{
			StartsAt:  s,
			EndsAt:    e,
			Reasoning: r.Reasoning,
		})
	}

	return suggestions, nil
}
