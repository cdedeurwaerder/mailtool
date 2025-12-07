package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/cdedeurwaerder/mailtool/internal/business"
	"github.com/cdedeurwaerder/mailtool/internal/implem/analyzer"
	"github.com/cdedeurwaerder/mailtool/internal/implem/provider"
	"github.com/cdedeurwaerder/mailtool/internal/implem/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func main() {

	// param: provider & tenant id
	var (
		prov     = os.Getenv("MAILTOOL_API_PROVIDER")
		tenantID = os.Getenv("MAILTOOL_TENANTID")
		dataDir  = os.Getenv("MAILTOOL_DATADIR")
		pa       = provider.NewInMemApi()
		sig      = make(chan os.Signal, 1)
	)

	if prov == "" || tenantID == "" {
		panic("missing env var MAILTOOL_API_PROVIDER or MAILTOOL_TENANTID")
	}

	if dataDir == "" {
		log.Info().Msg("data dir not set, defaultind to ./data")
		dataDir = "./data"
	}
	repo, err := repository.NewSqliteRepository(dataDir)
	if err != nil {
		panic("unable to start repository " + err.Error())
	}

	// provider should be use to choose an implem, but not with our in mem.

	uid, err := uuid.Parse(tenantID)
	if err != nil {
		panic(fmt.Sprintf("invalid tenant id %s: %s", tenantID, err.Error()))
	}
	srv := business.NewService(pa, uid, &analyzer.DummyAnalyzer{}, repo, runtime.NumCPU(), time.Second*5)

	signal.Notify(sig, os.Interrupt)

	srv.Start()

	<-sig
	log.Info().Msg("Handling interrupt signal")
	srv.Stop()
	os.Exit(1)
}
