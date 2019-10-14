package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"errors"
	"github.com/wynemo/kiss-proxy/util"
)


func main() {
	l, err := net.Listen("tcp4", "0.0.0.0:8118")
	if err != nil {
		fmt.Println("listen error", err.Error())
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		fmt.Println("socks connection from ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}


func changeFirstLine(conn net.Conn, data []byte) (util.Foo, error) {
	var pos int
	var foo util.Foo
	CRLF := []byte("\r\n")
	tmpBuf := make([]byte, 64*1024)
	i := 0
	for {
		pos = bytes.Index(data, CRLF)
		i++
		if pos == -1 {
			n, err := conn.Read(tmpBuf)
			if err != nil {
				fmt.Println("read error", err.Error())
				return foo, err
			}
			data = append(data, tmpBuf[:n]...)
		} else {
			break
		}
		if i == 50 {
			return foo, errors.New("dont fuck with me")
		}
	}
	tmp := bytes.Split(data[:pos], []byte(" "))
	if len(tmp) != 3 {
		return foo, errors.New("dont fuck with me")
	}

	method, uri, version := string(tmp[0]), string(tmp[1]), string(tmp[2])
	//fmt.Println("method uri", method, uri)
	var host string
	if strings.ToUpper(method) == "CONNECT" {
		tmp := strings.Split(uri, ":")
		if len(tmp) != 2 {
			fmt.Println("dont fuck with me")
			conn.Close()
		}
		host = uri
		a := []byte("HTTP/1.1 200 Connection established")
		b := []byte("Proxy-agent: proxy.py v1")
		c := [][]byte{a, b, CRLF}
		resp := bytes.Join(c, CRLF)
		conn.Write(resp)
	} else {
		myurl, err := url.Parse(uri)
		if err != nil {
			fmt.Println("invalid url")
		}
		host = myurl.Host
		path := myurl.EscapedPath()
		query := myurl.RawQuery
		if len(query) > 0 {
			path += "?" + query
		}
		fragment := myurl.Fragment
		if len(fragment) > 0 {
			path += "#" + fragment
		}
		if strings.Index(host, ":") == -1 {
			host = host + ":80"
		}
		newLine := strings.Join([]string{method, path, version}, " ")
		//fmt.Println("path is", path)
		log.Println("new line is", newLine)
		data = append([]byte(newLine), data[pos:]...)
	}
	foo.Data = data
	foo.Host = host
	foo.Method = method
	return foo, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var data []byte

	foo, err := changeFirstLine(conn, data)

	if err != nil {
		fmt.Println(err)
		return
	}

	host := foo.Host
	//method := foo.Method
	log.Println("host", host)
	remote, err := net.Dial("tcp4", host)
	if err != nil {
		fmt.Println("can't connect to remote")
		fmt.Println(err)
		return
	}

	defer remote.Close()

	if strings.ToUpper(foo.Method) != "CONNECT" {
		remote.Write(foo.Data)
		util.PipeAndChangeLater(remote, conn, changeFirstLine)
	} else {
		util.Pipe(remote, conn)
	}
}
