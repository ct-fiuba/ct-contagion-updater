package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

type RuleChecker interface {
	Execute(v1, v2 visits.Visit)
	SetInput(ic InputConnector)
	SetOutput(oc OutputConnector)
	SetResultExit(rc ResultConnector)
}
