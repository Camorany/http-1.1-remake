package request

import (
	"errors"
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

	if requestLine.Method == "" || requestLine.RequestTarget == "" || requestLine.HttpVersion == "" {
		err = errors.New("incorrect number of parts in request line ")
	}

	return requestLine, err
}
