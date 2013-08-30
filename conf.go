package main

import (
	"github.com/politician/goose/web"
	"log"
	"net/http"
)

type ServerConf struct {
	name    string
	addr    string
	handler http.Handler
}

func NewServerConf(addr, name string, handler http.Handler) *ServerConf {
	return &ServerConf{name, addr, handler}
}

func (c *ServerConf) IsWatchApi() bool {
	return c.name == "watchapi"
}

func (c *ServerConf) Start() (srv *web.HttpServer, err error) {
	srv, err = web.StartHttpServer(c.addr, c.handler)
	if err != nil {
		return
	}

	if c.IsWatchApi() {
		log.Printf("Watch API available at http://%s/watches", c.addr)
	} else {
		log.Printf("%s available at http://%s", c.name, c.addr)
	}

	return
}
