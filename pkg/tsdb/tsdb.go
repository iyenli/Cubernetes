package tsdb

import (
	"container/list"
	"sync"
	"time"
)

const second = 1000

var db map[string]*TimeSeries
var dbMtx sync.Mutex

type TimeSeries struct {
	mtx    sync.Mutex
	series list.List
}

func Init() {
	db = make(map[string]*TimeSeries)
}

func Insert(ms int64, name string) {
	curTime := time.Now().UnixMilli()

	dbMtx.Lock()
	ts := db[name]
	if ts == nil {
		ts = &TimeSeries{
			mtx:    sync.Mutex{},
			series: list.List{},
		}
		db[name] = ts
	}
	dbMtx.Unlock()

	ts.mtx.Lock()
	for ts.series.Len() > 0 && curTime-ts.series.Front().Value.(int64) > second {
		ts.series.Remove(ts.series.Front())
	}

	it := ts.series.Back()
	if it == nil {
		ts.series.PushBack(ms)
		ts.mtx.Unlock()
		return
	}

	for {
		if it.Value.(int64) <= ms {
			ts.series.InsertAfter(ms, it)
			break
		}

		it = it.Prev()
		if it == nil {
			ts.series.PushFront(ms)
			break
		}
	}
	ts.mtx.Unlock()
}

func QueryLastMinute(name string) int {
	curTime := time.Now().UnixMilli()

	dbMtx.Lock()
	ts := db[name]
	dbMtx.Unlock()

	if ts == nil {
		return 0
	}

	ts.mtx.Lock()

	for ts.series.Len() > 0 && curTime-ts.series.Front().Value.(int64) > second {
		ts.series.Remove(ts.series.Front())
	}

	stat := ts.series.Len()
	ts.mtx.Unlock()
	return stat
}
