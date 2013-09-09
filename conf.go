package main

import (
	"github.com/politician/goose/web"
	"log"
	"net/http"
)

type ServerConf struct {
	name      string
	addr      string
	handler   http.Handler
	verbosity int
}

func NewServerConf(addr, name string, verbosity int, handler http.Handler) *ServerConf {
	return &ServerConf{name, addr, handler, verbosity}
}

func (c *ServerConf) IsWatchApi() bool {
	return c.name == "watchapi"
}

func (c *ServerConf) Start() (srv *web.HttpServer, err error) {
	srv = web.NewHttpServer(c.addr, c.handler)
	srv.ServerName = c.name
	srv.Verbosity = c.verbosity

	if err = srv.Start(); err != nil {
		return
	}

	if c.IsWatchApi() {
		log.Printf("Watch API available at http://%s/watches", c.addr)
	} else {
		log.Printf("%s available at http://%s", c.name, c.addr)
	}

	return
}
