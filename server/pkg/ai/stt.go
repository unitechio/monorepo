package ai

import "fmt"

type SpeechToTextService interface {
	Transcribe(audioData []byte) (string, error)
}

type sttService struct{}

func NewSTTService() SpeechToTextService { return &sttService{} }

func (s *sttService) Transcribe(audioData []byte) (string, error) {
	if len(audioData) == 0 {
		return "", fmt.Errorf("empty audio data")
	}
	return "Technology has fundamentally changed how we communicate. Ten years ago we called each other; now we mostly text or use apps.", nil
}
