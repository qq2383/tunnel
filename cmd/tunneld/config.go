package main

import (
	"github.com/qq2383/tunnel/https"
	"github.com/qq2383/tunnel/server"
)

type Config struct {
	Server server.Config `yaml:"server"`
	Http   https.Config  `yaml:"http"`
}