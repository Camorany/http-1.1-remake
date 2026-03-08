package server

import (
	"fmt"
	"http_task_module/internal/response"
	"net"
)

const (
	running = iota
	closed  = iota
)

type Server struct {
	listener net.Listener
	state    int
}

func Serve(port int) (*Server, error) {
	tcpListener, listenerErr := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))

	if listenerErr != nil {
		return nil, listenerErr
	}

	server := &Server{tcpListener, running}

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

	headers := response.GetDefaultHeaders(0)

	response.WriteStatusLine(connection, response.StatusOk)
	response.WriteHeaders(connection, headers)
}

func (s *Server) Close() error {
	err := s.listener.Close()

	if err != nil {
		return err
	}

	s.state = closed

	return nil
}
