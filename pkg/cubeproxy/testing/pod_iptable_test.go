package testing

//func TestAddPod(t *testing.T) {
//	ptr, err := proxyruntime.InitPodChain()
//	assert.NoError(t, err)
//
//	pod, err := proxyruntime.GetPodByService(nil)
//	assert.NoError(t, err)
//	assert.NotNil(t, pod)
//	assert.Equal(t, 1, len(pod))
//
//	containerIP := net.ParseIP("10.0.0.4")
//	t.Log(pod[0].Status.IP.String())
//	err = ptr.AddPod(&pod[0], containerIP)
//	assert.NoError(t, err)
//
//	// Now check IP Table manually...
//}
//
//func TestDeletePod(t *testing.T) {
//	ptr, err := proxyruntime.InitPodChain()
//	assert.NoError(t, err)
//
//	pod, err := proxyruntime.GetPodByService(nil)
//	assert.NoError(t, err)
//	assert.NotNil(t, pod)
//	assert.Equal(t, 1, len(pod))
//
//	containerIP := net.ParseIP("10.0.0.4")
//	err = ptr.DeletePod(&pod[0], containerIP)
//	assert.NoError(t, err)
//
//	// Now check IP Table manually...
//}
