package out

import (
	"go.uber.org/zap"
	"risk-management/internal/domain"
	"risk-management/internal/infra/kafka"
)

type KafkaRiskAnalysisRepository struct {
	cli kafka.CloudEventsSender
	log *zap.Logger
}

func (krar *KafkaRiskAnalysisRepository) Store(analysis *domain.RiskAnalysis) error {

	return nil
}

func NewKafkaRiskAnalysisRepository(cli kafka.CloudEventsSender, log *zap.Logger) *KafkaRiskAnalysisRepository {
	return &KafkaRiskAnalysisRepository{cli: cli, log: log}
}
