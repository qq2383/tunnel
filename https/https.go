package https

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/https/html"
)

type HttpServer struct {
	*http.Server
	ln net.Listener
}

func New() *HttpServer {
	http.Handle("/", &Handler{root: html.FS})
	return &HttpServer{Server: &http.Server{}}
}

func (hs *HttpServer) Close() error {
	return hs.Server.Close()
}

func (hs *HttpServer) Start() error {
	cnf := config.Get("Http").(*Config)
	var err error
	hs.ln, err = net.Listen("tcp", fmt.Sprintf(":%d", cnf.Port))
	if err != nil {
		return err
	}

	err = hs.Serve(hs.ln)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (hs *HttpServer) Stop() error {
	return hs.ln.Close()
}
