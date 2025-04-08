package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	LoadEnv(".env")
	NewServer(NewConfig())

	select {}
}
