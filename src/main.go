package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/server"
)

func main() {
	// Initialize the database
	database := db.InitDB()
	var wg sync.WaitGroup

	ser := server.NewServer("EquiSplit", database, &wg)
	ser.CreateRouterInstance()
	// db.MigrateTables(ser)
	ser.MigrateTables()
	log.Fatal(ser.App.Listen(":8080"))

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
