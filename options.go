package faultline

import (
	"math/rand"
	"time"
)

type Option func(fl *FaultLine)

func Disabled(fl *FaultLine) {
	fl.disabled = true
}

func Enabled(fl *FaultLine) {
	fl.disabled = false
}

func DisabledIf(b bool) Option {
	if b {
		return Disabled
	}
	return Enabled
}

func EnabledIf(b bool) Option {
	return DisabledIf(!b)
}

func Possibility(f float64) Option {
	return func(fl *FaultLine) {
		fl.possibleProbability = f
	}
}

func MaximumDelay(duration time.Duration) Option {
	return func(fl *FaultLine) {
		fl.delayMax = duration
	}
}

func ModifiedDelay(f func(float64) float64) Option {
	return func(fl *FaultLine) {
		fl.delayMod = f
	}
}

func randomSource(src rand.Source) Option {
	return func(fl *FaultLine) {
		fl.rand = rand.New(src)
	}
}

func AuditFunction(f func(*FaultLine, string, string)) Option {
	return func(fl *FaultLine) {
		fl.auditFunction = f
	}
}
