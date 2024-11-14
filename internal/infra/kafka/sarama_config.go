package kafka

import "os"

type SaramaConfig struct {
	Host                string
	RiskManagementTopic string
	FraudDetectionTopic string
	GroupId             string
}

func NewSaramaConfig() *SaramaConfig {
	return &SaramaConfig{
		Host:                os.Getenv("KAFKA_HOST"),
		RiskManagementTopic: os.Getenv("KAFKA_RISK_MANAGEMENT_TOPIC"),
		FraudDetectionTopic: os.Getenv("KAFKA_FRAUD_DETECTION_TOPIC"),
		GroupId:             os.Getenv("KAFKA_GROUP_ID"),
	}
}
