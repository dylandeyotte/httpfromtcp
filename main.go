package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				if line != "" {
					strChan <- line
				}
				fmt.Print(err)
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
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ch := getLineChannels(file)

	for i := range ch {
		fmt.Println("read:", i)
	}
}
