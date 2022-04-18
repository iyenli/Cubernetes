package etcd_helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCase = []string{"test1", "pod-context"}

func TestStorePod(t *testing.T) {
	ctx := ETCDContext{client: newETCDClient()}
	defer closeETCDClient(ctx.client)

	res1, err1 := storePod(ctx, testCase[0], testCase[1])
	assert.Equal(t, nil, err1)
	assert.Equal(t, true, res1)

	res2, err2 := getPods(ctx, testCase[0])
	assert.Equal(t, nil, err2)
	assert.Equal(t, 1, len(res2))
	assert.Equal(t, res2[0], []byte(testCase[1]))
}

func TestHealthCheck(t *testing.T) {
	ctx := ETCDContext{client: newETCDClient()}
	defer closeETCDClient(ctx.client)

	res := ETCDHealthCheck(ctx)
	assert.Equal(t, true, res)
}

/*
 Can't pass compilation
*/
//func RunEtcd(t *testing.T, cfg *embed.Config) *clientv3.Client {
//	t.Helper()
//
//	e, err := embed.StartEtcd(cfg)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Cleanup(e.Close)
//
//	select {
//	case <-e.Server.ReadyNotify():
//	case <-time.After(60 * time.Second):
//		e.Server.Stop()
//		t.Fatal("server took too long to start")
//	}
//	go func() {
//		err := <-e.Err()
//		if err != nil {
//			t.Error(err)
//		}
//	}()
//
//	client, err := clientv3.New(clientv3.Config{
//		Endpoints:   e.Server.Cluster().ClientURLs(),
//		DialTimeout: 10 * time.Second,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	return client
//}
//
//func getAvailablePorts(count int) ([]int, error) {
//	ports := []int{}
//	for i := 0; i < count; i++ {
//		l, err := net.Listen("tcp", ":0")
//		if err != nil {
//			return nil, fmt.Errorf("could not bind to a port: %v", err)
//		}
//		// It is possible but unlikely that someone else will bind this port before we get a chance to use it.
//		defer l.Close()
//		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
//	}
//	return ports, nil
//}
