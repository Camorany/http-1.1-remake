package response

import (
	"fmt"
	"http_task_module/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOk            = 200
	StatusNotFound      = 400
	StatusInternalError = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error

	switch statusCode {
	case StatusOk:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))

	case StatusNotFound:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))

	case StatusInternalError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
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

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var err error
	for fieldLine, fieldValue := range headers {
		_, err = w.Write([]byte(fmt.Sprintf("%s: %s\r\n", fieldLine, fieldValue)))

		if err != nil {
			return err
		}
	}

	_, carriageLineErr := w.Write([]byte("\r\n"))

	if carriageLineErr != nil {
		return carriageLineErr
	}

	return nil
}
