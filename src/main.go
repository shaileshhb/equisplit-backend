package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/log"
	"github.com/shaileshhb/equisplit/src/security"
	"github.com/shaileshhb/equisplit/src/server"
)

func main() {
	logger := log.InitializeLogger()

	// Initialize the database
	database := db.InitDB()
	var wg sync.WaitGroup

	auth := security.NewAuthentication(logger)
	ser := server.NewServer("EquiSplit", database, logger, auth, &wg)
	ser.CreateRouterInstance()
	// db.MigrateTables(ser)
	ser.MigrateTables()
	logger.Error().Err(ser.App.Listen(":8080")).Msg("")

	// Stop Server On System Call or Interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	stopApp(ser)
}

func stopApp(ser *server.Server) {
	// app.Stop()
	ser.WG.Wait()
	fmt.Println("After wait")
	os.Exit(0)
}
