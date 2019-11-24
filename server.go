package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Run() error {
	fmt.Println("start tcp listen...")

	// Listen ポートの生成をする
	listen, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		return errors.WithStack(err)
	}
	defer listen.Close()

	// コネクションを受け付ける
	conn, err := listen.Accept()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	fmt.Println(">>> start")

	reader := bufio.NewReader(conn)
	scanner := textproto.NewReader(reader)

	// 一行ずつ処理する
	// リクエストヘッダー
	req, err := NewHttpRequest(scanner)
	if err != nil {
		return errors.WithStack(err)
	}

	// リクエストボディ
	switch req.headers["Method"] {
	case "GET":
		path, ok := req.headers["Path"]
		if !ok {
			return errors.New("no path found")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		p := filepath.Join(cwd, filepath.Clean(path))
		if err != nil {
			return errors.WithStack(err)
		}

		// file not found
		if !fileExists(p) {
			NotFoundError(conn)
		} else {
			data, err := ioutil.ReadFile(p)
			if err != nil {
				return errors.WithStack(err)
			}
			GetOk(conn, data)
		}
	case "POST", "PUT":
		if err := req.GetRequestBody(reader, scanner); err != nil {
			return errors.WithStack(err)
		}
		PostOK(conn)
		return nil
	default:
		return errors.New("no match method")
	}
	// completed
	fmt.Println("<<< end")

	return nil
}
