package request

import (
	"io"
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
	var request Request
	var err error

	request.RequestLine.HttpVersion = "1.1"
	request.RequestLine.RequestTarget = "/"
	request.RequestLine.Method = "GET"

	return &request, err
}
