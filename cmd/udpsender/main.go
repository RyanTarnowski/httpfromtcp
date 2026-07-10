package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = ":42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("Error resolving udp addr: %s\n", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error dialing udp addr: %s\n", err.Error())
		return
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", port)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")

		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading string: %s\n", err.Error())
			return
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Fatalf("Error writing string: %s\n", err.Error())
			return
		}

		fmt.Printf("Message sent: %s", msg)
	}
}
