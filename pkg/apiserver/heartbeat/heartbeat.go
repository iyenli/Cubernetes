package heartbeat

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/apiserver/health"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const MSG_DELIM byte = 26
const MSG_HEARTBEAT = "alive"

const INTERVAL = 5 * time.Second
const TIMEOUT = 15 * time.Second

var lostCnt int
var connected bool
var node object.Node
var lastUpdate time.Time

var timeLock sync.Mutex
var connLock sync.Mutex

// CheckConn
// package heartbeat is designed for cubelet
func CheckConn() bool {
	connLock.Lock()
	ret := connected
	connLock.Unlock()
	return ret
}

func closeConn(conn *net.Conn, cnt int) {
	_ = (*conn).Close()

	connLock.Lock()
	if lostCnt == cnt {
		lostCnt += 1
		connected = false
		log.Println("Heartbeat connection closed, stop all watching")
		watchobj.StopAll()
		go retryConnection()
	}
	connLock.Unlock()
}

func updateTime() {
	timeLock.Lock()
	lastUpdate = time.Now()
	timeLock.Unlock()
}

func getTime() time.Time {
	timeLock.Lock()
	t := lastUpdate
	timeLock.Unlock()
	return t
}

func maintainHealth(conn *net.Conn, cnt int) {
	defer closeConn(conn, cnt)

	for {
		if time.Since(getTime()) > TIMEOUT {
			log.Println("Timeout, close conn")
			return
		}

		buf, err := json.Marshal(node)
		if err != nil {
			log.Println("Fail to marshal Node, err: ", err)
			return
		}
		buf = append(buf, MSG_DELIM)
		_, err = (*conn).Write(buf)
		if err != nil {
			log.Println("Fail to send Node message, this is NORMAL when apiserver is lost, err: ", err)
			return
		}
		time.Sleep(INTERVAL)
	}
}

func updateHeartBeat() {
	connLock.Lock()
	cnt := lostCnt
	connLock.Unlock()

	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", cubeconfig.APIServerIp+":"+strconv.Itoa(cubeconfig.HeartbeatPort))
		if err == nil {
			break
		}
		log.Println("Cannot connect with ApiServer, retry...")
		time.Sleep(INTERVAL)
	}

	log.Println("Connected with ApiServer, sending heartbeat...")
	updateTime()

	connLock.Lock()
	connected = true
	connLock.Unlock()

	defer closeConn(&conn, cnt)

	go maintainHealth(&conn, cnt)

	reader := bufio.NewReader(conn)
	for {
		_, err := reader.ReadBytes(MSG_DELIM)
		if err != nil {
			log.Println("Fail to read from conn")
			return
		}
		updateTime()
	}
}

func retryConnection() {
	for {
		log.Println("Cannot access ApiServer, retry...")
		time.Sleep(INTERVAL)
		if health.CheckApiServerHealth() {
			break
		}
	}
	log.Println("ApiServer alive, connecting...")

	go updateHeartBeat()
}

// InitNode
// package heartbeat is designed for cubelet
func InitNode(n object.Node) {
	node = n
	node.Status.Condition.Ready = true
	connected = false
	lostCnt = 0
	go updateHeartBeat()
}
