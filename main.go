package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b := make([]byte, 8)
	line := ""

	for {
		n, err := file.Read(b)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			if line != "" {
				fmt.Printf("read: %s\n", line)
				line = ""
			}
			log.Fatal(err)
		}
		str := string(b[:n])
		parts := strings.Split(str, "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", line, parts[i])
			line = ""
		}
		line += parts[len(parts)-1]
	}

}
