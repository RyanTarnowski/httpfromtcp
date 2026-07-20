package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RyanTarnowski/httpfromtcp/internal/request"
	"github.com/RyanTarnowski/httpfromtcp/internal/response"
	"github.com/RyanTarnowski/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w io.Writer, req *request.Request) *server.HandlerError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.StatusCodeBadRequest,
				Message: `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`,
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.StatusCodeInternalServerError,
				Message: `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`,
			}
		default:
			w.Write([]byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`))
			return nil
		}
	}

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
