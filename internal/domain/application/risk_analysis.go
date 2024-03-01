package application

import (
	"go.uber.org/zap"
	"risk-management/internal/domain"
	"risk-management/internal/domain/repositories"
	"time"
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

func (ra *RiskAnalysisService) Assessment(scoring *domain.ScoringResult) error {
	total := scoring.Score.CurrencyScore.Value() + scoring.Score.SellerScore.Value() + scoring.Score.ValueScore.Value() + scoring.Score.AverageValueScore.Score
	var r *RiskLevel
	if low.Contains(total) {
		r = low
	} else if medium.Contains(total) {
		r = medium
	} else if high.Contains(total) {
		r = high
	}

	var s domain.Status

	if r == medium || r == low {
		s = domain.Approved
	} else {
		s = domain.Rejected
	}

	a := &domain.RiskAnalysis{
		Status: s,
		At:     time.Now(),
		Level:  r,
	}
	err := ra.repo.Store(a)
	if err != nil {
		ra.log.Error("error to store risk analysis", zap.String("error", err.Error()))
	}
	return nil
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
