package api

import "github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"

type Result struct {
	CompromisedCode compromisedCodes.CompromisedCode
	Error           error
}
