package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/RyanTarnowski/httpfromtcp/internal/request"
	"github.com/RyanTarnowski/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error creating listener: %v", err)
	}

	server := &Server{
		listener: l,
		handler:  handler,
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

	writer := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		he := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		he.Write(conn)
		return
	}
	var buff bytes.Buffer
	he := s.handler(&buff, req)
	if he != nil {
		he.Write(conn)
	} else {
		err := writer.WriteStatusLine(response.StatusCodeSuccess)
		//err := response.WriteStatusLine(conn, response.StatusCodeSuccess)
		if err != nil {
			fmt.Printf("error writing status line: %v\n", err)
		}
		err = writer.WriteHeaders(response.GetDefaultHeaders(buff.Len()))
		//err = response.WriteHeaders(conn, response.GetDefaultHeaders(buff.Len()))
		if err != nil {
			fmt.Printf("error writing headers: %v\n", err)
		}
		conn.Write(buff.Bytes())
	}
}

func (he HandlerError) Write(w io.Writer) {
	writer := response.NewWriter(w)

	err := writer.WriteStatusLine(he.StatusCode)
	//err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		fmt.Printf("error writing status line: %v\n", err)
	}
	err = writer.WriteHeaders(response.GetDefaultHeaders(len(he.Message)))
	//err = response.WriteHeaders(w, response.GetDefaultHeaders(len(he.Message)))
	if err != nil {
		fmt.Printf("error writing headers: %v\n", err)
	}

	w.Write([]byte(he.Message))
}
