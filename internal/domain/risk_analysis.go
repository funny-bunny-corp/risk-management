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
	Status Status
	At     time.Time
	Level  *RiskLevel
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
