package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/qq2383/socks5"
	"github.com/qq2383/tunnel/config"
)

var (
	State bool = false
)

type Client struct {
	ln net.Listener
}

func New() *Client {
	return &Client{}
}

func (c *Client) Start() error {
	cnf := config.Get("").(*Config)

	var err error
	c.ln, err = net.Listen("tcp", fmt.Sprintf(":%d", cnf.Local.Port))
	if err != nil {
		log.Println(err)
	}

	State = true
	for {
		conn, err := c.ln.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}

		go NewHandle(conn)
	}
	return nil
}

func (c *Client) Stop() error {
	State = false
	return c.ln.Close()
}

type Handle struct {
}

func NewHandle(conn net.Conn) {
	handle := &Handle{}
	handle.Do(conn)
}

func (h *Handle) Do(conn net.Conn) {
	defer conn.Close()

	s5 := socks5.New(conn)
	methods, err := s5.AuthGetMethods()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("methods: %v\n", methods)

	cnf := config.Get("").(*Config)
	remote, err := h.dial(cnf)
	if err != nil {
		log.Println(err)
		return
	}

	rs5 := socks5.New(remote)
	if err = rs5.AuthSendMethods(socks5.AuthUser); err != nil {
		log.Println(err)
		return
	}

	method, err := rs5.AuthGetMethod()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("r mehtod: %v\n", method)
	if method != socks5.AuthUser {
		s5.AuthRepies(socks5.ReplyFail)
		return
	}
	
	log.Printf("send user\n")
	usr := &cnf.User[0]
	if err = rs5.AuthSendUser(usr.Name, usr.Passwd); err != nil {
		log.Println(err)
		return
	}

	log.Printf("authStatus\n")
	if err := rs5.AuthGetStatus(); err != nil {
		log.Println(err)
		return
	}

	log.Printf("c autho none\n")
	if err = s5.AuthRepies(socks5.AuthNone); err != nil {
		log.Println(err)
		return
	}
	log.Printf("forward\n")
	s5.Forward(remote, conn)
}

func (h *Handle) dial(cnf *Config) (*tls.Conn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	addr := fmt.Sprintf("%s:%d", cnf.Remote.Addr, cnf.Remote.Port)
	remote, err := tls.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed connected server")
	}
	return remote, nil
}
