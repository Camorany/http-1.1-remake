package response

import (
	"fmt"
	"http_task_module/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOk            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

type responseState int

const (
	Initialized responseState = iota
	WritingStatusLine
	WritingHeaders
	WritingBody
	Done
)

type Writer struct {
	Connection io.Writer
	State      responseState
}

func (s responseState) String() string {
	switch s {
	case Initialized:
		return "Initialized"
	case WritingStatusLine:
		return "WritingStatusLine"
	case WritingHeaders:
		return "WritingHeaders"
	case WritingBody:
		return "WritingBody"
	case Done:
		return "Done"
	default:
		return "unknown"
	}
}

func WriteStateError(actualState responseState, expectedState responseState) error {
	return fmt.Errorf("Incorrect response writer state: expected %s, actual: %s", expectedState.String(), actualState.String())
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	newHeaders := headers.NewHeaders()

	newHeaders["content-length"] = strconv.Itoa(contentLen)
	newHeaders["connection"] = "close"
	newHeaders["content-type"] = "text/plain"

	return newHeaders
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var err error

	if w.State != WritingStatusLine {
		return WriteStateError(w.State, WritingHeaders)
	}

	switch statusCode {
	case StatusOk:
		_, err = w.Connection.Write([]byte("HTTP/1.1 200 OK\r\n"))

	case StatusBadRequest:
		_, err = w.Connection.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))

	case StatusInternalError:
		_, err = w.Connection.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}

	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var err error

	if w.State != WritingHeaders {
		return WriteStateError(w.State, WritingHeaders)
	}

	for fieldLine, fieldValue := range headers {
		_, err = w.Connection.Write([]byte(fmt.Sprintf("%s: %s\r\n", fieldLine, fieldValue)))

		if err != nil {
			return err
		}
	}

	// Add one carriage line after headers are done being written
	_, carriageLineErr := w.Connection.Write([]byte("\r\n"))
	if carriageLineErr != nil {
		return carriageLineErr
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.State != WritingBody {
		return 0, WriteStateError(w.State, WritingHeaders)
	}

	bytesWritten, bodyErr := w.Connection.Write(p)

	if bodyErr != nil {
		return bytesWritten, bodyErr
	}

	return bytesWritten, nil
}
