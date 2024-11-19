package out

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"risk-management/internal/domain"
	"risk-management/internal/infra/kafka"
)

const (
	eventTypeApproved    = "funny-bunny.xyz.risk-management.v1.risk.decision.approved"
	eventSubjectApproved = "risk-decision-approved"
	eventTypeRejected    = "funny-bunny.xyz.risk-management.v1.risk.decision.rejected"
	eventSubjectRejected = "risk-decision-rejected"
	eventSource          = "risk-management"
	eventContextData     = "domain"
	eventAudienceData    = "external-bounded-context"
	eventContextName     = "eventcontext"
	eventAudienceName    = "audience"
)

type KafkaRiskAnalysisRepository struct {
	cli kafka.CloudEventsSender
	log *zap.Logger
}

func (krar *KafkaRiskAnalysisRepository) Store(analysis *domain.RiskAnalysis) error {
	var evt cloudevents.Event
	if analysis.Approved() {
		evt = approvedEvent(analysis)
	} else if analysis.Rejected() {
		evt = rejectedEvent(analysis)
	}
	if result := krar.cli.Send(
		kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(evt.ID())),
		evt,
	); cloudevents.IsUndelivered(result) {
		krar.log.Error("failed to send", zap.String("error", result.Error()))
	} else {
		krar.log.Info("message sent", zap.String("id", evt.ID()), zap.Bool("ack", cloudevents.IsACK(result)))
	}
	return nil
}

func approvedEvent(analysis *domain.RiskAnalysis) cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetType(eventTypeApproved)
	e.SetSource(eventSource)
	e.SetSubject(eventSubjectApproved)
	_ = e.SetData(cloudevents.ApplicationJSON, analysis)
	return e
}

func rejectedEvent(analysis *domain.RiskAnalysis) cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetType(eventTypeRejected)
	e.SetSource(eventSource)
	e.SetExtension(eventAudienceName, eventAudienceData)
	e.SetExtension(eventContextName, eventContextData)
	e.SetSubject(eventSubjectRejected)
	_ = e.SetData(cloudevents.ApplicationJSON, analysis)
	return e
}

func NewKafkaRiskAnalysisRepository(cli kafka.CloudEventsSender, log *zap.Logger) *KafkaRiskAnalysisRepository {
	return &KafkaRiskAnalysisRepository{cli: cli, log: log}
}
