package server

import (
	"bytes"
	"fmt"
	"http_task_module/internal/headers"
	"http_task_module/internal/request"
	"http_task_module/internal/response"
	"io"
	"net"
)

const (
	running = iota
	closed  = iota
)

type Server struct {
	listener net.Listener
	state    int
	handler  Handler
}

type HandlerError struct {
	StatusCode   response.StatusCode
	ErrorMessage string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(handlerFunc Handler, port int) (*Server, error) {
	tcpListener, listenerErr := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))

	if listenerErr != nil {
		return nil, listenerErr
	}

	server := &Server{tcpListener, running, handlerFunc}

	go server.listen()

	return server, nil
}

func (s *Server) listen() {
	for {
		connection, connErr := s.listener.Accept()

		if connErr != nil {
			fmt.Print("Error accepting connection: ", connErr)
			continue
		}

		go s.handle(connection)

	}
}

func (s *Server) handle(connection net.Conn) {
	defer connection.Close()

	request, err := request.RequestFromReader(connection)

	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	var headers headers.Headers

	handlerError := s.handler(&buffer, request)

	if handlerError != nil {
		response.WriteStatusLine(connection, handlerError.StatusCode)
		headers = response.GetDefaultHeaders(len(handlerError.ErrorMessage))
		response.WriteHeaders(connection, headers)
		WriteError(connection, *handlerError)
	}

	headers = response.GetDefaultHeaders(buffer.Len())
	response.WriteStatusLine(connection, response.StatusOk)
	response.WriteHeaders(connection, headers)
	connection.Write(buffer.Bytes())
}

func (s *Server) Close() error {
	err := s.listener.Close()

	if err != nil {
		return err
	}

	s.state = closed

	return nil
}

func WriteError(w io.Writer, err HandlerError) {

	_, writeErr := fmt.Fprint(w, err.ErrorMessage)

	if writeErr != nil {
		panic(writeErr)
	}
}
