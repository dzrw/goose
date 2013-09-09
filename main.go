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
	Verbose   bool   `short:"v" long:"verbose" default:"false" optional:"true"`

	verbosity int
}

func main() {
	// Parse the command line.
	opts := parseArgs()

	run(opts)

	println("\ngoodbye")
}

func run(opts *Options) {
	// Dial the selected event backend.
	ep, err := dialEventProvider(opts)
	if err != nil {
		log.Fatal(err)
	}

	defer ep.Close()

	// Start a WatchService which coordiates access to watchdb.
	ws, err := worker.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Stop()

	// Start the HTTP servers.
	for _, conf := range servers(ws, ep, opts) {
		srv, err := conf.Start()
		if err != nil {
			log.Fatal(err)
		}

		defer srv.Stop()
	}

	// Wait for OS signals.
	await(opts)
}

// Gets the list of servers to start
func servers(ws worker.WatchService, ep eventdb.EventProvider, opts *Options) []*ServerConf {
	return []*ServerConf{
		NewServerConf(":8080", "watchapi", opts.verbosity,
			web.NewWatchHandler(ws)),
		NewServerConf(":8081", "data-access-service", opts.verbosity,
			web.NewMatchHandler("data-access-service", ws, ep)),
	}
}

// Blocks until SIGINT or SIGTERM.
func await(opts *Options) {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	println("CTRL-C to exit...")

	// Block until we receive a signal.
	sig := <-ch

	if opts.verbosity > 0 {
		log.Println("Got signal: ", sig.String())
	}
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

	opts.verbosity = 0
	if opts.Verbose {
		opts.verbosity += 1
	}

	return opts
}

func dialEventProvider(opts *Options) (ep eventdb.EventProvider, err error) {
	ep = eventdb.NopEventProvider()
	if opts.RedisAddr != "" {
		ep = eventdb.RedisEventProvider("tcp",
			opts.RedisAddr, opts.RedisDb, opts.verbosity)
	}

	return ep, ep.Dial()
}
