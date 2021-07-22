package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

type RuleChain interface {
	AddFilter(id string, filter RuleChecker) bool
	RemoveFilter(id string) bool
	Process(visit visits.Visit)
}
