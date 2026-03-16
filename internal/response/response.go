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

const (
	WritingStatusLine = iota
	WritingHeaders
	WritingBody
	Done
)

type Writer struct {
	Connection io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var err error

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

func GetDefaultHeaders(contentLen int) headers.Headers {
	newHeaders := headers.NewHeaders()

	newHeaders["content-length"] = strconv.Itoa(contentLen)
	newHeaders["connection"] = "close"
	newHeaders["content-type"] = "text/plain"

	return newHeaders
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var err error

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

	bytesWritten, bodyErr := w.Connection.Write(p)

	if bodyErr != nil {
		return bytesWritten, bodyErr
	}

	return bytesWritten, nil
}
