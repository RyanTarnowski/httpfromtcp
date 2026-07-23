package response

import (
	"fmt"
	"io"
	"strings"

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
	WritingBodyEnd
	WritingTrailers
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

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/html")
	return h
}

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
		_, err := fmt.Fprintf(w.writer, "%s: %s%s", key, value, crlf)
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

	w.state = WritingBodyEnd
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != WritingBody {
		return 0, fmt.Errorf("Attempted to write body out of order")
	}

	return fmt.Fprintf(w.writer, "%X%s%s%s", len(p), crlf, p, crlf)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != WritingBodyEnd {
		return 0, fmt.Errorf("Attempted to write end of body out of order")
	}

	w.state = WritingTrailers
	return fmt.Fprintf(w.writer, "0%s", crlf)
}

func (w *Writer) SetChunkedBodeEndState() {
	w.state = WritingBodyEnd
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.state != WritingTrailers {
		return fmt.Errorf("Attempted to write trailers out of order")
	}

	trailers, ok := h.GetValueByKey("trailer")
	if ok {
		for key := range strings.SplitSeq(trailers, ",") {
			key = strings.TrimSpace(key)
			value, ok := h.GetValueByKey(key)
			if ok {
				_, err := fmt.Fprintf(w.writer, "%s: %s%s", key, value, crlf)
				if err != nil {
					return err
				}
			}
		}
	}
	_, err := w.writer.Write([]byte(crlf))
	if err != nil {
		return err
	}

	return nil
}
