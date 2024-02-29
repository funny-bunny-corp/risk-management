package main

import (
	"github.com/google/wire"
	"risk-management/internal/adapter/kafka/in"
	"risk-management/internal/adapter/kafka/out"
	"risk-management/internal/domain/application"
	"risk-management/internal/domain/repositories"
	"risk-management/internal/infra/kafka"
	"risk-management/internal/infra/logger"
)

func buildAppContainer() (manager *Manager, err error) {
	wire.Build(kafka.NewSaramaConfig,
		logger.NewLogger,
		out.NewKafkaRiskAnalysisRepository,
		wire.Bind(new(repositories.RiskAnalysisRepository), new(*out.KafkaRiskAnalysisRepository)),
		application.NewRiskAnalysisService,
		in.NewFraudScoringReceiver,
		kafka.NewCloudEventsKafkaSender,
		kafka.NewCloudEventsKafkaConsumer,
		NewManager,
	)
	return nil, nil
}
