package services

import "github.com/kida21/telegram-langbot/internal/repositories"

type ProgressService struct {
	repo *repositories.ProgressRepository
}

func NewProgressService(repo *repositories.ProgressRepository) *ProgressService {
	return &ProgressService{repo: repo}
}

func (s *ProgressService) GetUserStats(userID int64) (int, int, float64, error) {
	progress, err := s.repo.GetByUser(userID)
	if err != nil {
		return 0, 0, 0, err
	}

	totalAttempts := 0
	totalCorrect := 0
	for _, p := range progress {
		totalAttempts += p.Attempts
		totalCorrect += p.CorrectCount
	}

	accuracy := 0.0
	if totalAttempts > 0 {
		accuracy = float64(totalCorrect) / float64(totalAttempts) * 100
	}
	return totalAttempts, totalCorrect, accuracy, nil
}