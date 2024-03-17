package domain

import (
	"time"
)

const (
	Rejected = iota
	Approved
)

type Status uint8

type RiskAnalysis struct {
	Status      Status              `json:"status"`
	At          time.Time           `json:"at"`
	Level       *RiskLevel          `json:"level"`
	Transaction TransactionAnalysis `json:"transaction"`
}

func (ra *RiskAnalysis) Approved() bool {
	return Approved == ra.Status
}

func (ra *RiskAnalysis) Rejected() bool {
	return Rejected == ra.Status
}

func (s Status) String() string {
	switch s {
	case Approved:
		return "approved"
	case Rejected:
		return "rejected"
	}
	return ""
}
