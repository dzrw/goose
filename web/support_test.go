package web

import (
	_ "github.com/davecgh/go-spew/spew"
	"log"
	"net"
	"net/http"
)

type fakeHandler struct{}

func (d *fakeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "actually not an error", 200)
}

func StartHttpServer(addr string, handler http.Handler) (srv *HttpServer, err error) {
	srv = NewHttpServer(addr, handler)
	err = srv.Start()
	return
}

func CheckConnection(addr string) (err error) {
	c := NewHttpTestClient(addr)
	err = c.Dial()
	if err != nil {
		return err
	}

	c.conn.Close()
	return
}

// A minimal HTTP client.
type HttpTestClient struct {
	Addr string
	//Rate int

	conn net.Conn
	quit chan chan bool
}

func NewHttpTestClient(addr string) *HttpTestClient {
	return &HttpTestClient{
		Addr: addr,
		quit: make(chan chan bool),
	}
}

func (c *HttpTestClient) Dial() (err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.Addr)
	if err != nil {
		return
	}

	log.Printf("dialing: %+v", tcpAddr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	c.conn = conn
	return
}

func (c *HttpTestClient) Start() {
	go c.forever()
}

func (c *HttpTestClient) forever() {
	// const (
	//  REQUEST = "ping"
	// )

	defer c.conn.Close()

	//recv, fault := Pump(c.conn, 32)

	for {
		select {
		// case buf := <-recv:
		//  spew.Dump(buf)

		// case err := <-fault:
		//  log.Println("client read error: ", err)
		//  break

		// case <-time.After(time.Duration(rate) * time.Second):
		//  _, err := c.conn.Write([]byte(REQUEST))
		//  if err != nil {
		//    log.Println("client write error: ", err)
		//    break
		//  }

		case done := <-c.quit:
			done <- true
			break
		}
	}
}

func (c *HttpTestClient) Stop() {
	log.Println("stopping client")
	done := make(chan bool)
	c.quit <- done
	<-done
}

// ---------------------------------
