package faultline

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

type FaultLine struct {
	disabled bool

	// possibleProbability guides the probability that faults will be injected or not
	possibleProbability float64

	// delayMax controls the maximum sleep time that a Delay could cause
	delayMax time.Duration
	// delayMod function can be used to change the delay duration (in order to create alternative timing distributions)
	delayMod func(float64) float64

	// auditFunction is used to report the fault injections that occur
	auditFunction func(fl *FaultLine, eventDescription string, callStack string)

	rand *rand.Rand
}

// With is used to apply options to the FaultLine. It creates a new copy of the parent object
func (fl *FaultLine) With(options ...Option) *FaultLine {
	other := &FaultLine{
		auditFunction:       fl.auditFunction,
		disabled:            fl.disabled,
		possibleProbability: fl.possibleProbability,
		delayMax:            fl.delayMax,
		delayMod:            fl.delayMod,
		rand:                fl.rand,
	}
	for _, o := range options {
		o(other)
	}
	return other
}

// Possible determines whether the possibility probability evaluated to true
func (fl *FaultLine) Possible() bool {
	return !fl.disabled && fl.rand.Float64() <= fl.possibleProbability
}

func (fl *FaultLine) doDelay(duration time.Duration, ctx context.Context) bool {
	if duration > 0 {
		t := time.NewTimer(duration)
		select {
		case <-t.C:
			return true
		case <-ctx.Done():
			t.Stop()
			return false
		}
	}
	return false
}

func (fl *FaultLine) prepareDelay() time.Duration {
	if fl.disabled {
		return 0
	}
	delayFactor := fl.rand.Float64()
	if fl.delayMod != nil {
		delayFactor = fl.delayMod(delayFactor)
	}
	delayFactor = math.Min(math.Max(0, delayFactor), 1)
	return time.Duration(delayFactor * float64(fl.delayMax))
}

// PossibleDelay will maybe add a small sleep time at the invocation site. Whether to sleep or not, and the exactly
// sleep time, are determined by the configuration of the FaultLine.
func (fl *FaultLine) PossibleDelay(ctx context.Context) bool {
	if fl.disabled {
		return false
	}
	if fl.Possible() {
		d := fl.prepareDelay()
		if fl.auditFunction != nil {
			var buff []byte
			runtime.Stack(buff, false)
			fl.auditFunction(fl, fmt.Sprintf("injecting delay of %s", d), string(buff))
		}
		fl.doDelay(d, ctx)
		return true
	}
	return false
}

// PossibleError will maybe return a specific error here
func (fl *FaultLine) PossibleError(err error) error {
	if fl.disabled {
		return nil
	}
	if fl.Possible() {
		if fl.auditFunction != nil {
			var buff []byte
			runtime.Stack(buff, false)
			fl.auditFunction(fl, fmt.Sprintf("injecting error '%s'", err), string(buff))
		}
		return err
	}
	return nil
}

// PossibleErrorOr is a convenience feature that will return the existing error if not nil
func (fl *FaultLine) PossibleErrorOr(existing error, possible error) error {
	if existing != nil {
		return existing
	}
	return fl.PossibleError(possible)
}

func (fl *FaultLine) DoInTheFuture(before time.Duration, f func(), ctx context.Context) bool {
	if fl.disabled {
		return false
	}
	subFl := fl.With(MaximumDelay(before))
	go func() {
		if subFl.doDelay(subFl.prepareDelay(), ctx) {
			f()
		}
	}()
	return true
}

func (fl *FaultLine) PossiblyDoInTheFuture(before time.Duration, f func(), ctx context.Context) bool {
	if fl.disabled {
		return false
	}
	if fl.Possible() {
		return fl.DoInTheFuture(before, f, ctx)
	}
	return false
}

func (fl *FaultLine) PossiblyPanicInTheFuture(before time.Duration, ctx context.Context) bool {
	return fl.PossiblyDoInTheFuture(before, func() {
		panic(fmt.Errorf("injected a random panic"))
	}, ctx)
}

func (fl *FaultLine) PossiblySignalInTheFuture(signal os.Signal, before time.Duration, ctx context.Context) bool {
	return fl.PossiblyDoInTheFuture(before, func() {
		if p, err := os.FindProcess(os.Getpid()); err == nil {
			p.Signal(signal)
		}
	}, ctx)
}
