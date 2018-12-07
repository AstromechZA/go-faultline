package faultline

import (
	"fmt"
	"math/rand"
	"time"
)

var globalFaultLine *FaultLine

func init() {
	ReplaceGlobals(new(FaultLine).With(
		randomSource(rand.NewSource(time.Now().UTC().UnixNano())),
		Possibility(0.1),
		AuditFunction(nil),
		MaximumDelay(time.Second),
		ModifiedDelay(nil),
	))
}

// G returns the global fault line object
func G() *FaultLine {
	return globalFaultLine
}

// ReplaceGlobals will replace the core fault line object
func ReplaceGlobals(fl *FaultLine) {
	if fl == nil {
		panic(fmt.Errorf("global faultline cannot be nil"))
	}
	globalFaultLine = fl
}
