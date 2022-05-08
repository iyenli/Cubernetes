package socketrw

import (
	"bytes"
	"net"
)

const MSG_DELIM byte = 26

func Write(conn net.Conn, content string) (int, error) {
	var buf bytes.Buffer
	buf.WriteString(content)
	buf.WriteByte(MSG_DELIM)

	return conn.Write(buf.Bytes())
}

func Read(conn net.Conn) (string, error) {
	var buf bytes.Buffer
	arr := make([]byte, 1)
	for {
		if _, err := conn.Read(arr); err != nil {
			return "", err
		}
		item := arr[0]
		if item == MSG_DELIM {
			break
		}
		buf.WriteByte(item)
	}
	str := buf.String()
	return str, nil
}
