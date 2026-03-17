package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLineChannels(f io.ReadCloser) <-chan string {
	strChan := make(chan string)
	b := make([]byte, 8)
	line := ""

	go func() {
		defer f.Close()
		defer close(strChan)
		for {
			n, err := f.Read(b)
			if err != nil {
				if line != "" {
					strChan <- line
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println(err)
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				strChan <- fmt.Sprintf("%s%s", line, parts[i])
				line = ""
			}
			line += parts[len(parts)-1]
		}
	}()

	return strChan
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection successful")

		ch := getLineChannels(conn)
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}

}
