package server

type Config struct {
	Port     int    `yaml:"port"`
	CertPath string `yaml:"cert_path"`
}
