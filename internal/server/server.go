package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/RyanTarnowski/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error creating listener: %v", err)
	}

	server := &Server{
		listener: l,
	}
	server.isClosed.Store(false)

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for s.isClosed.Load() == false {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() == false {
				log.Printf("Failed to accept connection: %v", err)
			}

			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.StatusCodeSuccess)
	if err != nil {
		fmt.Printf("error writing status line: %v\n", err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		fmt.Printf("error writing headers: %v\n", err)
	}
}
