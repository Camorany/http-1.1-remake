package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	n = 0
	done = false
	err = nil

	// If data is missing CRLF, not enough data to parse yet
	endIndex := strings.Index(string(data), "\r\n")
	if endIndex == -1 {
		return n, done, err
	}

	// Get individual headers as strings
	headerStrings := strings.Split(string(data), "\r\n")

	// Parse each headerString into header map
	for _, headerString := range headerStrings {

		//
		if headerString == "" {
			break
		}

		// Getting field-line and (trimmed) field-value
		splitStrings := strings.SplitN(headerString, ":", 2)
		headerFieldLine := splitStrings[0]
		headerFieldValue := strings.TrimSpace(splitStrings[1])

		// If header hasn't been split into a field-line and field-value only, throw error
		if len(splitStrings) != 2 {
			err = errors.New("string splitting error")
			return n, done, err
		}

		// If field-line contains whitespace
		if strings.Contains(headerFieldLine, " ") {
			err = errors.New("field-line formatting error")
			return n, done, err
		}

		// Add field-line and field-value to headers map
		h[headerFieldLine] = headerFieldValue

	}

	n = len(headerStrings[0]) + 2
	return n, done, err
}
