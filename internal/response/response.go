package response

import (
	"fmt"
	"io"

	"github.com/RyanTarnowski/httpfromtcp/internal/headers"
)

const crlf = "\r\n"

type Writer struct {
	writer io.Writer
	state  writerState
}

type writerState int

const (
	WritingStatusLine writerState = iota
	WritingHeaders
	WritingBody
)

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  WritingStatusLine,
	}
}

// func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
// 	statusLine := ""
// 	switch statusCode {
// 	case StatusCodeSuccess:
// 		statusLine = "HTTP/1.1 200 OK"
// 	case StatusCodeBadRequest:
// 		statusLine = "HTTP/1.1 400 Bad Request"
// 	case StatusCodeInternalServerError:
// 		statusLine = "HTTP/1.1 500 Internal Server Error"
// 	}
//
// 	_, err := w.Write([]byte(statusLine + crlf))
// 	return err
// }

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/html")
	return h
}

// func WriteHeaders(w io.Writer, headers headers.Headers) error {
// 	for key, value := range headers {
// 		_, err := w.Write([]byte(fmt.Sprintf("%s: %s%s", key, value, crlf)))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	_, err := w.Write([]byte(crlf))
//
// 	return err
// }

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WritingStatusLine {
		return fmt.Errorf("Attempted to write status line out of order")
	}

	statusLine := ""
	switch statusCode {
	case StatusCodeSuccess:
		statusLine = "HTTP/1.1 200 OK"
	case StatusCodeBadRequest:
		statusLine = "HTTP/1.1 400 Bad Request"
	case StatusCodeInternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	}

	_, err := w.writer.Write([]byte(statusLine + crlf))
	w.state = WritingHeaders
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WritingHeaders {
		return fmt.Errorf("Attempted to write headers out of order")
	}

	for key, value := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s%s", key, value, crlf)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte(crlf))

	w.state = WritingBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WritingBody {
		return 0, fmt.Errorf("Attempted to write body out of order")
	}

	return w.writer.Write(p)
}
