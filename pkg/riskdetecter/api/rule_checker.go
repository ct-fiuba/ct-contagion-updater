package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

type RuleChecker interface {
	Process(compromised, infected *visits.Visit, s *spaces.Space) error
	SetNext(RuleChecker)
	SetResultExit(ResultConnector)
}
