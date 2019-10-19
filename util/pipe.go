package util

import (
	"net"
)

func chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 64*1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}

type Foo struct {
	Data   []byte
	Method string
	Host   string
}

type change func(conn net.Conn, data []byte) (Foo, error)

//Pipe pipe too connections
func Pipe(conn1 net.Conn, conn2 net.Conn) []byte {
	chan1 := chanFromConn(conn1)
	chan2 := chanFromConn(conn2)
	closed1 := false
	closed2 := false

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				conn2.Close()
				closed2 = true
			} else {
				conn2.Write(b1)
			}
		case b2 := <-chan2:
			if b2 == nil {
				conn1.Close()
				closed1 = true
			} else {
				conn1.Write(b2)
			}
		}
		if closed1 && closed2 {
			return nil
		}
	}
}

//Pipe pipe too connections
func PipeAndChangeLater(conn1 net.Conn, conn2 net.Conn, fn change) []byte {
	chan1 := chanFromConn(conn1)
	chan2 := chanFromConn(conn2)
	closed1 := false
	closed2 := false
	connHasSent1 := false

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				conn2.Close()
				closed2 = true
			} else {
				conn2.Write(b1)
				connHasSent1 = true
			}
		case b2 := <-chan2:
			if b2 == nil {
				conn1.Close()
				closed1 = true
			} else {
				if connHasSent1 {
					foo, err := fn(conn2, b2)
					if err != nil {
						conn2.Close()
						conn1.Close()
						<-chan1
						<-chan2
						break
					}
					connHasSent1 = false
					b2 = foo.Data
				}
				conn1.Write(b2)
			}
		}
		if closed1 && closed2 {
			return nil
		}
	}
}
