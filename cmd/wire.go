//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"risk-management/internal/adapter/kafka/in"
	"risk-management/internal/adapter/kafka/out"
	"risk-management/internal/domain"
	"risk-management/internal/infra/kafka"
	"risk-management/internal/infra/logger"
)

func buildAppContainer() (manager *Manager, err error) {
	wire.Build(kafka.NewSaramaConfig,
		kafka.NewCloudEventsKafkaSender,
		kafka.NewCloudEventsKafkaConsumer,
		logger.NewLogger,
		out.NewKafkaRiskAnalysisRepository,
		wire.Bind(new(domain.RiskAnalysisRepository), new(*out.KafkaRiskAnalysisRepository)),
		domain.NewRiskAnalysisService,
		in.NewFraudScoringReceiver,
		NewManager,
	)
	return nil, nil
}
