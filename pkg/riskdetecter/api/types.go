package api

import "github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"

type ContactSeverity = int

const (
	HighRisk   ContactSeverity = 0
	MediumRisk ContactSeverity = 1
	LowRisk    ContactSeverity = 2
)

var RiskStringsToSeverity map[string]ContactSeverity = map[string]ContactSeverity{
	"Alto":  HighRisk,
	"Medio": MediumRisk,
	"Bajo":  LowRisk,
}

type Result struct {
	CompromisedCode compromisedCodes.CompromisedCode
	Error           error
}
