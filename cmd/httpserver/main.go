package main

import (
	"http_task_module/internal/request"
	"http_task_module/internal/response"
	"http_task_module/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42067

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode:   400,
			ErrorMessage: "Your problem is not my problem\r\n",
		}

	case "/myproblem":
		return &server.HandlerError{
			StatusCode:   500,
			ErrorMessage: "Woopsie, my bad\r\n",
		}

	default:
		w.WriteStatusLine(200)
		w.WriteHeaders(response.GetDefaultHeaders(len("Yippie it works!!!!\r\n")))
		w.WriteBody([]byte("Yippie it works!!!!\r\n"))
		return nil
	}
}
