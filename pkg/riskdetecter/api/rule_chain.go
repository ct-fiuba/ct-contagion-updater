package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

type RuleChain interface {
	AddFilter(id string, rule rules.Rule) bool
	RemoveFilter(id string) bool
	Process(v1, v2 *visits.Visit, s *spaces.Space) error
}
