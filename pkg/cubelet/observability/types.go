package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}
