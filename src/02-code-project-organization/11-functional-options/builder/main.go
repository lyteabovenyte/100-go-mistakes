package main

import (
	"errors"
	"net/http"
)

const defaultHTTPPort = 8080

type Config struct {
	Port int
}

type ConfigBuilder struct {
	port *int
}

// such a configuration method returns the builder itself so that
// we can use method chaining
func (b *ConfigBuilder) Port(port int) *ConfigBuilder {
	b.port = &port
	return b
}

func (b *ConfigBuilder) Build() (Config, error) {
	cfg := Config{}

	// the logic for port management goes
	// inside the builder so we can configure
	// our Config struct here to return.
	if b.port == nil {
		cfg.Port = defaultHTTPPort
	} else {
		if *b.port == 0 {
			cfg.Port = randomPort()
		} else if *b.port < 0 {
			return Config{}, errors.New("port should be positive")
		} else {
			cfg.Port = *b.port
		}
	}

	return cfg, nil
}

func NewServer(addr string, config Config) (*http.Server, error) {
	return nil, nil
}

func client() error {
	builder := ConfigBuilder{}
	builder.Port(8080)
	cfg, err := builder.Build()
	if err != nil {
		return err
	}

	server, err := NewServer("localhost", cfg)
	if err != nil {
		return err
	}

	_ = server
	return nil
}

func randomPort() int {
	return 4 // Chosen by fair dice roll, guaranteed to be random.
}
