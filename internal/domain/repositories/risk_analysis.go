package repositories

import "risk-management/internal/domain"

type RiskAnalysisRepository interface {
	Store(analysis *domain.RiskAnalysis) error
}
