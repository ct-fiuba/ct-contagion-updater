package api

type ContactSeverity = string

const (
	HighRisk   ContactSeverity = "HIGH"
	MediumRisk ContactSeverity = "MEDIUM"
	LowRisk    ContactSeverity = "LOW"
)

type Result struct {
	Severity ContactSeverity
	Error    error
}
