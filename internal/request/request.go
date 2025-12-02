package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
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
	var err error
	data, dataErr := io.ReadAll(reader)

	if dataErr != nil {
		err = dataErr
	}

	parsedRequestLine, parseErr := ParseRequestLine(string(data))

	if parseErr != nil {
		err = parseErr
	}

	var request Request
	request.RequestLine = parsedRequestLine

	return &request, err
}

func ParseRequestLine(data string) (RequestLine, error) {
	var requestLine RequestLine
	var err error
	httpRequestString := strings.Split(data, "\r\n")

	requestLine.Method = strings.Trim(strings.Split(httpRequestString[0], "/")[0], " ")
	requestLine.RequestTarget = strings.Trim(strings.Split(httpRequestString[0], " ")[1], " ")
	requestLine.HttpVersion = strings.Trim(strings.Split(httpRequestString[0], "/")[2], " ")

	if !IsUpperCase(requestLine.Method) || !strings.Contains(requestLine.RequestTarget, "/") || !(requestLine.HttpVersion == "1.1") || IsMissingPart(requestLine) {
		err = errors.New("request-line formatting error")
	}

	return requestLine, err
}

func IsUpperCase(data string) bool {
	isUpper := true
	for _, c := range data {
		if !unicode.IsUpper(c) {
			isUpper = false
		}
	}
	return isUpper
}

func IsMissingPart(data RequestLine) bool {
	if data.HttpVersion == "" || data.Method == "" || data.RequestTarget == "" {
		return true
	}
	return false
}
