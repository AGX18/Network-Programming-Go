package main

import (
	"fmt"
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

		for {
			data := make([]byte, 1024)
			n, err := syscall.Read(nfd, data)
			if err != nil || n == 0 {
				fmt.Fprintf(os.Stderr, "read error: %v\n", err)
				syscall.Close(nfd)
				break
			}
			fmt.Printf("Received %d bytes: %s\n", n, string(data[:n]))
			syscall.Write(nfd, data[:n])
			
			fmt.Printf("Handled connection from %v\n", sa)

		}
		go func() {
			defer syscall.Close(nfd)
		}()
	}
}