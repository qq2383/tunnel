package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/qq2383/tunnel/client"
	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/logger"
)

var (
	tunnelClient *client.Client = client.New()
)

func TestTunnel(t *testing.T) {
	root, _ := os.Getwd()
	root = filepath.FromSlash(root)
	
	logger.New(root, "tunnel", log.LstdFlags|log.Lshortfile|log.LUTC)

	fp := filepath.Join(root, "tunnel.cnf")
	var cnf client.Config
	config.Add(&cnf)
	err := config.Read(fp)
	if err != nil {
		log.Panicln(err)
	}

	err = tunnelClient.Start()
	if err != nil {
		log.Panic(err)
	}
}