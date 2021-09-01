package api

type ContactSeverity = string

const (
	HighRisk   ContactSeverity = "High"
	MediumRisk ContactSeverity = "Medium"
	LowRisk    ContactSeverity = "Low"
)

type Result struct {
	Severity ContactSeverity
	Error    error
}
