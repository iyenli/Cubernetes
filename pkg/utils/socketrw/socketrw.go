package socketrw

import (
	"bytes"
	"fmt"
	"net"
)

const delimiter byte = 0

func Write(conn net.Conn, content string) (int, error) {
	fmt.Printf("send %v: %v\n", conn.RemoteAddr(), content)
	var buf bytes.Buffer
	buf.WriteString(content)
	buf.WriteByte(delimiter)

	return conn.Write(buf.Bytes())
}

func Read(conn net.Conn) (string, error) {
	var str string
	var buf bytes.Buffer
	arr := make([]byte, 1)
	for {
		if _, err := conn.Read(arr); err != nil {
			return str, err
		}
		item := arr[0]
		if item == delimiter {
			break
		}
		buf.WriteByte(item)
	}
	str = buf.String()
	fmt.Printf("recv %v: %v\n", conn.RemoteAddr(), str)
	return str, nil
}
