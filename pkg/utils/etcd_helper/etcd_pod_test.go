package etcd_helper

import (
	"Cubernetes/pkg/object"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
	"testing"
)

var testCasesName = []string{"test1", "test2"}
var testCasePod = []object.Pod{
	{
		TypeMeta: object.TypeMeta{
			Kind:       "pod",
			APIVersion: "0.1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: testCasesName[0],
		},
		Spec: object.PodSpec{
			Containers: []object.Container{
				{
					Name:  testCasesName[0],
					Image: testCasesName[0],
				},
			},
		},
		Status: nil,
	},
	{
		TypeMeta: object.TypeMeta{
			Kind:       "pod",
			APIVersion: "0.1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: testCasesName[1],
		},
		Spec: object.PodSpec{
			Containers: []object.Container{
				{
					Name:  testCasesName[1],
					Image: testCasesName[1],
				},
			},
		},
		Status: nil,
	}}

func TestStorePod(t *testing.T) {
	ctx := ETCDContext{Client: NewETCDClient()}
	defer CloseETCDClient(ctx.Client)

	for i := 0; i < len(testCasesName); i++ {
		res1, err1 := StorePod(&ctx, testCasesName[i], testCasePod[i])
		assert.Equal(t, nil, err1)
		assert.Equal(t, true, res1)

		res2, err2 := GetPods(&ctx, testCasesName[i])
		assert.Equal(t, nil, err2)
		assert.Equal(t, 1, len(res2))
		assert.Equal(t, res2[0], testCasePod[i])
	}
}

func TestHealthCheck(t *testing.T) {
	ctx := ETCDContext{Client: NewETCDClient()}
	defer CloseETCDClient(ctx.Client)

	res := ETCDHealthCheck(&ctx)
	assert.Equal(t, true, res)
}

func TestDeletePodAndAllPods(t *testing.T) {
	ctx := ETCDContext{Client: NewETCDClient()}
	defer CloseETCDClient(ctx.Client)

	for i := 0; i < len(testCasesName); i++ {
		res1, _ := StorePod(&ctx, testCasesName[i], testCasePod[i])
		assert.Equal(t, true, res1)

		res2, _ := GetPods(&ctx, testCasesName[i])
		assert.Equal(t, res2[0], testCasePod[i])
	}

	res3, _ := GetAllPods(&ctx)
	assert.Equal(t, len(testCasesName), len(res3))

	res3, _ = GetPodsRange(&ctx, testCasesName[0], testCasesName[len(testCasesName)-1])
	assert.Equal(t, len(testCasesName)-1, len(res3))

	res4, _ := DeletePod(&ctx, testCasesName[0])
	assert.Equal(t, true, res4)

	res3, _ = GetAllPods(&ctx)
	assert.Equal(t, len(testCasesName)-1, len(res3))

	res5, _ := DeletePod(&ctx, "not-exist")
	assert.Equal(t, false, res5)

	res3, _ = GetAllPods(&ctx)
	assert.Equal(t, len(testCasesName)-1, len(res3))
}

func TestWatcher(t *testing.T) {
	ctx := ETCDContext{Client: NewETCDClient()}
	defer CloseETCDClient(ctx.Client)

	watcher := GetPodWatcher(&ctx, "test", false)

	go func() {
		for i := 0; i < len(testCasesName); i++ {
			_, _ = StorePod(&ctx, testCasesName[i], testCasePod[i])
		}
		res, _ := GetAllPods(&ctx)
		assert.Equal(t, len(testCasesName), len(res))
	}()

	go func() {
		index := 0
		for res := range watcher {
			assert.Equal(t, 1, len(res.Events))
			assert.Equal(t, testCasesName[index], res.Events[0].Kv.Key)
			assert.Equal(t, testCasePod[index], res.Events[0].Kv.Value)
			assert.Equal(t, clientv3.EventTypePut, res.Events[0].Type)
			index++
		}
	}()
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
//	Client, err := clientv3.New(clientv3.Config{
//		Endpoints:   e.Server.Cluster().ClientURLs(),
//		DialTimeout: 10 * time.Second,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	return Client
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
