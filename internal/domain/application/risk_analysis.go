package application

import (
	"go.uber.org/zap"
	"risk-management/internal/domain"
	"risk-management/internal/domain/repositories"
)

var high = &RiskLevel{
	Name: "High",
	From: 61,
	To:   100,
}

var medium = &RiskLevel{
	Name: "Medium",
	From: 31,
	To:   60,
}

var low = &RiskLevel{
	Name: "Low",
	From: 0,
	To:   30,
}

type RiskAnalysisService struct {
	repo repositories.RiskAnalysisRepository
	log  *zap.Logger
}

func (ra *RiskAnalysisService) Assessment(scoring *domain.ScoringResult) *RiskLevel {
	total := scoring.Score.CurrencyScore.Value() + scoring.Score.SellerScore.Value() + scoring.Score.ValueScore.Value() + scoring.Score.AverageValueScore.Score
	if low.Contains(total) {
		return low
	} else if medium.Contains(total) {
		return medium
	} else if high.Contains(total) {
		return high
	}
	return high
}

type RiskLevel struct {
	Name string
	From int
	To   int
}

func (rl *RiskLevel) Contains(val int) bool {
	return rl.From < val && rl.To > val
}

func NewRiskAnalysisService(repo repositories.RiskAnalysisRepository, log *zap.Logger) *RiskAnalysisService {
	return &RiskAnalysisService{
		repo: repo,
		log:  log,
	}
}
