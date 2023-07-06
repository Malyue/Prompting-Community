package client

import "time"

type Service struct {
	Id      string        `yaml:"id"`
	Addr    string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
}

type Config struct {
	ServiceList []*Service `yaml:"services"`
}
