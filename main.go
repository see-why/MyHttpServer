package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(reader io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer reader.Close()

		buffer := make([]byte, 8)
		currentLine := ""

		for {
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Error reading: %v\n", err)
				return
			}

			parts := strings.Split(string(buffer[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				ch <- currentLine
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]
		}

		if currentLine != "" {
			ch <- currentLine
		}
	}()

	return ch
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error creating listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Println("Connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

		fmt.Println("Connection closed")
	}
}
