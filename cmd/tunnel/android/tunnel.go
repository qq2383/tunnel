package android

import (
	"log"
	"path/filepath"

	"github.com/qq2383/tunnel/client"
	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/logger"
)

var (
	tunnelClient *client.Client = client.New()
)

func Init(r string, openLog bool) {
	root := filepath.FromSlash(r)
	if openLog {
		logger.New(root, "tunnel", log.LstdFlags|log.Lshortfile|log.LUTC)
	}

	fp := filepath.Join(root, "tunnel.cnf")
	var cnf client.Config
	config.Add(&cnf)
	err := config.Read(fp)
	if err != nil {
		log.Panicln(err)
	}
}

func Start() {
	err := tunnelClient.Start()
	if err != nil {
		log.Panic(err)
	}
}

func Stop() {
	err := tunnelClient.Stop()
	if err != nil {
		log.Panic(err)
	}
}
