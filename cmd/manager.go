package main

import (
	"context"
	"risk-management/internal/adapter/kafka/in"
	"risk-management/internal/infra/kafka"
)

type Manager struct {
	receiver *in.FraudScoringReceiver
	cli      kafka.CloudEventsReceiver
}

func (m *Manager) Start() error {
	err := m.cli.StartReceiver(context.Background(), m.receiver.Handle)
	if err != nil {
		return err
	}
	return nil
}

func NewManager(receiver *in.FraudScoringReceiver, cli kafka.CloudEventsReceiver) *Manager {
	return &Manager{
		receiver: receiver,
		cli:      cli,
	}
}
