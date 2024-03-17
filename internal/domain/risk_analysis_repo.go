package domain

type RiskAnalysisRepository interface {
	Store(analysis *RiskAnalysis) error
}
