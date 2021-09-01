package impl

import (
	"log"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	// "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
)

type SimpleResultHandler struct {
}

func NewSimpleResultHandler() *SimpleResultHandler {
	self := new(SimpleResultHandler)

	return self
}

func (self *SimpleResultHandler) Handle(r api.Result) {
	log.Printf("Llegó el resultado! ==> %v \n", r)
}

func (self *SimpleResultHandler) Push(r api.Result) bool {
	self.Handle(r)
	return true
}

// -- CASTERS

func (self *SimpleResultHandler) AsResultConnector() api.ResultConnector {
	return self
}
