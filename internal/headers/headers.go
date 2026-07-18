package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

const crlf = "\r\n"
const allowedSpecialChars = "!#$%&'*+-.^_`|~"

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
		return 0, false, fmt.Errorf("invalid key value pair %v", headerParts)
	}

	//Check for spaces and tabs
	if strings.ContainsAny(headerParts[0], " \t") {
		return 0, false, fmt.Errorf("invalid key name")
	}

	//Check for invalid chars
	isInvalid := func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune(allowedSpecialChars, r)
	}

	if strings.ContainsFunc(headerParts[0], isInvalid) {
		return 0, false, fmt.Errorf("invalid characters in key")
	}

	//value, ok := h[strings.ToLower(headerParts[0])]

	//if ok {
	//h[strings.ToLower(headerParts[0])] = value + ", " + strings.TrimSpace(headerParts[1])
	//} else {
	//h[strings.ToLower(headerParts[0])] = strings.TrimSpace(headerParts[1])
	//}

	h.Set(headerParts[0], headerParts[1])
	return idx + 2, false, nil
}

func (h Headers) GetValueByKey(key string) (string, bool) {
	v, ok := h[strings.ToLower(key)]
	return v, ok
}

func (h Headers) Set(key, value string) {
	value = strings.TrimSpace(value)
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}
