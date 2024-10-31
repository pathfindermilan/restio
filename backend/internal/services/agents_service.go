package services

type AIService interface {
	ProcessData(input string) (string, error)
}

type aiService struct {
}

func NewAIService() AIService {
	return &aiService{}
}

func (s *aiService) ProcessData(input string) (string, error) {
	return "Processed: " + input, nil
}
