package main

import (
	"github.com/politician/goose/eventdb"
	"github.com/politician/goose/web"
	"github.com/politician/goose/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("starting....")

	forever()

	log.Println("goodbye...")
}

func forever() {
	// Setup the event log.
	evt, err := eventdb.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalln("could not establish connection to the event database: ", err)
	}

	defer evt.Close()

	// Start a worker to coordiate access to the watchdb and redis event log.
	mgr := worker.Start(evt)

	defer mgr.Stop()

	// Start the HTTP servers.
	for _, conf := range servers(mgr) {
		srv, err := conf.Start()
		if err != nil {
			log.Fatal(err)
		}

		defer srv.Stop()
	}

	// Wait for OS signals.
	await()
}

// Gets the list of servers to start
func servers(mgr worker.Manager) []*ServerConf {
	watcher := mgr.(worker.WatchList)
	matcher := mgr.(worker.Matcher)

	return []*ServerConf{
		NewServerConf(":8080", "watchapi",
			web.NewWatchHandler(watcher)),
		NewServerConf(":8081", "data-access-service",
			web.NewMatchHandler("data-access-service", matcher)),
	}
}

// Blocks until SIGINT or SIGTERM.
func await() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	log.Println("CTRL-C to exit...")

	// Block until we receive a signal.
	sig := <-ch

	log.Println("Got signal: ", sig.String())
}
