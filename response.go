package main

import (
	"fmt"
	"io"
)

type Response struct {
	status  string
	headers map[string]string
	body    string
}

func (res *Response) WriteResponse(w io.Writer) {
	io.WriteString(w, res.status+"\r\n")
	for k, v := range res.headers {
		io.WriteString(w, fmt.Sprintf("%s: %s\r\n", k, v))
	}
	io.WriteString(w, "\r\n")
	if res.body != "" {
		io.WriteString(w, string(res.body))
	}
}

func (res *Response) SetHeader(key, value string) {
	if res.headers == nil {
		res.headers = make(map[string]string)
	}
	res.headers[key] = value
}

func NotFoundError(w io.Writer) {
	resp := Response{
		status: "HTTP/1.1 404 Not Found",
		body:   string("<h1>Not Found Error 404</h1>"),
	}
	resp.SetHeader("Content-Type", "text/html")
	resp.WriteResponse(w)
}

func InternalServerError(w io.Writer) {
	resp := Response{
		status: "HTTP/1.1 500 Internal Server Error",
		body:   string("<h1>Internal Server Error 500</h1>"),
	}
	resp.SetHeader("Content-Type", "text/html")
	resp.WriteResponse(w)
}

func GetOk(w io.Writer, data []byte) {
	resp := Response{
		status: "HTTP/1.1 200 OK",
		body:   string(data),
	}
	resp.SetHeader("Content-Type", "text/html")
	resp.WriteResponse(w)
}

func PostOK(w io.Writer) {
	resp := Response{
		status: "HTTP/1.1 200 OK",
		body:   "",
	}
	resp.WriteResponse(w)
}
