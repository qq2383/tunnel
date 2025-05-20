package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/qq2383/socks5"
	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/user"
)

type Server struct {
	net.Listener
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	cnf := config.Get("Server").(*Config)
	pem := filepath.Join(cnf.CertPath, "cert.pem")
	key := filepath.Join(cnf.CertPath, "cert.key")
	cert, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		return err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	addr := fmt.Sprintf(":%d", cnf.Port)
	s.Listener, err = tls.Listen("tcp", addr, config)
	if err != nil {
		return err
	}

	var conn net.Conn
	for {
		conn, err = s.Listener.Accept()
		if err != nil {
			_, ok := err.(*net.OpError)
			if ok {
				break
			}
			continue
		}
		go s.handle(conn)
	}
	return err
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	log.Printf("remote addr: %v\n", conn.RemoteAddr())

	s5 := socks5.New(conn)
	methods, err := s5.AuthGetMethods()
	log.Printf("methods: %v\n", methods)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("AuthRepies\n")
	s5.AuthRepies(socks5.AuthUser)

	uname, upwd, err := s5.AuthGetUser()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("u: %s, %s\n", uname, upwd)

	pass := user.CheckUser(uname, upwd)
	log.Printf("p: %v\n", pass)
	// check user
	if !user.CheckUser(uname, upwd) {
		// auth fail
		s5.AuthRepies(socks5.ReplyFail)
		log.Println("user verify failure")
		return
	}
	// auth success
	s5.AuthRepies(socks5.ReplyPass)

	// get remote addr or port
	host, port, atyp, err := s5.Requests()
	log.Printf("h: %s, %d, %v\n", host, port, atyp)
	if err != nil {
		log.Printf("host: %s, port: %d, %v\n", host, port, err)
		return
	}

	// dial to remote
	server, err := s5.Dial(host, port)
	if err != nil {
		log.Printf("host: %s, port: %d, %v\n", host, port, err)
		s5.Replies(socks5.ReplyFail, socks5.IPv4, nil)
		return
	}

	s5.Replies(socks5.ReplyPass, atyp, server.RemoteAddr())
	s5.Forward(server, conn)
}

func (s *Server) Stop() error {
	return s.Listener.Close()
}

func (s *Server) Close() error {
	return s.Stop()
}
