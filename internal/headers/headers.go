package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	//End of the headers
	if idx == 0 {
		return 2, true, nil
	}

	headerParts := strings.SplitN(string(data[:idx]), ":", 2)

	//Check for valid key value pair
	if len(headerParts) != 2 {
		return 0, false, fmt.Errorf("invalid key value pair")
	}

	//Check for spaces and tabs
	if strings.ContainsAny(headerParts[0], " \t") {
		return 0, false, fmt.Errorf("invalid key name")
	}

	h[headerParts[0]] = strings.TrimSpace(headerParts[1])

	return idx + 2, false, nil
}
