package proxyruntime

type ServiceChainElement struct {
	probabilityChainUid [][]string
	serviceChainUid     []string
	numberOfPods        int
}

type PodChainElement struct {
}
