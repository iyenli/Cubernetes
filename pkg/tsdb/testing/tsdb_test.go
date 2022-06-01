package testing

import (
	"Cubernetes/pkg/tsdb"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTSDB(t *testing.T) {
	tsdb.Init()

	go insert("a")
	go insert("c")
	go insert("b")
	go insert("a")
	go insert("b")
	go insert("c")

	fmt.Println(tsdb.QueryLastMinute("a"))
	fmt.Println(tsdb.QueryLastMinute("b"))
	fmt.Println(tsdb.QueryLastMinute("c"))

	time.Sleep(500 * time.Millisecond)

	fmt.Println(tsdb.QueryLastMinute("a"))
	fmt.Println(tsdb.QueryLastMinute("b"))
	fmt.Println(tsdb.QueryLastMinute("c"))

	time.Sleep(500 * time.Millisecond)

	fmt.Println(tsdb.QueryLastMinute("a"))
	fmt.Println(tsdb.QueryLastMinute("b"))
	fmt.Println(tsdb.QueryLastMinute("c"))

	time.Sleep(500 * time.Millisecond)

	fmt.Println(tsdb.QueryLastMinute("a"))
	fmt.Println(tsdb.QueryLastMinute("b"))
	fmt.Println(tsdb.QueryLastMinute("c"))
}

func insert(name string) {
	for i := 0; i < 1000; i += 1 {
		t := time.Now().UnixMilli()
		tsdb.Insert(t+rand.Int63n(20), name)
		time.Sleep(time.Millisecond)
	}
}
