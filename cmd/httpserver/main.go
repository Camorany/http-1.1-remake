package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"http_task_module/internal/request"
	"http_task_module/internal/response"
	"http_task_module/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	switch {
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/"):

		// Getting default headers, replacing content-length with transfer-encoding
		headers := response.GetDefaultHeaders(0)
		headers.RemoveHeader("content-length")
		headers.AddHeader("transfer-encoding", "chunked")
		headers.AddHeader("trailer", "X-Content-SHA256, X-Content-Length")

		// Writing status line and headers for response
		w.State = response.WritingStatusLine
		w.WriteStatusLine(200)
		w.State = response.WritingHeaders
		w.WriteHeaders(headers)

		resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		buf := make([]byte, 256)
		var totalContentBuf []byte

		for {
			bytesRead, readErr := resp.Body.Read(buf)
			if bytesRead > 0 {
				// Write body using http response stream
				w.State = response.WritingBody
				w.WriteChunkedBody(buf[:bytesRead])
				totalContentBuf = append(totalContentBuf, buf[:bytesRead]...)
			}

			fmt.Printf("%d\r\n", bytesRead)

			if readErr != nil {
				if readErr == io.EOF {
					break
				}
				panic(readErr)
			}
		}

		chunkedContentHash := sha256.Sum256(totalContentBuf)
		chunkedContentHashString := hex.EncodeToString(chunkedContentHash[:])

		trailers := response.BuildTrailers(chunkedContentHashString, strconv.Itoa(len(totalContentBuf)))

		w.State = response.WritingTrailers
		w.WriteTrailers(trailers)

		return nil
	case req.RequestLine.RequestTarget == "/video":

		mp4Data, mp4Err := os.ReadFile("assets/vim.mp4")

		if mp4Err != nil {
			return &server.HandlerError{
				StatusCode:   400,
				ErrorMessage: []byte(mp4Err.Error()),
			}
		}

		headers := response.GetDefaultHeaders(len(mp4Data))
		headers.OverrideHeader("content-type", "video/mp4")

		// Writing status line and headers for response
		w.State = response.WritingStatusLine
		w.WriteStatusLine(200)

		w.State = response.WritingHeaders
		w.WriteHeaders(headers)

		w.State = response.WritingBody
		w.WriteBody(mp4Data)

		return nil

	case req.RequestLine.RequestTarget == "/yourproblem":
		return &server.HandlerError{
			StatusCode:   400,
			ErrorMessage: response400(),
		}

	case req.RequestLine.RequestTarget == "/myproblem":
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
