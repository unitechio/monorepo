package ai

import "fmt"

type EvaluationResult struct {
	Score          float64 `json:"score"`
	Feedback       string  `json:"feedback"`
	ImprovedAnswer string  `json:"improved_answer"`
}

type LLMService interface {
	EvaluateSpeaking(transcript string) (*EvaluationResult, error)
	EvaluateWriting(prompt, text string) (*EvaluationResult, error)
}

type llmService struct{}

func NewLLMService() LLMService { return &llmService{} }

func (s *llmService) EvaluateSpeaking(transcript string) (*EvaluationResult, error) {
	if transcript == "" {
		return nil, fmt.Errorf("empty transcript")
	}
	return &EvaluationResult{
		Score:          7.5,
		Feedback:       "Good fluency. Avoid filler words. Grammar is generally accurate.",
		ImprovedAnswer: "Technology has fundamentally revolutionized our communication methods.",
	}, nil
}

func (s *llmService) EvaluateWriting(prompt, text string) (*EvaluationResult, error) {
	if text == "" {
		return nil, fmt.Errorf("empty response")
	}
	return &EvaluationResult{
		Score:          7.0,
		Feedback:       "Well-structured essay. Strengthen your argument in paragraph 2 with more specific evidence.",
		ImprovedAnswer: "",
	}, nil
}
