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

	readIndex := 0

	switch r.RequestStatus {
	case initalized:
		parsedRequestLine, parseN, parseErr := ParseRequestLine(data[readIndex:])

		if parseErr != nil {
			return 0, parseErr
		}

		if parseN == 0 {
			return 0, nil
		}

		if parseN != 0 {
			r.RequestStatus = done
		}

		r.RequestLine = *parsedRequestLine
		readIndex += parseN

	case done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}

	return readIndex, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, 8)
	readToIndex := 0

	var request Request
	request.RequestStatus = initalized

	for request.RequestStatus != done {

		if len(buf) == readToIndex {
			tmp := make([]byte, 8)
			buf = append(buf, tmp...)
		}

		readN, readErr := reader.Read(buf[readToIndex:])

		if readErr == io.EOF {
			return nil, readErr
		}

		readToIndex += readN

		parseN, parseErr := request.parse(buf[:readToIndex])

		if parseErr != nil {
			return nil, parseErr
		}

		copy(buf, buf[parseN:readToIndex])
		readToIndex -= parseN
	}

	return &request, nil
}

func ParseRequestLine(data []byte) (*RequestLine, int, error) {
	var requestLine RequestLine
	var err error
	var noOfBytes int

	httpRequestStringIndex := strings.Index(string(data), "\r\n")

	if httpRequestStringIndex == -1 {
		return nil, noOfBytes, nil
	}

	requestLineString := string(data)[:httpRequestStringIndex]

	if strings.Contains(string(data), "\r\n") {
		noOfBytes = len(requestLineString)
	} else {
		noOfBytes = 0
	}

	requestLine.Method = strings.Trim(strings.Split(requestLineString, "/")[0], " ")
	requestLine.RequestTarget = strings.Trim(strings.Split(requestLineString, " ")[1], " ")
	requestLine.HttpVersion = strings.Trim(strings.Split(strings.Split(requestLineString, " ")[2], "/")[1], " ") // Fix this to be more readable / less overcomplicated

	if !IsUpperCase(requestLine.Method) || !strings.Contains(requestLine.RequestTarget, "/") || !(requestLine.HttpVersion == "1.1") || IsMissingPart(requestLine) {
		err = errors.New("request-line formatting error")
	}

	return &requestLine, noOfBytes, err
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
