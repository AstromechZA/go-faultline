package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/astromechza/go-faultline"
)

func main() {

	// reconfigure the fault injector
	faultline.ReplaceGlobals(faultline.G().With(
		faultline.Enabled,
		faultline.AuditFunction(func(faultLine *faultline.FaultLine, description string, stack string) {
			log.Printf("faultline: %s at:\n%s", description, stack)
		})),
	)

	// setup a random panic in the future (this has a 50 percent chance of happening in the next hour)
	faultline.G().With(faultline.Possibility(0.5)).PossiblyPanicInTheFuture(time.Hour, context.Background())

	http.DefaultServeMux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("before")

		// randomly delay in the middle here
		faultline.G().PossibleDelay(req.Context())

		// maybe return a 500 error here (who knows!)
		if err := faultline.G().PossibleError(fmt.Errorf("something bad")); err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
		}

		log.Printf("after")
	})
	if err := http.ListenAndServe("0.0.0.0:8080", http.DefaultServeMux); err != nil {
		if err != http.ErrServerClosed {
			log.Printf("failed to listen: %s", err)
			os.Exit(1)
		}
	}
}
