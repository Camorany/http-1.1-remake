package server

import (
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

type Handler func(w *response.Writer, req *request.Request) *HandlerError

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

	responseWriter := response.Writer{
		Connection: connection,
	}

	if err != nil {
		handlerError := &HandlerError{
			StatusCode:   response.StatusBadRequest,
			ErrorMessage: err.Error(),
		}

		handlerError.Write(&responseWriter)
		return
	}

	handlerError := s.handler(&responseWriter, request)

	if handlerError != nil {
		handlerError.Write(&responseWriter)
		return
	}

}

func (s *Server) Close() error {
	err := s.listener.Close()

	if err != nil {
		return err
	}

	s.state = closed

	return nil
}

func (handlerError *HandlerError) Write(w *response.Writer) {

	w.WriteStatusLine(handlerError.StatusCode)
	w.WriteHeaders(response.GetDefaultHeaders(len(handlerError.ErrorMessage)))
	w.WriteBody([]byte(handlerError.ErrorMessage))
}

func WriteError(w io.Writer, err HandlerError) {

	_, writeErr := fmt.Fprint(w, err.ErrorMessage)

	if writeErr != nil {
		panic(writeErr)
	}
}
