package main

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Request struct {
	headers map[string]string
	body    string
}

func NewHttpRequest(r *textproto.Reader) (*Request, error) {
	h := &Request{headers: make(map[string]string)}

	isFirst := true
	for {
		line, err := r.ReadLine()
		if line == "" {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		// Request Line
		if isFirst {
			isFirst = false
			headerLine := strings.Fields(line)
			h.headers["Method"] = headerLine[0]
			h.headers["Path"] = headerLine[1]
			fmt.Println(headerLine[0], headerLine[1])
			continue
		}

		// Header Fields
		headerFields := strings.SplitN(line, ": ", 2)
		fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
		h.headers[headerFields[0]] = headerFields[1]
	}
	return h, nil
}

func (req *Request) GetRequestBody(r *bufio.Reader, s *textproto.Reader) error {
	lenStr, ok := req.headers["Content-Length"]

	if ok {
		len, err := strconv.Atoi(lenStr)
		if err != nil {
			return errors.WithStack(err)
		}
		buf := make([]byte, len)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Println("BODY:", string(buf))
		// chunked transfer
	} else {
		transferEncoding, ok := req.headers["Transfer-Encoding"]
		if !ok {
			return errors.New("no match operation")
		}
		if transferEncoding == "chunked" {
			for {
				line, err := s.ReadLine()
				if line == "0" {
					break
				}
				if err != nil {
					return errors.WithStack(err)
				}
				fmt.Println(line)
			}
		}
	}

	return nil
}
