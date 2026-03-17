package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	rl, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	newLineSplits := strings.Split(string(data), "\r\n")
	requestLineSplits := strings.Split(newLineSplits[0], " ")
	if len(requestLineSplits) != 3 {
		return nil, fmt.Errorf("Invalid request-line length: %v", len(requestLineSplits))
	}
	for _, i := range requestLineSplits[0] {
		if i < 'A' || i > 'Z' {
			return nil, fmt.Errorf("Invlaid request-line: %s", requestLineSplits[0])
		}
	}
	httpSplit := strings.Split(requestLineSplits[2], "/")
	if httpSplit[0] != "HTTP" {
		return nil, fmt.Errorf("Invalid HTTP version: %s", httpSplit[0])
	}
	if httpSplit[1] != "1.1" {
		return nil, fmt.Errorf("Invalid HTTP version: %s", httpSplit[1])
	}
	return &RequestLine{
		HttpVersion:   httpSplit[1],
		RequestTarget: requestLineSplits[1],
		Method:        requestLineSplits[0],
	}, nil
}
