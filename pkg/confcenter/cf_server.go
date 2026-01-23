package confcenter

import (
	"store/pkg/middlewares/tracer"
	"time"
)

type Server struct {
	Log    Log
	Tracer tracer.Config
	Http   *ServerConfig
	Grpc   *ServerConfig
}

type Meta struct {
	Namespace string
	Appname   string
	Domain    string
}

type Log struct {
	Level string
	Alarm struct {
		Threshold int
		Interval  int
	}
}

type ServerConfig struct {
	Network string
	Addr    string
	Timeout time.Duration
}
