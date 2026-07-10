package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLine := ""

		for {
			bSlice := make([]byte, 8)

			n, err := f.Read(bSlice)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Printf("Error reading file: %v", err)
				}
				return
			}

			str := string(bSlice[:n])
			parts := strings.Split(str, "\n")

			for i, part := range parts {
				if i > 0 {
					ch <- currentLine
					currentLine = part
				} else {
					currentLine += part
				}
			}
		}

	}()

	return ch
}
