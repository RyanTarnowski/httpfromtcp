package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}
	request_parts := strings.Split(string(request), "\r\n")

	request_line, err := parseRequestLine(request_parts[0])
	if err != nil {
		return &Request{}, err
	}

	return &Request{
		RequestLine: *request_line,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	request_line_part := strings.Split(line, " ")

	if len(request_line_part) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", line)
	}

	if strings.ToUpper(request_line_part[0]) != request_line_part[0] {
		return nil, fmt.Errorf("invalid method: %s", request_line_part[0])
	}

	versionParts := strings.Split(request_line_part[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", line)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: request_line_part[1],
		Method:        request_line_part[0],
	}, nil
}
