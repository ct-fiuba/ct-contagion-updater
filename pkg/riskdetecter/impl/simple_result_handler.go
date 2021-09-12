package impl

import (
	"fmt"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
)

type SimpleResultHandler struct {
	compromisedCollection *compromisedCodes.CompromisedCodesCollection
}

func NewSimpleResultHandler(compromisedCodesCollection *compromisedCodes.CompromisedCodesCollection) *SimpleResultHandler {
	self := new(SimpleResultHandler)
	self.compromisedCollection = compromisedCodesCollection

	return self
}

func (rh *SimpleResultHandler) Handle(result api.Result) {
	if result.Error != nil {
		fmt.Errorf("[ResultHandler] Error raised in rule chain --> %e\n", result.Error)
		return
	}

	fmt.Printf("[ResultHandler] Writing compromised code (%+v) \n", result.CompromisedCode)
	err := rh.compromisedCollection.Insert(result.CompromisedCode)
	if err != nil {
		fmt.Errorf("[ResultHandler] Error writing compromised code to DB --> %e\n", err)
	}
}

func (rh *SimpleResultHandler) Push(result api.Result) bool {
	rh.Handle(result)
	return true
}

// -- CASTERS

func (rh *SimpleResultHandler) AsResultConnector() api.ResultConnector {
	return rh
}
