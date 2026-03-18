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

func response400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func response500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func response200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

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
			ErrorMessage: response400(),
		}

	case "/myproblem":
		return &server.HandlerError{
			StatusCode:   500,
			ErrorMessage: response500(),
		}

	default:
		bodyContent := response200()

		headers := response.GetDefaultHeaders(len(bodyContent))
		headers.OverrideHeader("content-type", "text/html")

		w.State = response.WritingStatusLine
		w.WriteStatusLine(200)

		w.State = response.WritingHeaders
		w.WriteHeaders(headers)

		w.State = response.WritingBody
		w.WriteBody([]byte(bodyContent))

		return nil
	}
}
