package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/astromechza/go-faultline"
)

func main() {
	faultline.ReplaceGlobals(faultline.G().With(
		faultline.Enabled,
		faultline.AuditFunction(func(faultLine *faultline.FaultLine, description string, stack string) {
			log.Printf("%s at:\n%s", description, stack)
		})),
	)

	http.DefaultServeMux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("before")

		// randomly delay in the middle here
		faultline.G().PossibleDelay(req.Context())

		log.Printf("after")

		// maybe return a 500 error here (who knows!)
		if err := faultline.G().PossibleError(fmt.Errorf("something bad")); err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
		}
	})
	if err := http.ListenAndServe("0.0.0.0:8080", http.DefaultServeMux); err != nil {
		if err != http.ErrServerClosed {
			log.Printf("failed to listen: %s", err)
			os.Exit(1)
		}
	}
}
