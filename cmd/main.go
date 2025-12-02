package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/cdedeurwaerder/mailtool/internal/business"
	"github.com/cdedeurwaerder/mailtool/internal/implem/provider"
	"github.com/google/uuid"
)

func main() {

	// param: provider & tenant id
	var (
		pa  = provider.InMemApi{}
		srv = business.NewService(pa, uuid.New(), runtime.NumCPU(), time.Second*10)
		sig = make(chan os.Signal, 1)
	)

	signal.Notify(sig, os.Interrupt)

	srv.Start()

	<-sig
	log.Println("Handling interrupt signal")
	srv.Stop()
	os.Exit(1)
}
