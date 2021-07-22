package impl

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
)

type SimpleResultHandler struct {
}

func NewSimpleResultHandler() api.ResultHandler {
	self := new(SimpleResultHandler)

	return self
}

func (self *SimpleResultHandler) Handle(r api.Result) {
	// NOOP
}

func (self *SimpleResultHandler) Push(r api.Result) bool {
	self.Handle(r)
	return true
}

// -- CASTERS
func (self *SimpleResultHandler) asResultConnector() api.ResultConnector {
	return self
}
