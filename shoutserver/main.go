package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

func main()  {
	// syscall.AF_INET is for IPv4, SOCK_STREAM for TCP
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "socket error: %v\n", err)
		os.Exit(1)
	}
	// At this point, `fd` is just an integer. It's a blocking file descriptor.
	// To use it safely in Go, wrap it in an os.File
	file := os.NewFile(uintptr(fd), "my-socket")

	// Then, convert the os.File into a Go net.Conn
	conn, err := net.FileConn(file)
	if err != nil {
		log.Fatalf("net.FileConn failed: %v", err)
	}
	defer conn.Close()

	log.Printf("Successfully created and wrapped socket: %v", conn.LocalAddr())

}