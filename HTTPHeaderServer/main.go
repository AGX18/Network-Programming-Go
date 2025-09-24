package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func main()  {
	// syscall.AF_INET is for IPv4, SOCK_STREAM for TCP
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "socket error: %v\n", err)
		os.Exit(1)
	}

	defer syscall.Close(fd)

	syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: 7777,
		Addr: [4]byte{127, 0, 0, 1},
	})

	syscall.Listen(fd, 128)

	fmt.Println("Listening on port 7777...")

	// Accept connections in a loop

	for {
		nfd, sa, err := syscall.Accept(fd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept error: %v\n", err)
			continue
		}
		fmt.Printf("Accepted connection from %v\n", sa)

		// Now nfd is a blocking file descriptor for the accepted connection.
		// We can use syscall.Read and syscall.Write on it.
		// In a real application, you would likely want to set a timeout
		// and handle partial reads/writes.
		// For simplicity, we'll just read once and echo it back.

		data := make([]byte, 4096)
		n, err := syscall.Read(nfd, data)
		if err != nil || n == 0 {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			syscall.Close(nfd)
			continue
		}
		req := string(data[:n])
		HeadersAndBody := strings.Split(req, "\r\n\r\n")
		Headers := strings.Split(HeadersAndBody[0], "\r\n")
		body := HeadersAndBody[1]
		responseBody := make(map[string]string)
		for hline := range Headers {
			fmt.Printf("Header Line %d: %s\n", hline, Headers[hline])
			parts := strings.SplitN(Headers[hline], ": ", 2)
			if len(parts) == 2 {
				responseBody[parts[0]] = parts[1]
			}
		}
		fmt.Printf("Body: %s\n", body)
		fmt.Printf("Received %d bytes: %s\n", n, string(data[:n]))
		
		fmt.Printf("Handled connection from %v\n", sa)


		jsonBody, _ := json.Marshal(responseBody)
		response := "HTTP/1.1 200 OK\r\n\r\n" 
		syscall.Write(nfd, []byte(response))
		syscall.Write(nfd, jsonBody)
		syscall.Write(nfd, []byte("\r\n"))
		// Close the accepted connection

		syscall.Close(nfd)
	}
}