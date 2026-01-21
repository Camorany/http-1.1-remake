package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

const (
	initalized = iota
	done       = iota
)

type Request struct {
	RequestLine   RequestLine
	RequestStatus int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

	if r.RequestStatus == initalized {
		parsedRequestLine, noOfBytes, parseErr := ParseRequestLine(data)

		if parseErr != nil {
			return 0, parseErr
		}

		r.RequestLine = parsedRequestLine

		return noOfBytes, parseErr

	} else if r.RequestStatus == done {
		return 0, errors.New("error: trying to read data in a done state")
	} else {
		return 0, errors.New("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var err error
	data, dataErr := io.ReadAll(reader)

	if dataErr != nil {
		err = dataErr
	}

	parsedRequestLine, noOfBytes, parseErr := ParseRequestLine(data)

	if parseErr != nil {
		err = parseErr
	}

	var request Request
	request.RequestLine = parsedRequestLine

	return &request, err
}

func ParseRequestLine(data []byte) (RequestLine, int, error) {
	var requestLine RequestLine
	var err error
	var noOfBytes int
	httpRequestString := strings.Split(string(data), "\r\n")

	if strings.Contains(string(data), "\r\n") {
		noOfBytes = len(data)
	} else {
		noOfBytes = 0
	}

	requestLine.Method = strings.Trim(strings.Split(httpRequestString[0], "/")[0], " ")
	requestLine.RequestTarget = strings.Trim(strings.Split(httpRequestString[0], " ")[1], " ")
	requestLine.HttpVersion = strings.Trim(strings.Split(httpRequestString[0], "/")[2], " ")

	if !IsUpperCase(requestLine.Method) || !strings.Contains(requestLine.RequestTarget, "/") || !(requestLine.HttpVersion == "1.1") || IsMissingPart(requestLine) {
		err = errors.New("request-line formatting error")
	}

	return requestLine, noOfBytes, err
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
