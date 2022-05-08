package main

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/apiserver"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/etcdrw"
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var connMap sync.Map
var unknownMap sync.Map

type NodeConn struct {
	UID        string
	LastUpdate time.Time
	Conn       *net.Conn
}

func handle(conn *net.Conn) {
	defer func() {
		log.Printf("Connection with %s closed\n", (*conn).RemoteAddr().String())
		_ = (*conn).Close()
	}()
	reader := bufio.NewReader(*conn)

	key, _ := uuid.NewUUID()
	unknownConn := NodeConn{
		UID:        key.String(),
		LastUpdate: time.Now(),
		Conn:       conn,
	}
	unknownMap.Store(unknownConn.UID, unknownConn)

	buf, err := reader.ReadBytes(apiserver.MSG_DELIM)
	if err != nil {
		log.Println("Fail to read from conn, err:", err)
		return
	}
	unknownMap.Delete(unknownConn.UID)
	buf = buf[:len(buf)-1]
	var node object.Node
	err = json.Unmarshal(buf, &node)
	if err != nil || node.UID == "" || node.Status == nil {
		log.Println("Fail to parse Node, err:", err)
		return
	}

	oldBuf, err := etcdrw.GetObj("/apis/node/" + node.UID)
	if err != nil || oldBuf == nil {
		log.Printf("Node UID=%s not found, err: %v\n", node.UID, err)
		return
	}

	nodeStr := string(buf)
	nodeConn := NodeConn{
		UID:        node.UID,
		LastUpdate: time.Now(),
		Conn:       conn,
	}
	connMap.Store(node.UID, nodeConn)
	log.Printf("Updating Node UID=%s, ready=%v into etcd\n", node.UID, node.Status.Condition.Ready)
	err = etcdrw.PutObj("/apis/node/"+node.UID, nodeStr)
	if err != nil {
		log.Printf("Fail to put Node UID=%s into etcd, err: %v\n", node.UID, err)
		return
	}

	defer func() {
		node.Status.Condition.Ready = false
		buf, err = json.Marshal(node)
		if err != nil {
			log.Println("Fail to marshal Node, err: ", err)
			return
		}
		log.Printf("Updating Node UID=%s, ready=%v into etcd\n", node.UID, false)
		err = etcdrw.PutObj("/apis/node/"+node.UID, nodeStr)
		if err != nil {
			log.Printf("Fail to put Node UID=%s into etcd, err: %v\n", node.UID, err)
		}
	}()

	for {
		buf, err = reader.ReadBytes(apiserver.MSG_DELIM)
		if err != nil {
			log.Println("Fail to read from conn, err: ", err)
			return
		}
		buf = buf[:len(buf)-1]

		nodeConn.LastUpdate = time.Now()
		newNodeStr := string(buf)
		if nodeStr != newNodeStr {
			var newNode object.Node
			err = json.Unmarshal(buf, &newNode)
			if err != nil || newNode.UID != node.UID || newNode.Status == nil {
				log.Println("Fail to parse Node, err:", err)
				return
			}

			log.Printf("Updating Node UID=%s, ready=%v into etcd\n", newNode.UID, newNode.Status.Condition.Ready)
			err = etcdrw.PutObj("/apis/node/"+newNode.UID, newNodeStr)
			if err != nil {
				log.Printf("Fail to put Node UID=%s into etcd, err: %v\n", newNode.UID, err)
				return
			}
			nodeStr = newNodeStr
			node = newNode
		}
		connMap.Store(node.UID, nodeConn)
	}
}

func checkHealth() {
	for {
		time.Sleep(5 * time.Second)
		var toDelete []string
		var unkDelete []string
		connMap.Range(func(key, value any) bool {
			nodeConn := value.(NodeConn)
			if time.Since(nodeConn.LastUpdate) >= 15*time.Second {
				log.Println("Timeout, closing connection with ", (*nodeConn.Conn).RemoteAddr())
				toDelete = append(toDelete, nodeConn.UID)
				_ = (*nodeConn.Conn).Close()
				return true
			}

			_, err := (*nodeConn.Conn).Write([]byte(apiserver.MSG_HEARTBEAT))
			if err != nil {
				log.Println("Fail to send Heartbeat message, err: ", err)
				toDelete = append(toDelete, nodeConn.UID)
				_ = (*nodeConn.Conn).Close()
			}
			return true
		})

		unknownMap.Range(func(key, value any) bool {
			nodeConn := value.(NodeConn)
			if time.Since(nodeConn.LastUpdate) >= 15*time.Second {
				log.Println("Timeout, closing connection with ", (*nodeConn.Conn).RemoteAddr())
				unkDelete = append(unkDelete, nodeConn.UID)
				_ = (*nodeConn.Conn).Close()
			}
			return true
		})

		for _, UID := range toDelete {
			connMap.Delete(UID)
		}

		for _, UID := range unkDelete {
			unknownMap.Delete(UID)
		}
	}
}

func listenHeartbeat() {
	connMap = sync.Map{}
	unknownMap = sync.Map{}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(cubeconfig.HeartbeatPort))
	if err != nil {
		log.Fatal("Failure when listening heartbeat, err: ", err)
		return
	}
	defer func() { _ = listener.Close() }()

	go checkHealth()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failure when accepting heartbeat, err: ", err)
			continue
		}
		go handle(&conn)
	}
}
