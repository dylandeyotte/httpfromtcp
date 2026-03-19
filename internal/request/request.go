package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const bufSize = 8

func (r *Request) parse(data []byte) (int, error) {
	if r.State == requestStateInitialized {
		rl, num, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if num == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.State = requestStateDone
		return num, nil
	} else if r.State == requestStateDone {
		return 0, fmt.Errorf("Cannot read data in done state")
	}
	return 0, fmt.Errorf("Unknown state")
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufSize)
	readToIndex := 0
	r := &Request{
		State: requestStateInitialized,
	}
	for r.State != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		num, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.State = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += num
		parsedBytesNum, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[parsedBytesNum:])
		readToIndex -= parsedBytesNum
	}
	return r, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	num := bytes.Index(data, []byte("\r\n"))
	if num == -1 {
		return nil, 0, nil
	}
	newLineSplits := string(data[:num])
	requestLineSplits := strings.Split(newLineSplits, " ")
	if len(requestLineSplits) != 3 {
		return nil, 0, fmt.Errorf("Invalid request-line length: %v", len(requestLineSplits))
	}
	for _, i := range requestLineSplits[0] {
		if i < 'A' || i > 'Z' {
			return nil, 0, fmt.Errorf("Invlaid request-line: %s", requestLineSplits[0])
		}
	}
	httpSplit := strings.Split(requestLineSplits[2], "/")
	if len(httpSplit) != 2 {
		return nil, 0, fmt.Errorf("Invalid HTTP version: %s", httpSplit)
	}
	if httpSplit[0] != "HTTP" {
		return nil, 0, fmt.Errorf("Invalid HTTP version: %s", httpSplit[0])
	}
	if httpSplit[1] != "1.1" {
		return nil, 0, fmt.Errorf("Invalid HTTP version: %s", httpSplit[1])
	}
	return &RequestLine{
		HttpVersion:   httpSplit[1],
		RequestTarget: requestLineSplits[1],
		Method:        requestLineSplits[0],
	}, num + 2, nil
}
