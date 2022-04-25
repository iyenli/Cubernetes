package watchobj

import (
	cubeconfig "Cubernetes/config"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func watching(url string, ch chan string, ctx context.Context) {
	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	data := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readN, err := resp.Body.Read(data)
			if readN > 0 {
				ch <- string(data[:readN])
			}
			if err == io.EOF {
				log.Println("connection closed by server")
				return
			}
			if err != nil {
				panic(err)
			}
		}
	}
}

func WatchObj(path string) (chan string, context.CancelFunc) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + path
	ch := make(chan string)
	ctx, cancel := context.WithCancel(context.TODO())
	go watching(url, ch, ctx)
	return ch, cancel
}
