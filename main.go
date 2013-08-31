package main

import (
	goflags "github.com/jessevdk/go-flags"
	"github.com/politician/goose/eventdb"
	"github.com/politician/goose/web"
	"github.com/politician/goose/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Options struct {
	RedisAddr string `short:"r" long:"redis" value-name:"HOST" description:"Redis address to emit events to" optional:"true"`
	RedisDb   int    `long:"redis-db" value-name:"DB" default:"0" description:"Redis database to use" optional:"true"`

	ep eventdb.EventProvider
}

func main() {
	// Parse the command line.
	opts := parseArgs()

	log.Println("starting....")

	forever(opts)

	log.Println("goodbye...")
}

func forever(opts *Options) {
	// Start a worker to coordiate access to the watchdb and redis event log.
	mgr, err := worker.Start(opts.ep)
	if err != nil {
		log.Fatal(err)
	}

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

// Parses the command-line arguments, and validates them.
func parseArgs() *Options {
	opts := &Options{}

	_, err := goflags.Parse(opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.RedisDb < 0 || opts.RedisDb > 15 {
		log.Fatal("redis db out of range")
	}

	opts.ep = eventdb.NopEventProvider()
	if opts.RedisAddr != "" {
		opts.ep = eventdb.RedisEventProvider("tcp", opts.RedisAddr, opts.RedisDb)
	}

	return opts
}
