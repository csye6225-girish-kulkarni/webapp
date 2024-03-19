package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	routerPkg "webapp/router"
)

func main() {
	// Open the log file
	logFile, err := os.OpenFile("/var/log/webapp/webapp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure zerolog to write to the log file
	log.Logger = log.Output(logFile)

	router := routerPkg.InitializeRouter()
	err = router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting the server")
		return
	}
}
