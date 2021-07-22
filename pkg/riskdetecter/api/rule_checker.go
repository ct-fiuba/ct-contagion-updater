package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

type RuleChecker interface {
	Execute(visit visits.Visit)
	SetInput(ic InputConnector)
	SetOutput(oc OutputConnector)
	SetResultExit(rc ResultConnector)
}
