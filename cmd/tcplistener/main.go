package main

import (
	"fmt"
	"github.com/RyanTarnowski/httpfromtcp/internal/request"
	"log"
	"net"
)

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error starting listener: %s\n", err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s\n", err.Error())
			return
		}

		fmt.Println("Connection accepted from", conn.RemoteAddr())
		fmt.Println("*************************************************")

		request, err := request.RequestFromReader(conn)

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

		fmt.Println("\nConnection to", conn.RemoteAddr(), "closed")
		fmt.Println("*************************************************")
	}
}
