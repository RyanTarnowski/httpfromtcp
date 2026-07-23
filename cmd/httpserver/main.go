package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/RyanTarnowski/httpfromtcp/internal/request"
	"github.com/RyanTarnowski/httpfromtcp/internal/response"
	"github.com/RyanTarnowski/httpfromtcp/internal/server"
)

const port = 42069

func main() {
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

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		//target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

		resp, err := http.Get("https://httpbin.org/html")
		if err != nil {
			handler500(w, req)
			return
		}
		defer resp.Body.Close()

		w.WriteStatusLine(response.StatusCodeSuccess)
		h := response.GetDefaultHeaders(0)
		h.Replace("Content-Length", "Transfer-Encoding", "chunked")
		h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
		w.WriteHeaders(h)

		buf := make([]byte, 1024)
		fullRespBody := []byte{}

		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				_, err := w.WriteChunkedBody(buf[:n])
				if err != nil {
					fmt.Println("Error writing chunked body:", err)
					break
				}

				fullRespBody = append(fullRespBody, buf[:n]...)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading response body:", err)
				break
			}
		}

		w.SetChunkedBodeEndState()
		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			fmt.Println("Error writing chunked body done:", err)
		}
		h.Set("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256(fullRespBody)))
		h.Set("X-Content-Length", fmt.Sprintf("%d", len(fullRespBody)))
		err = w.WriteTrailers(h)
		if err != nil {
			fmt.Println("Error writing trailers:", err)
		}
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handler400(w, req)
		return
	case "/myproblem":
		handler500(w, req)
		return
	default:
		handler200(w, req)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}
