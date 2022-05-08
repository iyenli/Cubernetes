package apiserver

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)

const MSG_DELIM byte = 26
const MSG_HEARTBEAT = "alive" + string(MSG_DELIM)

var conn net.Conn
var node object.Node
var lastUpdate time.Time

func listenHeartbeat(reader *bufio.Reader) {
	defer func() { _ = conn.Close() }()

	for {
		lastUpdate = time.Now()
		_, err := reader.ReadBytes(MSG_DELIM)
		if err != nil {
			log.Println("Fail to read from conn")
			return
		}
		if time.Since(lastUpdate) > 15*time.Second {
			log.Println("Timeout, close conn")
			_ = conn.Close()
			return
		}
	}
}

func sendHeartBeat() {
	var err error
	conn, err = net.Dial("tcp", cubeconfig.APIServerIp+":"+strconv.Itoa(cubeconfig.HeartbeatPort))
	if err != nil {
		log.Fatal("Fail to dial heartbeat server", err)
	}

	defer func() { _ = conn.Close() }()

	reader := bufio.NewReader(conn)
	go listenHeartbeat(reader)

	for {
		time.Sleep(5 * time.Second)
		buf, err := json.Marshal(node)
		if err != nil {
			log.Println("Fail to marshal Node")
			return
		}
		buf = append(buf, MSG_DELIM)
		_, err = conn.Write(buf)
		if err != nil {
			log.Println("Fail to send Node message")
			return
		}
	}
}

func Init(n object.Node) {
	node = n
	go sendHeartBeat()
}
