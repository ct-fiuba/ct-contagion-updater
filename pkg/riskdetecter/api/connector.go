package api

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

type InputConnector interface {
	Pop() visits.Visit
}

type OutputConnector interface {
	Push(v1, v2 visits.Visit) bool
}

type ResultConnector interface {
	Push(r Result) bool
}
