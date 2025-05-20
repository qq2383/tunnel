package main

import "C"
import (
	"log"
	"path/filepath"

	"github.com/qq2383/tunnel/client"
	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/logger"
)

var (
	tunnelClient *client.Client = client.New()
	root string
)

//export Init
func Init(r *C.char, openLog C.int) {
	root = filepath.FromSlash(C.GoString(r))
	if openLog != 0 {
		logger.New(root, "tunnel", log.LstdFlags|log.Lshortfile|log.LUTC)
	}
}

//export Start
func Start() {
	fp := filepath.Join(root, "tunnel.cnf")
	var cnf client.Config
	config.Add(&cnf)
	err := config.Read(fp)
	if err != nil {
		log.Println(err)
	}

	err = tunnelClient.Start()
	if err != nil {
		log.Println(err)
	}
}

//export Stop
func Stop() {
	err := tunnelClient.Stop()
	if err != nil {
		log.Println(err)
	}
}

//export State
func State() C.int {
	if client.State {
		return C.int(1)
	}
	return C.int(0)
}

func main() {}
