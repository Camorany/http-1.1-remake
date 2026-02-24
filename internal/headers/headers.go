package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func IsInvalidFieldName(s string) bool {
	pattern := regexp.MustCompile("[^a-zA-Z0-9!#$%&'*+-.^_`|~]")
	return pattern.MatchString(s) || len(s) < 1
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
			done = true
			n = n + 2
			break
		}

		// Getting field-line and (trimmed) field-value
		splitStrings := strings.SplitN(headerString, ":", 2)
		headerFieldLine := strings.ToLower(splitStrings[0])
		headerFieldValue := strings.TrimSpace(splitStrings[1])

		// If header hasn't been split into a field-line and field-value only, throw error
		if len(splitStrings) != 2 {
			err = errors.New("string splitting error")
			return 0, done, err
		}

		// If field-line contains invalid characters or is empty
		if IsInvalidFieldName(headerFieldLine) {
			err = errors.New("field-line formatting error")
			return 0, done, err
		}

		// If field-line contains whitespace
		if strings.Contains(headerFieldLine, " ") {
			err = errors.New("field-line formatting error")
			return 0, done, err
		}

		// Add field-line and field-value to headers map
		_, exists := h[headerFieldLine]

		if exists {
			h[headerFieldLine] = h[headerFieldLine] + ", " + headerFieldValue
		} else {
			h[headerFieldLine] = headerFieldValue
		}

		n = n + len(headerString)

	}

	return n, done, err
}
