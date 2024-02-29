package kafka

import "os"

type SaramaConfig struct {
	Host    string
	Topic   string
	GroupId string
}

func NewSaramaConfig() *SaramaConfig {
	return &SaramaConfig{
		Host:    os.Getenv("KAFKA_HOST"),
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupId: os.Getenv("KAFKA_GROUP_ID"),
	}
}
