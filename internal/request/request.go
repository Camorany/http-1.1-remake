package request

import (
	"bytes"
	"errors"
	"fmt"
	"http_task_module/internal/headers"
	"io"
	"strings"
	"unicode"
)

type state int

// States monitored for parsing HTTP segments
const (
	initialized    state = iota
	parsingHeaders state = iota
	done           state = iota
)

// Returns string of state name
func (s state) String() string {
	switch s {
	case initialized:
		return "parsingRequestLine"
	case parsingHeaders:
		return "parsingHeaders"
	case done:
		return "done"
	default:
		return "unknown"
	}
}

type Request struct {
	RequestLine   RequestLine
	Headers       headers.Headers
	RequestStatus state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

	bytesParsed := 0

	switch r.RequestStatus {
	case initialized:
		parsedRequestLine, requestLineBytesParsed, parseErr := ParseRequestLine(data)

		if parseErr != nil {
			return 0, parseErr
		}

		if requestLineBytesParsed == 0 {
			return 0, nil
		}

		r.RequestStatus = parsingHeaders
		r.RequestLine = *parsedRequestLine
		bytesParsed = requestLineBytesParsed

	case parsingHeaders:
		headersBytesParsed, parseHeadersState, headersErr := r.Headers.Parse(data)

		if headersErr != nil {
			return 0, headersErr
		}

		if headersBytesParsed == 0 {
			return 0, nil
		}

		if parseHeadersState == true {
			r.RequestStatus = done
			bytesParsed = headersBytesParsed
		}

	case done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}

	return bytesParsed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, 8)
	readToIndex := 0

	var request Request
	request.RequestStatus = initialized
	request.Headers = headers.NewHeaders()

	for request.RequestStatus != done {

		if len(buf) == readToIndex {
			tmp := make([]byte, 8)
			buf = append(buf, tmp...)
		}

		bytesRead, readErr := reader.Read(buf[readToIndex:])

		if bytesRead == 0 && readErr == io.EOF {
			return nil, fmt.Errorf("EOF error occurred during '%s' state", request.RequestStatus.String())
		}

		readToIndex += bytesRead

		bytesParsed, parseErr := request.parse(buf[:readToIndex])

		if parseErr != nil {
			return nil, parseErr
		}

		// Removes parsed bytes from buffer (leaves unparsed bytes & bytes with no data read into them)
		copy(buf, buf[bytesParsed:readToIndex])
		// Shifts readToIndex to the new index of the final byte in buffer containing data
		readToIndex -= bytesParsed
	}

	return &request, nil
}

func ParseRequestLine(data []byte) (*RequestLine, int, error) {
	var requestLine RequestLine
	var err error
	var bytesParsed int

	httpRequestEndIndex := bytes.Index(data, []byte("\r\n"))

	if httpRequestEndIndex == -1 {
		return nil, 0, nil
	}
	// bytes parsed = number of bytes consumed plus 2 (include /r/n so it's not left in buffer on next iteration)
	bytesParsed = httpRequestEndIndex + 2

	requestLineString := string(data)[:httpRequestEndIndex]
	parts := strings.SplitN(requestLineString, " ", 3)

	if len(parts) != 3 {
		return nil, 0, errors.New("malformed request line")
	}

	requestLine.Method = parts[0]
	requestLine.RequestTarget = parts[1]
	requestLine.HttpVersion = strings.TrimPrefix(parts[2], "HTTP/")

	if !IsUpperCase(requestLine.Method) ||
		!strings.Contains(requestLine.RequestTarget, "/") ||
		!(requestLine.HttpVersion == "1.1") ||
		IsMissingPart(requestLine) ||
		!strings.HasPrefix(parts[2], "HTTP/") {

		return nil, 0, errors.New("request-line formatting error")

	}

	return &requestLine, bytesParsed, err
}

func IsUpperCase(data string) bool {
	for _, c := range data {
		if !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}

func IsMissingPart(data RequestLine) bool {
	if data.HttpVersion == "" || data.Method == "" || data.RequestTarget == "" {
		return true
	}
	return false
}
