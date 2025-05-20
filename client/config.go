package client

import (
	"github.com/qq2383/tunnel/user"
)

type Config struct {
	Local  Local  `yaml:"local"`
	Remote Remote `yaml:"remote"`
	User   []user.User `yaml:"user"`
}

type Local struct {
	Port     int    `yaml:"port"`
}

type Remote struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
	CertPath string `yaml:"cert_path"`
}