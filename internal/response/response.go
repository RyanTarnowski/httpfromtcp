package response

import (
	"fmt"
	"io"

	"github.com/RyanTarnowski/httpfromtcp/internal/headers"
)

const crlf = "\r\n"

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case StatusCodeSuccess:
		statusLine = "HTTP/1.1 200 OK"
	case StatusCodeBadRequest:
		statusLine = "HTTP/1.1 400 Bad Request"
	case StatusCodeInternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	}

	_, err := w.Write([]byte(statusLine + crlf))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s%s", key, value, crlf)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(crlf))

	return err
}
