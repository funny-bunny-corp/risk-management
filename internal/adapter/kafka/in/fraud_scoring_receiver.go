package in

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	"risk-management/internal/domain"
	"risk-management/internal/domain/application"
)

const eventType = "paymentic.io.fraud-detection.v1.transaction.scorecard.created"

type FraudScoringReceiver struct {
	log *zap.Logger
	rs  *application.RiskAnalysisService
}

func (fsr *FraudScoringReceiver) Handle(ctx context.Context, event cloudevents.Event) error {
	if eventType == event.Type() {
		data := &domain.ScoringResult{}
		if err := event.DataAs(data); err != nil {
			fsr.log.Error("error to retrieve deserialize cloud event data", zap.String("error", err.Error()))
			return err
		}
		err := fsr.rs.Assessment(data)
		if err != nil {
			fsr.log.Error("error to retrieve deserialize cloud event data", zap.String("error", err.Error()))
			return err
		}
		return nil
	}
	return nil
}

func NewFraudScoringReceiver(log *zap.Logger, rs *application.RiskAnalysisService) *FraudScoringReceiver {
	return &FraudScoringReceiver{log: log, rs: rs}
}
