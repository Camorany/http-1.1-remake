package main

import (
	"fmt"
	"http_task_module/internal/request"
	"http_task_module/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(handler, port)
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

func handler(w io.Writer, req *request.Request) *server.HandlerError {
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
		fmt.Fprint(w, "All good, frfr\r\n")
		return nil
	}
}
