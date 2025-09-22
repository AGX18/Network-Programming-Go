package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func main()  {
	// syscall.AF_INET is for IPv4, SOCK_STREAM for TCP
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "socket error: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	// Create the target address structure (sockaddr)
	// We want to connect to 127.0.0.1 on port 8080
	addr := &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}

	err = syscall.Bind(fd, addr)
	if err != nil {
		log.Fatalf("Error binding to socket: %v", err)
	}
	fmt.Println("Connected to server! Waiting for a message...")

	// Now fd is connected to the server. We can use syscall.Read and syscall.Write.
	data := make([]byte, 1024)
	for {
		n, sender, err := syscall.Recvfrom(fd, data, 0)
		if err != nil {
			log.Printf("Error receiving from socket: %v", err)
			continue
		}
		log.Printf("Received %d bytes from %v", n, sender)
	
		fmt.Printf("Server sent: %s\n", string(data[:n]))
		for ch := range data[:n] {
			if data[ch] >= 'a' && data[ch] <= 'z' {
				data[ch] = data[ch] - 32
			}
		}

		err = syscall.Sendto(fd, data[:n], 0, sender)
		if err != nil {
			log.Printf("Error writing to socket: %v", err)
		}
	}

}
// At this point, `fd` is just an integer. It's a blocking file descriptor.
// To use it safely in Go, wrap it in an os.File
// file := os.NewFile(uintptr(fd), "my-socket")

// // Then, convert the os.File into a Go net.Conn
// conn, err := net.FileConn(file)
// if err != nil {
// 	log.Fatalf("net.FileConn failed: %v", err)
// }
// defer conn.Close()