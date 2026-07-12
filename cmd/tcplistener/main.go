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
		//lines := getLinesChannelPG(conn)

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
				if currentLine != "" {
					ch <- currentLine
				}
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

//func getLinesChannelPG(f io.ReadCloser) <-chan string {
//	lines := make(chan string)
//	go func() {
//		defer f.Close()
//		defer close(lines)
//		currentLineContents := ""
//		for {
//			b := make([]byte, 8, 8)
//			n, err := f.Read(b)
//			if err != nil {
//				if currentLineContents != "" {
//					lines <- currentLineContents
//				}
//				if errors.Is(err, io.EOF) {
//					break
//				}
//				fmt.Printf("error: %s\n", err.Error())
//				return
//			}
//			str := string(b[:n])
//			parts := strings.Split(str, "\n")
//			for i := 0; i < len(parts)-1; i++ {
//				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
//				currentLineContents = ""
//			}
//			currentLineContents += parts[len(parts)-1]
//		}
//	}()
//	return lines
//}
