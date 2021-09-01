package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/contagions"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
)

type RuleChain interface {
	AddFilter(id string, rule rules.Rule) bool
	RemoveFilter(id string) bool
	Process(contagion contagions.Contagion)
}
