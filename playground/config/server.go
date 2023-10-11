package config

import "time"

var (
	defaultServerHost     = "localhost"
	defaultServerPort     = ":8080"
	defaultReadTimeout    = time.Second * 10
	defaultWriteTimeout   = time.Second * 120
	defaultMaxHeaderBytes = 1 << 20
)

type Server struct {
	Host string
	Port string
}

func New(host, port string,
	readTimeout, writeTimeout time.Duration,
	maxHeaderBytes int) *Server {
	if host == "" {
		host = defaultServerHost
	}

	if port == "" {
		port = defaultServerPort
	}

	return &Server{
		Host: host,
		Port: port,
	}
}
