// Package grpcserver implements gRPC server.
package server

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	port            string
	netListener     net.Listener
	server          *grpc.Server
	notify          chan error
	shutdownTimeout time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
}

// New -.
func New(server *grpc.Server, lis net.Listener, port string, opts ...Option) *Server {
	s := &Server{
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		port:            port,
	}

	s.server = server
	s.netListener = lis

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		fmt.Println("gRPC server running on port", s.port)
		s.notify <- s.server.Serve(s.netListener)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	fmt.Println("grpc server is stopped")
	s.server.Stop()
	return nil
}

// GetServer -.
func (s *Server) GetServer() *grpc.Server {
	return s.server
}
