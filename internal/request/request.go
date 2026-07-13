package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const buffSize = 8
const crlf = "\r\n"

type State int

const (
	Initialized State = iota
	Done
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	bSlice := make([]byte, buffSize)
	readToIndex := 0
	request := Request{
		State: Initialized,
	}

	for request.State != Done {
		if readToIndex >= len(bSlice) {
			newSlice := make([]byte, len(bSlice)*2)
			copy(newSlice, bSlice)
			bSlice = newSlice
		}

		n, err := reader.Read(bSlice[readToIndex:])
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("Error reading: %v", err)
			}

			request.State = Done
			break
		}

		readToIndex += n

		np, err := request.parse(bSlice[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(bSlice, bSlice[np:readToIndex])
		readToIndex -= np
	}

	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	byteLength := len(data)
	requestLineText := string(data[:idx])
	request_line_part := strings.Split(requestLineText, " ")

	if len(request_line_part) != 3 {
		return nil, 0, fmt.Errorf("poorly formatted request-line: %s", requestLineText)
	}

	if strings.ToUpper(request_line_part[0]) != request_line_part[0] {
		return nil, 0, fmt.Errorf("invalid method: %s", request_line_part[0])
	}

	versionParts := strings.Split(request_line_part[2], "/")
	if len(versionParts) != 2 {
		return nil, 0, fmt.Errorf("malformed start-line: %s", requestLineText)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, 0, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, 0, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: request_line_part[1],
		Method:        request_line_part[0],
	}, byteLength, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		request_line, bytes_consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("error parsing request line: %s", err.Error())
		}
		if bytes_consumed == 0 {
			return 0, nil
		}
		r.RequestLine = *request_line
		r.State = Done
		return bytes_consumed, nil

	case Done:
		return 0, fmt.Errorf("trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unrecognized status: %v", r.State)
	}
}
