package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) GetHeader(key string) string {
	return h[strings.ToLower(key)]
}

func (headers *Headers) OverrideHeader(overrideHeader string, newHeaderValue string) {
	(*headers)[overrideHeader] = newHeaderValue
}

func (headers *Headers) AddHeader(headerKey string, headerValue string) {
	(*headers)[headerKey] = headerValue
}

func (headers *Headers) RemoveHeader(headerToRemove string) {
	delete((*headers), headerToRemove)
}

func NewHeaders() Headers {
	return make(Headers)
}

func IsInvalidFieldName(s string) bool {
	pattern := regexp.MustCompile("[^a-zA-Z0-9!#$%&'*+-.^_`|~]")
	return pattern.MatchString(s) || len(s) < 1
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	bytesParsed := 0
	done = false
	err = nil

	// If data is missing CRLF, not enough data to parse yet
	headersEndIndex := bytes.Index(data, []byte("\r\n\r\n"))

	if headersEndIndex == -1 {
		return n, done, err
	}

	bytesParsed = headersEndIndex + 4

	// Get individual headers as strings
	headerStringsBlob := strings.Split(string(data[:headersEndIndex]), "\r\n\r\n")[0]
	var headerStrings []string

	if strings.Contains(headerStringsBlob, "\r\n") {
		headerStrings = strings.Split(headerStringsBlob, "\r\n")
	} else {
		headerStrings = []string{headerStringsBlob}
	}

	// Parse each headerString into header map
	for _, headerString := range headerStrings {

		// Getting field-line and (trimmed) field-value
		splitStrings := strings.SplitN(headerString, ":", 2)

		// If header hasn't been split into a field-line and field-value only, throw error
		if len(splitStrings) != 2 {
			err = errors.New("header field-line/field-value splitting error")
			return 0, done, err
		}

		// Get header field-line and field-value from header string
		headerFieldLine := strings.ToLower(splitStrings[0])
		headerFieldValue := strings.TrimSpace(splitStrings[1])

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

	}

	done = true
	return bytesParsed, done, err
}
