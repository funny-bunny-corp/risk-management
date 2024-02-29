package in

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
)

type FraudScoringReceiver struct {
	log *zap.Logger
}

func (fsr *FraudScoringReceiver) Handle(ctx context.Context, event cloudevents.Event) error {

	return nil
}

func NewFraudScoringReceiver(log *zap.Logger) *FraudScoringReceiver {
	return &FraudScoringReceiver{log: log}
}
