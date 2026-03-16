package server

import (
	"bytes"
	"fmt"
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

func Serve(port int, handlerFunc Handler) (*Server, error) {
	tcpListener, listenerErr := net.Listen("tcp", fmt.Sprintf(":%d", port))

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
			if s.state == closed {
				return
			}
			fmt.Printf("Error accepting connection: %v", connErr)
			continue
		}

		go s.handle(connection)

	}
}

func (s *Server) handle(connection net.Conn) {
	defer connection.Close()

	request, err := request.RequestFromReader(connection)

	if err != nil {
		handlerError := &HandlerError{
			StatusCode:   response.StatusBadRequest,
			ErrorMessage: err.Error(),
		}

		handlerError.Write(connection)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(buffer, request)

	if handlerError != nil {
		handlerError.Write(connection)
		return
	}

	headers := response.GetDefaultHeaders(buffer.Len())
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

func (handlerError *HandlerError) Write(connection net.Conn) {

	response.WriteStatusLine(connection, handlerError.StatusCode)
	headers := response.GetDefaultHeaders(len(handlerError.ErrorMessage))
	response.WriteHeaders(connection, headers)
	WriteError(connection, *handlerError)
}

func WriteError(w io.Writer, err HandlerError) {

	_, writeErr := fmt.Fprint(w, err.ErrorMessage)

	if writeErr != nil {
		panic(writeErr)
	}
}
