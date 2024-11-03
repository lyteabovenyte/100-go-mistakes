package main

import (
	"errors"
	"net/http"
)

const defaultHTTPPort = 8080

type options struct {
	port *int
}

type Option func(options *options) error

// closure --> an anonymous funciton that references a variable
// outside of its body.
func WithPort(port int) Option {
	// it returns pointer to the actual options type
	// which can be used to reference a variable
	// outside of its body to the options field
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		options.port = &port
		return nil
	}
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options) // it sets up the options struct.
		if err != nil {
			return nil, err
		}
	}

	// At this stage, the options struct is built and contains the config
	// Therefore, we can implement our logic related to port configuration
	var port int
	if options.port == nil {
		port = defaultHTTPPort
	} else {
		if *options.port == 0 {
			port = randomPort()
		} else {
			port = *options.port
		}
	}

	_ = port
	return nil, nil
}

func client() {
	_, _ = NewServer("localhost", WithPort(8080))
}

func randomPort() int {
	return 4 // Chosen by fair dice roll, guaranteed to be random.
}
